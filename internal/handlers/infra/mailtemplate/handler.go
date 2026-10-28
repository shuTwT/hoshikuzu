package mailtemplate

import (
	"errors"
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	mailtemplate_service "github.com/shuTwT/hoshikuzu/internal/services/infra/mailtemplate"

	"github.com/gofiber/fiber/v2"
)

type MailTemplateHandler struct {
	mailTemplateService mailtemplate_service.MailTemplateService
}

func NewMailTemplateHandler(mailTemplateService mailtemplate_service.MailTemplateService) *MailTemplateHandler {
	return &MailTemplateHandler{
		mailTemplateService: mailTemplateService,
	}
}

// @Summary 查询邮件模板
// @Description 查询指定ID的邮件模板详情
// @Tags 后台管理接口/邮件模板
// @Accept json
// @Produce json
// @Param id path int true "邮件模板ID"
// @Success 200 {object} model.HttpSuccess{data=ent.MailTemplate}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/mail-template/query/{id} [get]
func (h *MailTemplateHandler) QueryMailTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	t, err := h.mailTemplateService.QueryMailTemplate(c.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrMailTemplateNotFound) {
			return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", t))
}

// @Summary 分页查询邮件模板
// @Description 分页查询邮件模板列表，支持按模板名称、模板编码模糊搜索
// @Tags 后台管理接口/邮件模板
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param name query string false "模板名称（模糊匹配）"
// @Param code query string false "模板编码（模糊匹配）"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.MailTemplate]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/mail-template/page [get]
func (h *MailTemplateHandler) QueryMailTemplatePage(c *fiber.Ctx) error {
	var pageQuery model.MailTemplatePageQuery
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if pageQuery.Page < 1 {
		pageQuery.Page = 1
	}
	if pageQuery.Size < 1 {
		pageQuery.Size = 10
	}

	count, templates, err := h.mailTemplateService.QueryMailTemplatePage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.MailTemplate]{
		Total:   int64(count),
		Records: templates,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建邮件模板
// @Description 创建新的邮件模板，模板编码全局唯一
// @Tags 后台管理接口/邮件模板
// @Accept json
// @Produce json
// @Param request body model.MailTemplateCreateReq true "邮件模板信息"
// @Success 200 {object} model.HttpSuccess{data=ent.MailTemplate}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/mail-template/create [post]
func (h *MailTemplateHandler) CreateMailTemplate(c *fiber.Ctx) error {
	var createReq model.MailTemplateCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if createReq.Name == "" || createReq.Code == "" || createReq.Subject == "" || createReq.Content == "" {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "模板名称、模板编码、邮件主题、邮件正文不能为空"))
	}

	t, err := h.mailTemplateService.CreateMailTemplate(c.Context(), createReq)
	if err != nil {
		if errors.Is(err, model.ErrMailTemplateCodeExists) {
			return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", t))
}

// @Summary 更新邮件模板
// @Description 更新指定ID的邮件模板信息
// @Tags 后台管理接口/邮件模板
// @Accept json
// @Produce json
// @Param id path int true "邮件模板ID"
// @Param request body model.MailTemplateUpdateReq true "邮件模板信息"
// @Success 200 {object} model.HttpSuccess{data=ent.MailTemplate}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/mail-template/update/{id} [put]
func (h *MailTemplateHandler) UpdateMailTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var updateReq model.MailTemplateUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	t, err := h.mailTemplateService.UpdateMailTemplate(c.Context(), id, updateReq)
	if err != nil {
		if errors.Is(err, model.ErrMailTemplateNotFound) {
			return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
		}
		if errors.Is(err, model.ErrMailTemplateCodeExists) {
			return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", t))
}

// @Summary 删除邮件模板
// @Description 删除指定ID的邮件模板
// @Tags 后台管理接口/邮件模板
// @Accept json
// @Produce json
// @Param id path int true "邮件模板ID"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/mail-template/delete/{id} [delete]
func (h *MailTemplateHandler) DeleteMailTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	err = h.mailTemplateService.DeleteMailTemplate(c.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrMailTemplateNotFound) {
			return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
