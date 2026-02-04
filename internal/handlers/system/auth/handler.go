package auth_handler

import (
	"github.com/shuTwT/hoshikuzu/internal/services/system/auth"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type AuthHandlerImpl struct {
	authService auth.AuthService
}

func NewAuthHandlerImpl(authService auth.AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{authService: authService}
}

// @Summary 用户登录
// @Description 验证用户凭据并返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "登录请求"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 401 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/auth/login/password [post]
func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	var req *model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	loginResp, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(
			fiber.StatusUnauthorized,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("Login successful", loginResp))
}
