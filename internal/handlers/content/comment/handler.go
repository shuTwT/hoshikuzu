package comment

import (
	"strconv"

	comment_service "github.com/shuTwT/hoshikuzu/internal/services/content/comment"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	commentService comment_service.CommentService
}

func NewCommentHandler(commentService comment_service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// @Summary 获取评论列表
// @Description 获取评论列表
// @Tags 后台管理接口/评论
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Comment]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/comment/page [get]
func (h *CommentHandler) ListCommentPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	resp, err := h.commentService.ListCommentPage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("评论列表获取成功", resp))
}

// @Summary 获取评论
// @Description 获取指定评论
// @Tags 后台管理接口/评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Comment}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/comment/query/{id} [get]
func (h *CommentHandler) GetComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp, err := h.commentService.GetComment(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("评论列表获取成功", resp))
}
