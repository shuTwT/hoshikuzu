package knowledgebase

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/internal/services/content/knowledgebase"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type KnowledgeBaseHandler interface {
	CreateKnowledgeBase(c *fiber.Ctx) error
	UpdateKnowledgeBase(c *fiber.Ctx) error
	GetKnowledgeBasePage(c *fiber.Ctx) error
	GetKnowledgeBase(c *fiber.Ctx) error
	DeleteKnowledgeBase(c *fiber.Ctx) error
	GetKnowledgeBaseList(c *fiber.Ctx) error
}

type KnowledgeBaseHandlerImpl struct {
	service knowledgebase.KnowledgeBaseService
}

func NewKnowledgeBaseHandlerImpl(service knowledgebase.KnowledgeBaseService) KnowledgeBaseHandler {
	return &KnowledgeBaseHandlerImpl{service: service}
}

// @Summary 创建知识库
// @Description 创建知识库
// @Tags 后台管理接口/知识库
// @Accept json
// @Produce json
// @Param knowledge_base_create_req body model.KnowledgeBaseCreateReq true "知识库创建请求体"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/knowledge-base/create [post]
func (h *KnowledgeBaseHandlerImpl) CreateKnowledgeBase(c *fiber.Ctx) error {
	createReq := model.KnowledgeBaseCreateReq{}
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	_, err := h.service.CreateKnowledgeBase(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 更新知识库
// @Description 更新知识库
// @Tags 后台管理接口/知识库
// @Accept json
// @Produce json
// @Param id path int true "知识库 ID"
// @Param knowledge_base_update_req body model.KnowledgeBaseUpdateReq true "知识库更新请求体"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/knowledge-base/update/{id} [put]
func (h *KnowledgeBaseHandlerImpl) UpdateKnowledgeBase(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	updateReq := model.KnowledgeBaseUpdateReq{}
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.UpdateKnowledgeBase(c.Context(), id, updateReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 获取知识库分页列表
// @Description 获取知识库分页列表
// @Tags 后台管理接口/知识库
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Param page query int false "页码" default(1)
// @Success 200 {object} model.HttpSuccess{data=[]model.KnowledgeBaseResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/knowledge-base/page [get]
func (h *KnowledgeBaseHandlerImpl) GetKnowledgeBasePage(c *fiber.Ctx) error {
	queryReq := model.KnowledgeBaseQueryReq{}
	if err := c.QueryParser(&queryReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	kbs, total, err := h.service.GetKnowledgeBasePage(c.Context(), queryReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.KnowledgeBaseResp, 0, len(kbs))
	for _, kb := range kbs {
		resp = append(resp, model.KnowledgeBaseResp{
			ID:                    kb.ID,
			Name:                  kb.Name,
			ModelProvider:         kb.ModelProvider.String(),
			Model:                 kb.Model,
			VectorDimension:       kb.VectorDimension,
			MaxBatchDocumentCount: kb.MaxBatchDocumentCount,
			CreatedAt:             (*model.LocalTime)(&kb.CreatedAt),
			UpdatedAt:             (*model.LocalTime)(&kb.UpdatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", model.PageResult[model.KnowledgeBaseResp]{
		Total:   int64(total),
		Records: resp,
	}))
}

// @Summary 获取知识库
// @Description 获取指定知识库
// @Tags 后台管理接口/知识库
// @Accept json
// @Produce json
// @Param id path int true "知识库 ID"
// @Success 200 {object} model.HttpSuccess{data=model.KnowledgeBaseResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/knowledge-base/{id} [get]
func (h *KnowledgeBaseHandlerImpl) GetKnowledgeBase(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	kb, err := h.service.GetKnowledgeBase(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := model.KnowledgeBaseResp{
		ID:                    kb.ID,
		Name:                  kb.Name,
		ModelProvider:         kb.ModelProvider.String(),
		Model:                 kb.Model,
		VectorDimension:       kb.VectorDimension,
		MaxBatchDocumentCount: kb.MaxBatchDocumentCount,
		CreatedAt:             (*model.LocalTime)(&kb.CreatedAt),
		UpdatedAt:             (*model.LocalTime)(&kb.UpdatedAt),
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 删除知识库
// @Description 删除指定知识库
// @Tags 后台管理接口/知识库
// @Accept json
// @Produce json
// @Param id path int true "知识库 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/knowledge-base/delete/{id} [delete]
func (h *KnowledgeBaseHandlerImpl) DeleteKnowledgeBase(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.DeleteKnowledgeBase(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 获取知识库列表
// @Description 获取所有知识库
// @Tags 后台管理接口/知识库
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.KnowledgeBaseResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/knowledge-base/list [get]
func (h *KnowledgeBaseHandlerImpl) GetKnowledgeBaseList(c *fiber.Ctx) error {
	kbs, err := h.service.GetKnowledgeBaseList(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.KnowledgeBaseResp, 0, len(kbs))
	for _, kb := range kbs {
		resp = append(resp, model.KnowledgeBaseResp{
			ID:                    kb.ID,
			Name:                  kb.Name,
			ModelProvider:         kb.ModelProvider.String(),
			Model:                 kb.Model,
			VectorDimension:       kb.VectorDimension,
			MaxBatchDocumentCount: kb.MaxBatchDocumentCount,
			CreatedAt:             (*model.LocalTime)(&kb.CreatedAt),
			UpdatedAt:             (*model.LocalTime)(&kb.UpdatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", resp))
}
