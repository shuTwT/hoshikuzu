package prompt

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/content/prompt"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type PromptHandler struct {
	promptService prompt.PromptService
}

func NewPromptHandler(promptService prompt.PromptService) *PromptHandler {
	return &PromptHandler{
		promptService: promptService,
	}
}

// @Summary 查询提示词分页
// @Description 查询提示词分页列表，包含模型图片与标签
// @Tags 后台管理接口/提示词画廊
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "标题或正文关键词"
// @Param tag_id query int false "标签ID"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Prompt]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/prompt/page [get]
func (h *PromptHandler) ListPromptPage(c *fiber.Ctx) error {
	var pageQuery model.PromptPageQuery
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if pageQuery.Page < 1 {
		pageQuery.Page = 1
	}
	if pageQuery.Size < 1 || pageQuery.Size > 100 {
		pageQuery.Size = 10
	}

	count, prompts, err := h.promptService.ListPromptPage(c.Context(), &pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	pageResult := model.PageResult[*ent.Prompt]{
		Total:   int64(count),
		Records: prompts,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建提示词
// @Description 创建提示词，同时保存模型图片与标签
// @Tags 后台管理接口/提示词画廊
// @Accept json
// @Produce json
// @Param req body model.PromptCreateReq true "提示词创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Prompt}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/prompt/create [post]
func (h *PromptHandler) CreatePrompt(c *fiber.Ctx) error {
	var req model.PromptCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	newPrompt, err := h.promptService.CreatePrompt(c.Context(), &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", newPrompt))
}

// @Summary 更新提示词
// @Description 更新提示词，模型图片与标签为全量覆盖
// @Tags 后台管理接口/提示词画廊
// @Accept json
// @Produce json
// @Param id path string true "提示词ID"
// @Param req body model.PromptUpdateReq true "提示词更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Prompt}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/prompt/update/{id} [put]
func (h *PromptHandler) UpdatePrompt(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	var req model.PromptUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedPrompt, err := h.promptService.UpdatePrompt(c.Context(), id, &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", updatedPrompt))
}

// @Summary 查询提示词
// @Description 查询指定提示词，包含模型图片与标签
// @Tags 后台管理接口/提示词画廊
// @Accept json
// @Produce json
// @Param id path string true "提示词ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Prompt}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/prompt/query/{id} [get]
func (h *PromptHandler) QueryPrompt(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	p, err := h.promptService.QueryPrompt(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", p))
}

// @Summary 删除提示词
// @Description 删除指定提示词及其模型图片
// @Tags 后台管理接口/提示词画廊
// @Accept json
// @Produce json
// @Param id path string true "提示词ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/prompt/delete/{id} [delete]
func (h *PromptHandler) DeletePrompt(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	if err := h.promptService.DeletePrompt(c.Context(), id); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 查询提示词标签列表
// @Description 查询全部提示词标签，供下拉选择
// @Tags 后台管理接口/提示词画廊
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.PromptTag}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/prompt-tag/list [get]
func (h *PromptHandler) ListPromptTag(c *fiber.Ctx) error {
	tags, err := h.promptService.ListPromptTag(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", tags))
}
