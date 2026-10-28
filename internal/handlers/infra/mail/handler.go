package mail

import (
	mail_service "github.com/shuTwT/hoshikuzu/internal/services/infra/mail"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MailHandler struct {
	mailService mail_service.MailService
}

func NewMailHandler(mailService mail_service.MailService) *MailHandler {
	return &MailHandler{
		mailService: mailService,
	}
}

// @Summary 发送测试邮件
// @Description 使用当前保存的 SMTP 配置发送一封测试邮件，用于验证邮件服务配置是否正确
// @Tags 后台管理接口/邮件
// @Accept json
// @Produce json
// @Param req body model.SendTestMailReq true "测试邮件请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/mail/test [post]
func (h *MailHandler) SendTestMail(c *fiber.Ctx) error {
	var req model.SendTestMailReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if req.To == "" {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "收件人邮箱不能为空"))
	}

	if err := h.mailService.SendTestEmail(c.Context(), req.To, req.Subject, req.Content); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}
