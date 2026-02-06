package plugin

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	plugin_ent "github.com/shuTwT/hoshikuzu/ent/plugin"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	plugin_lib "github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v3"
)

type PluginService interface {
	ListPluginPage(ctx context.Context, page, size int) (int, []*ent.Plugin, error)
	ListPluginPageWithQuery(ctx context.Context, req model.PluginPageReq) (int, []*ent.Plugin, error)
	QueryPlugin(ctx context.Context, id int) (*ent.Plugin, error)
	CreatePlugin(ctx context.Context, fileHeader *multipart.FileHeader) (*ent.Plugin, error)
	DeletePlugin(ctx context.Context, id int) error
	StartPlugin(ctx context.Context, id int) error
	StopPlugin(ctx context.Context, id int) error
	RestartPlugin(ctx context.Context, id int) error
	AutoStartPlugins(ctx context.Context) error
	RegisterPlugin(ctx context.Context, pluginInfo *model.PluginRegisterReq) error
	HeartbeatPlugin(ctx context.Context, heartbeatInfo *model.PluginHeartbeatReq) error
	CheckPluginTimeout(ctx context.Context) error
}

type PluginServiceImpl struct {
	client           *ent.Client
	pluginClients    map[int]*plugin_lib.Client
	pluginCache      map[string]map[string]interface{}
	pluginHeartbeats map[string]time.Time
	mu               sync.RWMutex
}

type emptyPlugin struct {
	plugin_lib.Plugin
}

func (emptyPlugin) Server(*plugin_lib.MuxBroker) (interface{}, error) {
	return &emptyPluginRPCServer{Impl: &emptyPlugin{}}, nil
}

func (emptyPlugin) Client(b *plugin_lib.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &emptyPluginRPCClient{client: c}, nil
}

type emptyPluginRPCServer struct {
	Impl *emptyPlugin
}

func (s *emptyPluginRPCServer) Init() error {
	return nil
}

func (s *emptyPluginRPCServer) Ping() error {
	return nil
}

func (s *emptyPluginRPCServer) Close() error {
	return nil
}

type emptyPluginRPCClient struct {
	client *rpc.Client
}

func (c *emptyPluginRPCClient) Close() error {
	return nil
}

func NewPluginServiceImpl(client *ent.Client) *PluginServiceImpl {
	return &PluginServiceImpl{
		client:           client,
		pluginClients:    make(map[int]*plugin_lib.Client),
		pluginCache:      make(map[string]map[string]interface{}),
		pluginHeartbeats: make(map[string]time.Time),
	}
}

