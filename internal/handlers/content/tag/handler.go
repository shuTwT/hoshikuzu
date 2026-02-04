package tag

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	tag_service "github.com/shuTwT/hoshikuzu/internal/services/content/tag"

	"github.com/gofiber/fiber/v2"
)

type TagHandler interface {
	QueryTag(c *fiber.Ctx) error
	QueryTagList(c *fiber.Ctx) error
	QueryTagPage(c *fiber.Ctx) error
	CreateTag(c *fiber.Ctx) error
	UpdateTag(c *fiber.Ctx) error
	DeleteTag(c *fiber.Ctx) error
}

type TagHandlerImpl struct {
	tagService tag_service.TagService
}

func NewTagHandlerImpl(tagService tag_service.TagService) *TagHandlerImpl {
	return &TagHandlerImpl{
		tagService: tagService,
	}
}

// @Summary 查询标签
// @Description 查询指定标签
// @Tags 后台管理接口/标签
// @Accept json
// @Produce json
// @Param id path int true "标签 ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Tag}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/tag/{id} [get]
func (h *TagHandlerImpl) QueryTag(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	tag, err := h.tagService.QueryTag(c, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Tag not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", tag))
}

// @Summary 查询标签列表
// @Description 查询所有标签
// @Tags 后台管理接口/标签
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Tag}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/tag/list [get]
func (h *TagHandlerImpl) QueryTagList(c *fiber.Ctx) error {
	tags, err := h.tagService.QueryTagList(c)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", tags))
}

// @Summary 查询标签分页
// @Description 查询标签分页列表
// @Tags 后台管理接口/标签
// @Accept json
// @Produce json
// @Param limit query int false "返回数据条数限制"
// @Param offset query int false "返回数据偏移量"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Tag]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/tag/page [get]
func (h *TagHandlerImpl) QueryTagPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, tags, err := h.tagService.QueryTagPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	pageResult := model.PageResult[*ent.Tag]{
		Total:   int64(count),
		Records: tags,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建标签
// @Description 创建新标签
// @Tags 后台管理接口/标签
// @Accept json
// @Produce json
// @Param tag body model.TagCreateReq true "标签创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Tag}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/tag/create [post]
func (h *TagHandlerImpl) CreateTag(c *fiber.Ctx) error {
	var createReq model.TagCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	tag, err := h.tagService.CreateTag(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", tag))
}

// @Summary 更新标签
// @Description 更新指定标签
// @Tags 后台管理接口/标签
// @Accept json
// @Produce json
// @Param id path int true "标签 ID"
// @Param tag body model.TagUpdateReq true "标签更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Tag}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/tag//update/{id} [put]
func (h *TagHandlerImpl) UpdateTag(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var updateReq model.TagUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	updatedTag, err := h.tagService.UpdateTag(c, id, updateReq)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Tag not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", updatedTag))
}

// @Summary 删除标签
// @Description 删除指定标签
// @Tags 后台管理接口/标签
// @Accept json
// @Produce json
// @Param id path int true "标签 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/tag/delete/{id} [delete]
func (h *TagHandlerImpl) DeleteTag(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.tagService.DeleteTag(c, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Tag not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
