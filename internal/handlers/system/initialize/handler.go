package initialize

import (
	"errors"
	"log/slog"

	initialize_service "github.com/shuTwT/hoshikuzu/internal/services/system/initialize"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type InitializeHandler struct {
	initializeService initialize_service.InitializeService
}

func NewInitializeHandler(initializeService initialize_service.InitializeService) *InitializeHandler {
	return &InitializeHandler{initializeService: initializeService}
}

// PreInit 查询初始化前置信息
// @Summary 查询初始化前置信息
// @Description 查询系统初始化所需的数据库类型等前置信息
// @Tags 公开接口/初始化
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=model.PreInitResp}
// @Failure 500 {object} model.HttpError
// @Router /api/preinit [get]
func (h *InitializeHandler) PreInit(c *fiber.Ctx) error {
	DBType := config.GetString(config.DATABASE_TYPE)
	return c.JSON(model.NewSuccess("success", model.PreInitResp{DBType: DBType}))
}

// Initialize 处理系统初始化请求
// @Summary 系统初始化
// @Description 首次运行系统时进行初始化设置，包括数据库配置和创建管理员账户
// @Tags 公开接口/初始化
// @Accept json
// @Produce json
// @Param request body model.InitializeRequest true "初始化请求"
// @Success 200 {object} map[string]string "初始化成功"
// @Failure 400 {object} model.HttpError "请求参数错误"
// @Failure 409 {object} model.HttpError "系统已初始化"
// @Failure 500 {object} model.HttpError "服务器内部错误"
// @Router /api/initialize [post]
func (h *InitializeHandler) Initialize(c *fiber.Ctx) error {
	var req model.InitializeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if req.AdminPassword != req.ConfirmPassword {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "两次输入的密码不一致"))
	}

	if err := h.initializeService.Initialize(c.UserContext(), &req); err != nil {
		if errors.Is(err, model.ErrSystemInitialized) {
			return c.JSON(model.NewError(fiber.StatusConflict, err.Error()))
		}
		slog.Error("系统初始化失败", "error", err)
		return c.JSON(model.NewError(fiber.StatusInternalServerError, "系统初始化失败，请查看服务端日志"))
	}

	return c.JSON(model.NewSuccess("系统初始化成功", nil))
}