func (s *PluginServiceImpl) ListPluginPage(ctx context.Context, page, size int) (int, []*ent.Plugin, error) {
	count, err := s.client.Plugin.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	plugins, err := s.client.Plugin.Query().
		Order(ent.Desc(plugin_ent.FieldID)).
		Offset((page - 1) * size).
		Limit(size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, plugins, nil
}

func (s *PluginServiceImpl) ListPluginPageWithQuery(ctx context.Context, req model.PluginPageReq) (int, []*ent.Plugin, error) {
	query := s.client.Plugin.Query()

	if req.Name != "" {
		query.Where(plugin_ent.NameContains(req.Name))
	}

	if req.Key != "" {
		query.Where(plugin_ent.KeyContains(req.Key))
	}

	if req.Status != "" {
		query.Where(plugin_ent.StatusEQ(plugin_ent.Status(req.Status)))
	}

	if req.Enabled != nil {
		query.Where(plugin_ent.Enabled(*req.Enabled))
	}

	if req.AutoStart != nil {
		query.Where(plugin_ent.AutoStart(*req.AutoStart))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	plugins, err := query.
		Order(ent.Desc(plugin_ent.FieldID)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, plugins, nil
}

func (s *PluginServiceImpl) QueryPlugin(ctx context.Context, id int) (*ent.Plugin, error) {
	pluginEntity, err := s.client.Plugin.Query().
		Where(plugin_ent.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return pluginEntity, nil
}

func (s *PluginServiceImpl) CreatePlugin(ctx context.Context, fileHeader *multipart.FileHeader) (*ent.Plugin, error) {
	if fileHeader == nil {
		return nil, errors.New("文件不能为空")
	}

	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer srcFile.Close()

	tempFile, err := os.CreateTemp("", "plugin-*.zip")
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, srcFile); err != nil {
		return nil, fmt.Errorf("复制文件失败: %w", err)
	}

	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("打开压缩包失败: %w", err)
	}
	defer zipReader.Close()

	var configContent []byte
	var binaryFile *zip.File
	pluginDir := ""

	for _, f := range zipReader.File {
		if f.Name == "plugin-config.yaml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("打开配置文件失败: %w", err)
			}
			configContent, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("读取配置文件失败: %w", err)
			}
		} else if strings.HasSuffix(f.Name, "/") {
			if pluginDir == "" {
				pluginDir = strings.TrimSuffix(f.Name, "/")
			}
		} else {
			ext := filepath.Ext(f.Name)
			if ext == "" || strings.Contains(strings.ToLower(f.Name), "bin") {
				binaryFile = f
			}
		}
	}

	if configContent == nil {
		return nil, errors.New("压缩包中未找到 plugin-config.yaml 文件")
	}

	var pluginConfig model.PluginConfig
	if err := yaml.Unmarshal(configContent, &pluginConfig); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if err := validatePluginConfig(&pluginConfig); err != nil {
		return nil, err
	}

	exists, err := s.client.Plugin.Query().Where(plugin_ent.Key(pluginConfig.Key)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("检查插件是否存在失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("插件 key '%s' 已存在", pluginConfig.Key)
	}

	if len(pluginConfig.Dependencies) > 0 {
		for _, depKey := range pluginConfig.Dependencies {
			depExists, err := s.client.Plugin.Query().Where(plugin_ent.Key(depKey)).Exist(ctx)
			if err != nil {
				return nil, fmt.Errorf("检查依赖插件 '%s' 失败: %w", depKey, err)
			}
			if !depExists {
				return nil, fmt.Errorf("依赖插件 '%s' 不存在", depKey)
			}
		}
	}

	pluginsDir := "./data/plugins"
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return nil, fmt.Errorf("创建插件目录失败: %w", err)
	}

	targetDir := filepath.Join(pluginsDir, pluginConfig.Key)
	if err := os.RemoveAll(targetDir); err != nil {
		return nil, fmt.Errorf("清理旧插件目录失败: %w", err)
	}

	for _, f := range zipReader.File {
		targetPath := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return nil, fmt.Errorf("创建目录失败: %w", err)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return nil, fmt.Errorf("创建父目录失败: %w", err)
			}
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("打开压缩文件失败: %w", err)
			}
			defer rc.Close()

			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return nil, fmt.Errorf("创建文件失败: %w", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, rc); err != nil {
				return nil, fmt.Errorf("解压文件失败: %w", err)
			}
		}
	}

	var binPath string
	if binaryFile != nil {
		binPath = filepath.Join(targetDir, binaryFile.Name)
		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			return nil, errors.New("压缩包中未找到二进制文件")
		}
		if err := os.Chmod(binPath, 0755); err != nil {
			return nil, fmt.Errorf("设置二进制文件权限失败: %w", err)
		}
	} else {
		return nil, errors.New("压缩包中未找到二进制文件")
	}

	builder := s.client.Plugin.Create().
		SetKey(pluginConfig.Key).
		SetName(pluginConfig.Name).
		SetVersion(pluginConfig.Version).
		SetBinPath(binPath).
		SetMagicCookieValue(pluginConfig.MagicCookieValue).
		SetDependencies(pluginConfig.Dependencies).
		SetEnabled(true).
		SetAutoStart(pluginConfig.AutoStart).
		SetStatus("stopped")

	if pluginConfig.Description != "" {
		builder.SetDescription(pluginConfig.Description)
	}

	if pluginConfig.ProtocolVersion != "" {
		parsedVersion, err := strconv.ParseUint(pluginConfig.ProtocolVersion, 10, 32)
		if err == nil {
			builder.SetProtocolVersion(uint(parsedVersion))
		} else {
			builder.SetProtocolVersion(1)
		}
	} else {
		builder.SetProtocolVersion(1)
	}

	if pluginConfig.MagicCookieKey != "" {
		builder.SetMagicCookieKey(pluginConfig.MagicCookieKey)
	} else {
		builder.SetMagicCookieKey("GO_PLUGIN")
	}

	if pluginConfig.Config != "" {
		builder.SetConfig(pluginConfig.Config)
	}

	pluginEntity, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("保存插件信息失败: %w", err)
	}

	return pluginEntity, nil
}

