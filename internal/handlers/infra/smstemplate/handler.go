package smstemplate

import (
	"errors"
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	smstemplate_service "github.com/shuTwT/hoshikuzu/internal/services/infra/smstemplate"

	"github.com/gofiber/fiber/v2"
)

type SmsTemplateHandler struct {
	smsTemplateService smstemplate_service.SmsTemplateService
}

func NewSmsTemplateHandler(smsTemplateService smstemplate_service.SmsTemplateService) *SmsTemplateHandler {
	return &SmsTemplateHandler{
		smsTemplateService: smsTemplateService,
	}
}

// @Summary 查询短信模板
// @Description 查询指定ID的短信模板详情
// @Tags 后台管理接口/短信模板
// @Accept json
// @Produce json
// @Param id path int true "短信模板ID"
// @Success 200 {object} model.HttpSuccess{data=ent.SmsTemplate}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/sms-template/query/{id} [get]
func (h *SmsTemplateHandler) QuerySmsTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	t, err := h.smsTemplateService.QuerySmsTemplate(c.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrSmsTemplateNotFound) {
			return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", t))
}

// @Summary 分页查询短信模板
// @Description 分页查询短信模板列表，支持按模板名称、模板编码模糊搜索
// @Tags 后台管理接口/短信模板
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param name query string false "模板名称（模糊匹配）"
// @Param code query string false "模板编码（模糊匹配）"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.SmsTemplate]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/sms-template/page [get]
func (h *SmsTemplateHandler) QuerySmsTemplatePage(c *fiber.Ctx) error {
	var pageQuery model.SmsTemplatePageQuery
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if pageQuery.Page < 1 {
		pageQuery.Page = 1
	}
	if pageQuery.Size < 1 {
		pageQuery.Size = 10
	}

	count, templates, err := h.smsTemplateService.QuerySmsTemplatePage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.SmsTemplate]{
		Total:   int64(count),
		Records: templates,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建短信模板
// @Description 创建新的短信模板，模板编码全局唯一
// @Tags 后台管理接口/短信模板
// @Accept json
// @Produce json
// @Param request body model.SmsTemplateCreateReq true "短信模板信息"
// @Success 200 {object} model.HttpSuccess{data=ent.SmsTemplate}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/sms-template/create [post]
func (h *SmsTemplateHandler) CreateSmsTemplate(c *fiber.Ctx) error {
	var createReq model.SmsTemplateCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if createReq.Name == "" || createReq.Code == "" || createReq.Content == "" {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "模板名称、模板编码、模板内容不能为空"))
	}

	t, err := h.smsTemplateService.CreateSmsTemplate(c.Context(), createReq)
	if err != nil {
		if errors.Is(err, model.ErrSmsTemplateCodeExists) {
			return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", t))
}

// @Summary 更新短信模板
// @Description 更新指定ID的短信模板信息
// @Tags 后台管理接口/短信模板
// @Accept json
// @Produce json
// @Param id path int true "短信模板ID"
// @Param request body model.SmsTemplateUpdateReq true "短信模板信息"
// @Success 200 {object} model.HttpSuccess{data=ent.SmsTemplate}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/sms-template/update/{id} [put]
func (h *SmsTemplateHandler) UpdateSmsTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var updateReq model.SmsTemplateUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	t, err := h.smsTemplateService.UpdateSmsTemplate(c.Context(), id, updateReq)
	if err != nil {
		if errors.Is(err, model.ErrSmsTemplateNotFound) {
			return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
		}
		if errors.Is(err, model.ErrSmsTemplateCodeExists) {
			return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", t))
}

// @Summary 删除短信模板
// @Description 删除指定ID的短信模板
// @Tags 后台管理接口/短信模板
// @Accept json
// @Produce json
// @Param id path int true "短信模板ID"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/sms-template/delete/{id} [delete]
func (h *SmsTemplateHandler) DeleteSmsTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	err = h.smsTemplateService.DeleteSmsTemplate(c.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrSmsTemplateNotFound) {
			return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
