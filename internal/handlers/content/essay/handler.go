package essay

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/internal/middleware"
	"github.com/shuTwT/hoshikuzu/internal/services/content/essay"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type EssayHandler interface {
	CreateEssay(c *fiber.Ctx) error
	UpdateEssay(c *fiber.Ctx) error
	GetEssayPage(c *fiber.Ctx) error
	ListEssay(c *fiber.Ctx) error
	DeleteEssay(c *fiber.Ctx) error
}

type EssayHandlerImpl struct {
	service essay.EssayService
}

func NewEssayHandlerImpl(service essay.EssayService) EssayHandler {
	return &EssayHandlerImpl{service: service}
}

// @Summary 创建说说
// @Description 创建一个新说说
// @Tags 后台管理接口/说说
// @Accept json
// @Produce json
// @Param req body model.EssayCreateReq true "说创建请求"
// @Success 200 {object} model.HttpSuccess{data=model.EssayCreateReq}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/essay/create [post]
func (h *EssayHandlerImpl) CreateEssay(c *fiber.Ctx) error {
	createReq := &model.EssayCreateReq{}
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	userId := middleware.GetCurrentUser(c).ID
	_, err := h.service.CreateEssay(c.Context(), userId, createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 更新说说
// @Description 更新指定说说的信息
// @Tags 后台管理接口/说说
// @Accept json
// @Produce json
// @Param id path string true "说说ID"
// @Param req body model.EssayUpdateReq true "说说更新请求"
// @Success 200 {object} model.HttpSuccess{data=model.EssayUpdateReq}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/essay/update/{id} [put]
func (h *EssayHandlerImpl) UpdateEssay(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	updateReq := &model.EssayUpdateReq{}
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.UpdateEssay(c.Context(), id, updateReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 获取说说分页列表
// @Description 获取说说分页列表
// @Tags 后台管理接口/说说
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.EssayResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/essay/page [get]
func (h *EssayHandlerImpl) GetEssayPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	essays, total, err := h.service.GetEssayPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.EssayResp, 0, len(essays))
	for _, e := range essays {
		resp = append(resp, model.EssayResp{
			ID:       e.ID,
			Content:  e.Content,
			Draft:    e.Draft,
			Images:   e.Images,
			CreateAt: (*model.LocalTime)(&e.CreatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", model.PageResult[model.EssayResp]{
		Total:   int64(total),
		Records: resp,
	}))
}

// @Summary 获取说说列表
// @Description 获取说说列表
// @Tags 后台管理接口/说说
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} model.HttpSuccess{data=[]model.EssayResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/essay/list [get]
func (h *EssayHandlerImpl) ListEssay(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	essays, err := h.service.GetEssayList(c.Context(), limit)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.EssayResp, 0, len(essays))
	for _, e := range essays {
		resp = append(resp, model.EssayResp{
			ID:       e.ID,
			Content:  e.Content,
			Draft:    e.Draft,
			Images:   e.Images,
			CreateAt: (*model.LocalTime)(&e.CreatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 删除说说
// @Description 删除指定说说
// @Tags 后台管理接口/说说
// @Accept json
// @Produce json
// @Param id path string true "说说ID"
// @Success 200 {object} model.HttpSuccess{data=model.EssayResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/essay/delete/{id} [delete]
func (h *EssayHandlerImpl) DeleteEssay(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.DeleteEssay(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}