func (s *PluginServiceImpl) DeletePlugin(ctx context.Context, id int) error {
	pluginEntity, err := s.client.Plugin.Query().Where(plugin_ent.ID(id)).First(ctx)
	if err != nil {
		return err
	}

	if pluginEntity.Status == "running" {
		if err := s.StopPlugin(ctx, id); err != nil {
			return fmt.Errorf("停止插件失败: %w", err)
		}
	}

	pluginDir := filepath.Join("./data/plugins", pluginEntity.Key)
	if err := os.RemoveAll(pluginDir); err != nil {
		return fmt.Errorf("删除插件目录失败: %w", err)
	}

	err = s.client.Plugin.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	s.mu.Lock()
	delete(s.pluginClients, id)
	s.mu.Unlock()

	return nil
}

func (s *PluginServiceImpl) StartPlugin(ctx context.Context, id int) error {
	pluginEntity, err := s.client.Plugin.Query().Where(plugin_ent.ID(id)).First(ctx)
	if err != nil {
		return err
	}

	if !pluginEntity.Enabled {
		return errors.New("插件未启用")
	}

	if _, err := os.Stat(pluginEntity.BinPath); os.IsNotExist(err) {
		return fmt.Errorf("插件二进制文件不存在: %s", pluginEntity.BinPath)
	}

	if len(pluginEntity.Dependencies) > 0 {
		for _, depKey := range pluginEntity.Dependencies {
			depPlugin, err := s.client.Plugin.Query().Where(plugin_ent.Key(depKey)).First(ctx)
			if err != nil {
				return fmt.Errorf("获取依赖插件 '%s' 失败: %w", depKey, err)
			}
			if depPlugin.Status != "running" {
				return fmt.Errorf("依赖插件 '%s' 未运行", depKey)
			}
		}
	}

	s.mu.RLock()
	if client, exists := s.pluginClients[id]; exists {
		s.mu.RUnlock()
		if _, err := client.Client(); err == nil {
			return errors.New("插件已在运行中")
		}
		client.Kill()
		delete(s.pluginClients, id)
		s.mu.RUnlock()
	} else {
		s.mu.RUnlock()
	}

	now := time.Now()
	err = s.client.Plugin.UpdateOneID(id).
		SetStatus("loading").
		SetLastStartedAt(now).
		SetLastError("").
		Exec(ctx)
	if err != nil {
		return err
	}

	go func(pluginID int, pluginKey string) {
		handshakeConfig := plugin_lib.HandshakeConfig{
			ProtocolVersion:  pluginEntity.ProtocolVersion,
			MagicCookieKey:   pluginEntity.MagicCookieKey,
			MagicCookieValue: pluginEntity.MagicCookieValue,
		}

		client := plugin_lib.NewClient(&plugin_lib.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         map[string]plugin_lib.Plugin{pluginKey: &emptyPlugin{}},
			Cmd:             exec.Command(pluginEntity.BinPath),
			Managed:         true,
		})

		rpcClient, err := client.Client()
		if err != nil {
			s.client.Plugin.UpdateOneID(pluginID).
				SetStatus("error").
				SetLastError(fmt.Sprintf("插件连接失败: %s", err.Error())).
				Exec(context.Background())
			return
		}

		_, err = rpcClient.Dispense(pluginKey)
		if err != nil {
			s.client.Plugin.UpdateOneID(pluginID).
				SetStatus("error").
				SetLastError(fmt.Sprintf("插件分发失败: %s", err.Error())).
				Exec(context.Background())
			client.Kill()
			return
		}

		s.mu.Lock()
		s.pluginClients[pluginID] = client
		s.mu.Unlock()

		s.client.Plugin.UpdateOneID(pluginID).
			SetStatus("running").
			Exec(context.Background())
	}(id, pluginEntity.Key)

	return nil
}

