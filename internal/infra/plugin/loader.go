package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	go_plugin "github.com/hashicorp/go-plugin"
	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/plugin"
	"github.com/shuTwT/hoshikuzu/internal/infra/plugin/util"
	plugin_consts "github.com/shuTwT/hoshikuzu/pkg/plugin/consts"
	plugin_shared "github.com/shuTwT/hoshikuzu/pkg/plugin/shared"
)

// PluginLoader 插件加载器（管理插件进程、RPC连接）
type PluginLoader struct {
	entClient      *ent.Client
	binDir         string                       // 插件二进制解压目录（如./plugins/bin）
	pluginMap      map[string]go_plugin.Plugin  // go-plugin插件映射
	runningPlugins map[string]*go_plugin.Client // 正在运行的插件客户端（key：pluginName_version）
}

// NewPluginLoader 创建插件加载器实例
func NewPluginLoader(entClient *ent.Client, binDir string) *PluginLoader {
	return &PluginLoader{
		entClient: entClient,
		binDir:    binDir,
		pluginMap: map[string]go_plugin.Plugin{
			"plugin_store": &plugin_shared.StoreGRPCPlugin{}, // 公共包定义的插件
		},
		runningPlugins: make(map[string]*go_plugin.Client),
	}
}

// LoadPluginByNameVersion 按插件名+版本加载插件实例
func (l *PluginLoader) LoadPluginByNameVersion(ctx context.Context, pluginName, version string) (plugin_shared.PluginStore, error) {
	// 1. 校验参数
	if !util.ValidateSemanticVersion(version) {
		return nil, errors.New("invalid semantic version")
	}

	// 2. 查询数据库：获取插件版本信息
	pluginVer, err := l.entClient.Plugin.
		Query().
		Where(
			plugin.Version(version),
			plugin.Enabled(true),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("query plugin version failed: %w", err)
	}

	// 4. 构建插件二进制路径（解压后的可执行文件）
	zipPath := pluginVer.BinPath
	pluginKey := fmt.Sprintf("%s_%s", pluginName, version)
	binDstDir := filepath.Join(l.binDir, pluginName, version)
	binPath := filepath.Join(binDstDir, pluginName) // 插件二进制文件（无后缀：Linux/Mac；.exe：Windows）

	// 5. 解压zip包（若未解压）
	if _, err := os.Stat(binPath); err != nil {
		if os.IsNotExist(err) {
			// 解压zip到目标目录
			if err := util.Unzip(zipPath, binDstDir); err != nil {
				return nil, fmt.Errorf("unzip plugin failed: %w", err)
			}
		} else {
			return nil, fmt.Errorf("check bin file failed: %w", err)
		}
	}

	// 6. 检查插件是否已运行，若已运行直接返回实例
	if client, ok := l.runningPlugins[pluginKey]; ok {
		if client.Exited() {
			delete(l.runningPlugins, pluginKey)
		} else {
			clientProtocol, err := client.Client()
			raw, err := clientProtocol.Dispense("plugin_store")
			if err != nil {
				return nil, err
			}
			return raw.(plugin_shared.PluginStore), nil
		}
	}

	// 7. 配置go-plugin客户端
	client := go_plugin.NewClient(&go_plugin.ClientConfig{
		HandshakeConfig: go_plugin.HandshakeConfig{
			ProtocolVersion:  plugin_consts.ProtocolVersion,
			MagicCookieKey:   plugin_consts.MagicCookieKey,
			MagicCookieValue: plugin_consts.MagicCookieValue,
		},
		Plugins:          l.pluginMap,
		Cmd:              exec.Command(binPath), // 启动插件进程
		AllowedProtocols: []go_plugin.Protocol{go_plugin.ProtocolGRPC},
	})

	// 8. 建立RPC连接，获取插件实例
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("create rpc client failed: %w", err)
	}

	raw, err := rpcClient.Dispense("plugin_store")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("dispense plugin failed: %w", err)
	}

	// 9. 记录正在运行的插件
	l.runningPlugins[pluginKey] = client

	// 10. 插件初始化（注入配置）
	pluginInstance := raw.(plugin_shared.PluginStore)
	if err := pluginInstance.Init(ctx, map[string]string{
		"plugin_name": pluginName,
		"version":     version,
	}); err != nil {
		client.Kill()
		delete(l.runningPlugins, pluginKey)
		return nil, fmt.Errorf("plugin init failed: %w", err)
	}

	return pluginInstance, nil
}

// UnloadPluginByNameVersion 卸载插件（停止进程、释放资源）
func (l *PluginLoader) UnloadPluginByNameVersion(pluginName, version string) error {
	pluginKey := fmt.Sprintf("%s_%s", pluginName, version)
	if client, ok := l.runningPlugins[pluginKey]; ok {
		// 调用插件销毁方法
		clientProtocol, err := client.Client()
		raw, err := clientProtocol.Dispense("plugin_store")
		if err == nil {
			pluginInstance := raw.(plugin_shared.PluginStore)
			_ = pluginInstance.Destroy(context.Background())
		}
		// 停止插件进程
		client.Kill()
		delete(l.runningPlugins, pluginKey)
	}
	return nil
}
