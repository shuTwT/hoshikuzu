package plugin

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/plugin"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type PluginHandler interface {
	CreatePlugin(c *fiber.Ctx) error
	ListPluginPage(c *fiber.Ctx) error
	QueryPlugin(c *fiber.Ctx) error
	DeletePlugin(c *fiber.Ctx) error
	StartPlugin(c *fiber.Ctx) error
	StopPlugin(c *fiber.Ctx) error
	RestartPlugin(c *fiber.Ctx) error
	RegisterPlugin(c *fiber.Ctx) error
	HeartbeatPlugin(c *fiber.Ctx) error
}

type PluginHandlerImpl struct {
	pluginService plugin.PluginService
}

func NewPluginHandlerImpl(pluginService plugin.PluginService) *PluginHandlerImpl {
	return &PluginHandlerImpl{pluginService: pluginService}
}

func (h *PluginHandlerImpl) CreatePlugin(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		slog.Error("Failed to get uploaded file", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "获取上传文件失败"))
	}

	pluginEntity, err := h.pluginService.CreatePlugin(c.Context(), file)
	if err != nil {
		slog.Error("Failed to create plugin", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildPluginResp(pluginEntity)
	slog.Info("Plugin created successfully", "plugin_id", pluginEntity.ID, "plugin_key", pluginEntity.Key)
	return c.JSON(model.NewSuccess("插件创建成功", resp))
}

func (h *PluginHandlerImpl) ListPluginPage(c *fiber.Ctx) error {
	var pageQuery model.PageQuery
	if err := c.QueryParser(&pageQuery); err != nil {
		slog.Error("Failed to parse query parameters", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "查询参数解析失败"))
	}

	count, plugins, err := h.pluginService.ListPluginPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		slog.Error("Failed to list plugins", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pluginResps := make([]*model.PluginResp, 0, len(plugins))
	for _, p := range plugins {
		pluginResps = append(pluginResps, h.buildPluginResp(p))
	}

	pageResult := model.PageResult[*model.PluginResp]{
		Total:   int64(count),
		Records: pluginResps,
	}
	return c.JSON(model.NewSuccess("插件列表获取成功", pageResult))
}

func (h *PluginHandlerImpl) QueryPlugin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid plugin ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的插件ID"))
	}

	pluginEntity, err := h.pluginService.QueryPlugin(c.Context(), id)
	if err != nil {
		slog.Error("Failed to query plugin", "plugin_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildPluginResp(pluginEntity)
	return c.JSON(model.NewSuccess("插件查询成功", resp))
}

func (h *PluginHandlerImpl) DeletePlugin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid plugin ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的插件ID"))
	}

	err = h.pluginService.DeletePlugin(c.Context(), id)
	if err != nil {
		slog.Error("Failed to delete plugin", "plugin_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin deleted successfully", "plugin_id", id)
	return c.JSON(model.NewSuccess("插件删除成功", nil))
}

func (h *PluginHandlerImpl) StartPlugin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid plugin ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的插件ID"))
	}

	err = h.pluginService.StartPlugin(c.Context(), id)
	if err != nil {
		slog.Error("Failed to start plugin", "plugin_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin started successfully", "plugin_id", id)
	return c.JSON(model.NewSuccess("插件启动成功", nil))
}

func (h *PluginHandlerImpl) StopPlugin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid plugin ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的插件ID"))
	}

	err = h.pluginService.StopPlugin(c.Context(), id)
	if err != nil {
		slog.Error("Failed to stop plugin", "plugin_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin stopped successfully", "plugin_id", id)
	return c.JSON(model.NewSuccess("插件停止成功", nil))
}

func (h *PluginHandlerImpl) RestartPlugin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid plugin ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的插件ID"))
	}

	err = h.pluginService.RestartPlugin(c.Context(), id)
	if err != nil {
		slog.Error("Failed to restart plugin", "plugin_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin restarted successfully", "plugin_id", id)
	return c.JSON(model.NewSuccess("插件重启成功", nil))
}

func (h *PluginHandlerImpl) RegisterPlugin(c *fiber.Ctx) error {
	// 检查debug模式是否开启
	if !config.GetBool(config.SERVER_DEBUG) {
		slog.Warn("RegisterPlugin called but debug mode is not enabled")
		return c.JSON(model.NewError(fiber.StatusForbidden, "此接口仅在debug模式下可用"))
	}

	// 接收插件注册信息
	var pluginInfo model.PluginRegisterReq
	if err := c.BodyParser(&pluginInfo); err != nil {
		slog.Error("Failed to parse plugin registration info", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "解析插件注册信息失败"))
	}

	// 调用服务层方法存储插件注册信息
	err := h.pluginService.RegisterPlugin(c.Context(), &pluginInfo)
	if err != nil {
		slog.Error("Failed to register plugin", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin registered successfully", "plugin_name", pluginInfo.Name)
	return c.JSON(model.NewSuccess("插件注册成功", nil))
}

func (h *PluginHandlerImpl) HeartbeatPlugin(c *fiber.Ctx) error {
	// 检查debug模式是否开启
	if !config.GetBool(config.SERVER_DEBUG) {
		slog.Warn("HeartbeatPlugin called but debug mode is not enabled")
		return c.JSON(model.NewError(fiber.StatusForbidden, "此接口仅在debug模式下可用"))
	}

	// 接收插件心跳信息
	var heartbeatInfo model.PluginHeartbeatReq
	if err := c.BodyParser(&heartbeatInfo); err != nil {
		slog.Error("Failed to parse plugin heartbeat info", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "解析插件心跳信息失败"))
	}

	// 调用服务层方法更新插件的心跳时间
	err := h.pluginService.HeartbeatPlugin(c.Context(), &heartbeatInfo)
	if err != nil {
		slog.Error("Failed to update plugin heartbeat", "plugin_name", heartbeatInfo.Name, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin heartbeat updated successfully", "plugin_name", heartbeatInfo.Name)
	return c.JSON(model.NewSuccess("插件心跳更新成功", nil))
}

func (h *PluginHandlerImpl) buildPluginResp(p *ent.Plugin) *model.PluginResp {
	var lastStartedAt, lastStoppedAt *time.Time
	if !p.LastStartedAt.IsZero() {
		lastStartedAt = &p.LastStartedAt
	}
	if !p.LastStoppedAt.IsZero() {
		lastStoppedAt = &p.LastStoppedAt
	}

	return &model.PluginResp{
		ID:               p.ID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
		Key:              p.Key,
		Name:             p.Name,
		Version:          p.Version,
		Description:      p.Description,
		BinPath:          p.BinPath,
		ProtocolVersion:  fmt.Sprintf("%d", p.ProtocolVersion),
		MagicCookieKey:   p.MagicCookieKey,
		MagicCookieValue: p.MagicCookieValue,
		Dependencies:     p.Dependencies,
		Config:           p.Config,
		Enabled:          p.Enabled,
		AutoStart:        p.AutoStart,
		Status:           string(p.Status),
		LastError:        p.LastError,
		LastStartedAt:    lastStartedAt,
		LastStoppedAt:    lastStoppedAt,
	}
}