func (s *PluginServiceImpl) StopPlugin(ctx context.Context, id int) error {
	s.mu.RLock()
	client, exists := s.pluginClients[id]
	s.mu.RUnlock()

	if !exists {
		return errors.New("插件未运行")
	}

	client.Kill()

	s.mu.Lock()
	delete(s.pluginClients, id)
	s.mu.Unlock()

	now := time.Now()
	err := s.client.Plugin.UpdateOneID(id).
		SetStatus("stopped").
		SetLastStoppedAt(now).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *PluginServiceImpl) RestartPlugin(ctx context.Context, id int) error {
	pluginEntity, err := s.client.Plugin.Query().Where(plugin_ent.ID(id)).First(ctx)
	if err != nil {
		return err
	}

	if pluginEntity.Status == "running" {
		if err := s.StopPlugin(ctx, id); err != nil {
			return fmt.Errorf("停止插件失败: %w", err)
		}
		time.Sleep(1 * time.Second)
	}

	if err := s.StartPlugin(ctx, id); err != nil {
		return fmt.Errorf("启动插件失败: %w", err)
	}

	return nil
}

func (s *PluginServiceImpl) AutoStartPlugins(ctx context.Context) error {
	plugins, err := s.client.Plugin.Query().
		Where(plugin_ent.AutoStart(true)).
		Where(plugin_ent.Enabled(true)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("获取自动启动插件失败: %w", err)
	}

	for _, p := range plugins {
		if p.Status != "running" {
			go func(pluginID int) {
				if err := s.StartPlugin(context.Background(), pluginID); err != nil {
					fmt.Printf("自动启动插件 %d 失败: %v\n", pluginID, err)
				}
			}(p.ID)
		}
	}

	return nil
}

func (s *PluginServiceImpl) RegisterPlugin(ctx context.Context, pluginInfo *model.PluginRegisterReq) error {
	// 使用插件名称作为key
	pluginKey := pluginInfo.Name

	s.mu.Lock()
	s.pluginCache[pluginKey] = map[string]interface{}{
		"name":         pluginInfo.Name,
		"version":      pluginInfo.Version,
		"grpc_address": pluginInfo.GrpcAddress,
		"status":       pluginInfo.Status,
		"start_time":   pluginInfo.StartTime,
		"metadata":     pluginInfo.Metadata,
	}
	// 初始化心跳时间
	s.pluginHeartbeats[pluginKey] = time.Now()
	s.mu.Unlock()

	return nil
}

func (s *PluginServiceImpl) HeartbeatPlugin(ctx context.Context, heartbeatInfo *model.PluginHeartbeatReq) error {
	// 使用插件名称作为key
	pluginKey := heartbeatInfo.Name

	s.mu.Lock()
	// 更新心跳时间
	s.pluginHeartbeats[pluginKey] = time.Now()
	// 如果插件状态不为空，则更新插件缓存中的状态
	if heartbeatInfo.Status != "" {
		if pluginInfo, exists := s.pluginCache[pluginKey]; exists {
			pluginInfo["status"] = heartbeatInfo.Status
			s.pluginCache[pluginKey] = pluginInfo
		}
	}
	s.mu.Unlock()

	return nil
}

func (s *PluginServiceImpl) CheckPluginTimeout(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 超时时间设置为30秒
	timeout := 30 * time.Second
	now := time.Now()

	// 遍历所有插件的心跳时间
	for pluginKey, heartbeatTime := range s.pluginHeartbeats {
		// 检查心跳时间是否超过了超时时间
		if now.Sub(heartbeatTime) > timeout {
			// 如果超过了超时时间，则将插件状态置为停止(stopped)
			if pluginInfo, exists := s.pluginCache[pluginKey]; exists {
				pluginInfo["status"] = "stopped"
				s.pluginCache[pluginKey] = pluginInfo
			}
		}
	}

	return nil
}

func validatePluginConfig(config *model.PluginConfig) error {
	if config.Key == "" {
		return errors.New("插件 key 不能为空")
	}
	if config.Name == "" {
		return errors.New("插件名称不能为空")
	}
	if config.Version == "" {
		return errors.New("插件版本不能为空")
	}
	if config.MagicCookieValue == "" {
		return errors.New("Magic Cookie Value 不能为空")
	}
	return nil
}
