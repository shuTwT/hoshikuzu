package initialize

import (
	"fmt"

	"github.com/shuTwT/hoshikuzu/ent"
	setting_service "github.com/shuTwT/hoshikuzu/internal/services/system/setting"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type InitializeHandler interface {
	PreInit(c *fiber.Ctx) error
	Initialize(c *fiber.Ctx) error
}

type InitializeHandlerImpl struct {
	client         *ent.Client
	userService    user_service.UserService
	settingService setting_service.SettingService
}

func NewInitializeHandlerImpl(client *ent.Client, userService user_service.UserService, settingService setting_service.SettingService) *InitializeHandlerImpl {
	return &InitializeHandlerImpl{client: client, userService: userService, settingService: settingService}
}

func (h *InitializeHandlerImpl) PreInit(c *fiber.Ctx) error {
	DBType := config.GetString(config.DATABASE_TYPE)
	return c.JSON(model.NewSuccess("success", model.PreInitResp{DBType: DBType}))
}

// Initialize 处理系统初始化请求
// @Summary 系统初始化
// @Description 首次运行系统时进行初始化设置，包括数据库配置和创建管理员账户
// @Tags Initialize
// @Accept json
// @Produce json
// @Param request body model.InitializeRequest true "初始化请求"
// @Success 200 {object} map[string]string "初始化成功"
// @Failure 400 {object} model.HttpError "请求参数错误"
// @Failure 409 {object} model.HttpError "系统已初始化"
// @Failure 500 {object} model.HttpError "服务器内部错误"
// @Router /api/initialize [post]
func (h *InitializeHandlerImpl) Initialize(c *fiber.Ctx) error {
	var req *model.InitializeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if req.AdminPassword != req.ConfirmPassword {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "两次输入的密码不一致"))
	}

	// 检查系统是否已初始化
	isInitialized, err := h.settingService.IsSystemInitialized(c.UserContext())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, fmt.Sprintf("检查系统初始化状态失败: %v", err)))
	}
	if isInitialized {
		return c.JSON(model.NewError(fiber.StatusConflict, "系统已初始化，请勿重复操作"))
	}

	// 创建角色e
	h.initRole(c)
	// 创建超级管理员账户
	adminUser, err := h.userService.CreateUser(c.UserContext(), model.UserCreateReq{
		Name:     req.AdminUsername,
		Password: req.AdminPassword,
		Email:    req.AdminEmail,
		RoleID:   1,
	})
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, fmt.Sprintf("创建超级管理员账户失败: %v", err)))
	}
	fmt.Printf("超级管理员账户创建成功: %s\n", adminUser.Name)

	// 标记系统已初始化
	if err := h.settingService.SetSystemInitialized(c.UserContext()); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, fmt.Sprintf("标记系统初始化状态失败: %v", err)))
	}

	return c.JSON(model.NewSuccess("系统初始化成功", nil))
}

func (h *InitializeHandlerImpl) initRole(c *fiber.Ctx) *ent.Role {
	// 初始化角色
	role := h.client.Role.Create().
		SetID(1).
		SetName("超级管理员").
		SetCode("superAdmin").
		SetIsDefault(true).
		SaveX(c.UserContext())
	h.client.Role.Create().
		SetID(2).
		SetName("访客").
		SetCode("guest").
		SetIsDefault(true).
		SaveX(c.UserContext())
	h.client.Role.Create().
		SetID(3).
		SetName("普通用户").
		SetCode("common").
		SetIsDefault(true).
		SaveX(c.UserContext())

	return role
}
