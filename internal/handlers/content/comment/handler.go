package comment

import (
	"strconv"
	"time"

	comment_service "github.com/shuTwT/hoshikuzu/internal/services/content/comment"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

var TWIKOO_EVENT = struct {
	GetConfig        string
	CommentGET       string
	GetCommentsCount string
	CommentSubmit    string
}{
	GetConfig:        "GET_CONFIG",
	CommentGET:       "COMMENT_GET",
	GetCommentsCount: "GET_COMMENTS_COUNT",
	CommentSubmit:    "COMMENT_SUBMIT",
}

type CommentHandler interface {
	ListCommentPage(c *fiber.Ctx) error
	HandleTwikoo(c *fiber.Ctx) error
	RecentComment(c *fiber.Ctx) error
	GetComment(c *fiber.Ctx) error
}

type CommentHandlerImpl struct {
	commentService comment_service.CommentService
}

func NewCommentHandlerImpl(commentService comment_service.CommentService) *CommentHandlerImpl {
	return &CommentHandlerImpl{commentService: commentService}
}

// @Summary 获取评论列表
// @Description 获取评论列表
// @Tags comment
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Comment]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/comment/page [get]
func (h *CommentHandlerImpl) ListCommentPage(c *fiber.Ctx) error {
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
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Comment}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/comment/{id} [get]
func (h *CommentHandlerImpl) GetComment(c *fiber.Ctx) error {
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

type TwikooReqBody struct {
	AccessToken  *string   `json:"accessToken"`
	Event        string    `json:"event"`
	Url          *string   `json:"url"`
	EnvId        *string   `json:"envId"`
	IncludeReply *bool     `json:"includeReply"`
	Region       *string   `json:"region"`
	Urls         *[]string `json:"urls"`
	Comment      *string   `json:"comment"`
	Href         *string   `json:"href"`
	Link         *string   `json:"link"`
	Mail         *string   `json:"mail"`
	Nick         *string   `json:"nick"`
	UA           *string   `json:"ua"`
}

// @Summary 处理Twikoo评论事件
// @Description 处理Twikoo评论事件，包括获取配置、获取评论、获取评论数量、提交评论
// @Tags comment
// @Accept json
// @Produce json
// @Param twikoo_req_body body TwikooReqBody true "Twikoo请求体"
// @Success 200 {object} model.HttpSuccess{data=interface{}}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/twikoo [get]
// @Router /api/v1/twikoo [post]
// @Router /api/v1/twikoo [put]
// @Router /api/v1/twikoo [delete]
func (h *CommentHandlerImpl) HandleTwikoo(c *fiber.Ctx) error {

	var reqBody TwikooReqBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// 配置
	if reqBody.Event == TWIKOO_EVENT.GetConfig {
		return c.JSON(fiber.Map{
			"code":   0,
			"config": fiber.Map{},
		})
	}
	// 获取评论
	if reqBody.Event == TWIKOO_EVENT.CommentGET {
		commmentList := []fiber.Map{}
		comments, err := h.commentService.ListComment(c.Context(), *reqBody.Url)
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		for _, comment := range comments {
			browser, os := h.commentService.ParseUserAgent(*comment.UserAgent)
			commmentList = append(commmentList, fiber.Map{
				"id":        comment.ID,
				"parentId":  comment.ParentID,
				"url":       comment.URL,
				"userAgent": comment.UserAgent,
				"comment":   comment.Content,
				"browser":   browser,
				"os":        os,
				"created":   time.Time(comment.CreatedAt).UnixMilli(),
				"updated":   time.Time(comment.UpdatedAt).UnixMilli(),
				"replies":   []fiber.Map{},
			})
		}
		return c.JSON(fiber.Map{
			"count": len(commmentList),
			"data":  commmentList,
			"more":  false,
		})
	}
	// 获取评论数量
	if reqBody.Event == TWIKOO_EVENT.GetCommentsCount {
		data := []fiber.Map{}
		count, err := h.commentService.CountComment(c.Context(), *reqBody.IncludeReply, *reqBody.Urls)
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		for _, url := range *reqBody.Urls {
			data = append(data, fiber.Map{
				"url":   url,
				"count": count,
			})
		}
		return c.JSON(fiber.Map{
			"data": data,
		})
	}

	// 提交评论
	if reqBody.Event == TWIKOO_EVENT.CommentSubmit {
		ipAddress := c.IP()
		id, err := h.commentService.CreateComment(c.Context(), *reqBody.Comment, *reqBody.Href, *reqBody.Link, *reqBody.Mail, *reqBody.Nick, *reqBody.UA, *reqBody.Url, ipAddress)
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		return c.JSON(fiber.Map{
			"id": id,
		})
	}
	return c.JSON(fiber.Map{})
}

// @Summary 获取最近评论
// @Description 获取最近评论
// @Tags comment
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]interface{}}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/comment/recent [get]
func (h *CommentHandlerImpl) RecentComment(c *fiber.Ctx) error {
	comments, err := h.commentService.GetRecentComment(c.Context(), 10)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	commentList := []fiber.Map{}
	for _, entity := range comments {
		commentList = append(commentList, fiber.Map{
			"id":        entity.ID,
			"parentId":  entity.ParentID,
			"url":       entity.URL,
			"userAgent": entity.UserAgent,
			"comment":   entity.Content,
			"created":   time.Time(entity.CreatedAt).UnixMilli(),
			"updated":   time.Time(entity.UpdatedAt).UnixMilli(),
			"replies":   []fiber.Map{},
		})
	}
	return c.JSON(model.NewSuccess("最近评论获取成功", commentList))
}
