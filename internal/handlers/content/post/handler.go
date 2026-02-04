package post

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/shuTwT/hoshikuzu/internal/infra/logger"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	post_service "github.com/shuTwT/hoshikuzu/internal/services/content/post"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type PostHandler interface {
	ListPost(c *fiber.Ctx) error
	ListPostPage(c *fiber.Ctx) error
	CreatePost(c *fiber.Ctx) error
	UpdatePostContent(c *fiber.Ctx) error
	UpdatePostSetting(c *fiber.Ctx) error
	PublishPost(c *fiber.Ctx) error
	UnpublishPost(c *fiber.Ctx) error
	QueryPost(c *fiber.Ctx) error
	QueryPostBySlug(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
	GetSummaryForStream(c *fiber.Ctx) error
	GetPostMonthStats(c *fiber.Ctx) error
	GetRandomPost(c *fiber.Ctx) error
	SearchPosts(c *fiber.Ctx) error
}

type PostHandlerImpl struct {
	postService post_service.PostService
}

func NewPostHandlerImpl(postService post_service.PostService) *PostHandlerImpl {
	return &PostHandlerImpl{
		postService: postService,
	}
}

// @Summary 查询所有文章
// @Description 查询所有文章
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.PostResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/list [get]
func (h *PostHandlerImpl) ListPost(c *fiber.Ctx) error {
	var req model.PostListReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	posts, err := h.postService.QueryPostList(c.Context(), req)
	postResps := make([]*model.PostResp, 0, len(posts))
	for _, post := range posts {
		postResps = append(postResps, &model.PostResp{
			ID:                    post.ID,
			Title:                 post.Title,
			Slug:                  post.Slug,
			Content:               post.Content,
			MdContent:             post.MdContent,
			HtmlContent:           post.HTMLContent,
			ContentType:           string(post.ContentType),
			Status:                string(post.Status),
			IsAutogenSummary:      post.IsAutogenSummary,
			IsVisible:             post.IsVisible,
			IsPinToTop:            post.IsPinToTop,
			IsAllowComment:        post.IsAllowComment,
			IsVisibleAfterComment: post.IsVisibleAfterComment,
			IsVisibleAfterPay:     post.IsVisibleAfterPay,
			Price:                 float32(post.Price) / 100,
			PublishedAt:           post.PublishedAt,
			ViewCount:             post.ViewCount,
			CommentCount:          post.CommentCount,
			Cover:                 post.Cover,
			Keywords:              post.Keywords,
			Copyright:             post.Copyright,
			Author:                post.Author,
			Summary:               post.Summary,
			CreatedAt:             post.CreatedAt,
			Categories:            post.Edges.Categories,
			CategoryIds: func() []int {
				ids := make([]int, len(post.Edges.Categories))
				for i, cat := range post.Edges.Categories {
					ids[i] = cat.ID
				}
				return ids
			}(),
			Tags: post.Edges.Tags,
			TagIds: func() []int {
				ids := make([]int, len(post.Edges.Tags))
				for i, tag := range post.Edges.Tags {
					ids[i] = tag.ID
				}
				return ids
			}(),
		})
	}
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", postResps))
}

// @Summary 查询文章分页列表
// @Description 查询文章分页列表
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.PostResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/page [get]
func (h *PostHandlerImpl) ListPostPage(c *fiber.Ctx) error {
	var req model.PostPageReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	posts, count, err := h.postService.QueryPostPage(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	postResp := make([]*model.PostResp, 0, len(posts))
	for _, post := range posts {
		postResp = append(postResp, &model.PostResp{
			ID:                    post.ID,
			Title:                 post.Title,
			Slug:                  post.Slug,
			Content:               post.Content,
			MdContent:             post.MdContent,
			HtmlContent:           post.HTMLContent,
			ContentType:           string(post.ContentType),
			Status:                string(post.Status),
			IsAutogenSummary:      post.IsAutogenSummary,
			IsVisible:             post.IsVisible,
			IsPinToTop:            post.IsPinToTop,
			IsAllowComment:        post.IsAllowComment,
			IsVisibleAfterComment: post.IsVisibleAfterComment,
			IsVisibleAfterPay:     post.IsVisibleAfterPay,
			Price:                 float32(post.Price) / 100,
			PublishedAt:           post.PublishedAt,
			ViewCount:             post.ViewCount,
			CommentCount:          post.CommentCount,
			Cover:                 post.Cover,
			Keywords:              post.Keywords,
			Copyright:             post.Copyright,
			Author:                post.Author,
			Summary:               post.Summary,
			CreatedAt:             post.CreatedAt,
			Categories:            post.Edges.Categories,
			CategoryIds: func() []int {
				ids := make([]int, len(post.Edges.Categories))
				for i, cat := range post.Edges.Categories {
					ids[i] = cat.ID
				}
				return ids
			}(),
			Tags: post.Edges.Tags,
			TagIds: func() []int {
				ids := make([]int, len(post.Edges.Tags))
				for i, tag := range post.Edges.Tags {
					ids[i] = tag.ID
				}
				return ids
			}(),
		})
	}
	return c.JSON(model.NewSuccess("success", model.PageResult[*model.PostResp]{
		Total:   int64(count),
		Records: postResp,
	}))
}

// @Summary 创建文章
// @Description 创建一篇新文章
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param post body model.PostCreateReq true "文章创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/create [post]
func (h *PostHandlerImpl) CreatePost(c *fiber.Ctx) error {
	var post model.PostCreateReq
	if err := c.BodyParser(&post); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	newPost, err := h.postService.CreatePost(c.Context(), post.Title, post.Content)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", newPost))
}

// @Summary 更新文章内容
// @Description 更新指定文章的内容
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param post body model.PostUpdateReq true "文章更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/update/{id} [put]
func (h *PostHandlerImpl) UpdatePostContent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	var post model.PostUpdateReq
	if err = c.BodyParser(&post); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	newPost, err := h.postService.UpdatePostContent(c.Context(), id, post.Content, post.HtmlContent, post.MdContent)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", newPost))
}

// @Summary 更新文章设置
// @Description 更新指定文章的设置
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param post body model.PostUpdateReq true "文章更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/update/setting/{id} [put]
func (h *PostHandlerImpl) UpdatePostSetting(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	var post model.PostUpdateReq
	if err = c.BodyParser(&post); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	newPost, err := h.postService.UpdatePostSetting(c.Context(), id, post)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", newPost))
}

// @Summary 发布文章
// @Description 发布指定文章
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/publish/{id} [put]
func (h *PostHandlerImpl) PublishPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	newPost, err := h.postService.PublishPost(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", newPost))
}

// @Summary 取消发布文章
// @Description 取消发布指定文章
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/unpublish/{id} [put]
func (h *PostHandlerImpl) UnpublishPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	newPost, err := h.postService.UnpublishPost(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", newPost))
}

// @Summary 查询文章
// @Description 查询指定文章
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/query/{id} [get]
func (h *PostHandlerImpl) QueryPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	post, err := h.postService.QueryPostById(c.Context(), id)
	postResp := model.PostResp{
		ID:                    post.ID,
		Title:                 post.Title,
		Slug:                  post.Slug,
		Content:               post.Content,
		MdContent:             post.MdContent,
		HtmlContent:           post.HTMLContent,
		ContentType:           string(post.ContentType),
		Status:                string(post.Status),
		IsAutogenSummary:      post.IsAutogenSummary,
		IsVisible:             post.IsVisible,
		IsPinToTop:            post.IsPinToTop,
		IsAllowComment:        post.IsAllowComment,
		IsVisibleAfterComment: post.IsVisibleAfterComment,
		IsVisibleAfterPay:     post.IsVisibleAfterPay,
		Price:                 float32(post.Price) / 100,
		PublishedAt:           post.PublishedAt,
		ViewCount:             post.ViewCount,
		CommentCount:          post.CommentCount,
		Cover:                 post.Cover,
		Keywords:              post.Keywords,
		Copyright:             post.Copyright,
		Author:                post.Author,
		Summary:               post.Summary,
		Categories:            post.Edges.Categories,
		CategoryIds: func() []int {
			ids := make([]int, len(post.Edges.Categories))
			for i, cat := range post.Edges.Categories {
				ids[i] = cat.ID
			}
			return ids
		}(),
		Tags: post.Edges.Tags,
		TagIds: func() []int {
			ids := make([]int, len(post.Edges.Tags))
			for i, tag := range post.Edges.Tags {
				ids[i] = tag.ID
			}
			return ids
		}(),
		CreatedAt: post.CreatedAt,
	}
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", postResp))
}

// @Summary 根据Slug查询文章
// @Description 根据Slug查询指定文章
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Param slug path string true "文章Slug"
// @Success 200 {object} model.HttpSuccess{data=model.PostResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/posts/slug/{slug} [get]
func (h *PostHandlerImpl) QueryPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Slug is required"))
	}

	post, err := h.postService.QueryPostBySlug(c.Context(), slug)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	if post == nil {
		return c.JSON(model.NewError(fiber.StatusNotFound, "Post not found"))
	}

	postResp := model.PostResp{
		ID:                    post.ID,
		Title:                 post.Title,
		Slug:                  post.Slug,
		Content:               post.Content,
		MdContent:             post.MdContent,
		HtmlContent:           post.HTMLContent,
		ContentType:           string(post.ContentType),
		Status:                string(post.Status),
		IsAutogenSummary:      post.IsAutogenSummary,
		IsVisible:             post.IsVisible,
		IsPinToTop:            post.IsPinToTop,
		IsAllowComment:        post.IsAllowComment,
		IsVisibleAfterComment: post.IsVisibleAfterComment,
		IsVisibleAfterPay:     post.IsVisibleAfterPay,
		Price:                 float32(post.Price) / 100,
		PublishedAt:           post.PublishedAt,
		ViewCount:             post.ViewCount,
		CommentCount:          post.CommentCount,
		Cover:                 post.Cover,
		Keywords:              post.Keywords,
		Copyright:             post.Copyright,
		Author:                post.Author,
		Summary:               post.Summary,
		Categories:            post.Edges.Categories,
		CategoryIds: func() []int {
			ids := make([]int, len(post.Edges.Categories))
			for i, cat := range post.Edges.Categories {
				ids[i] = cat.ID
			}
			return ids
		}(),
		Tags: post.Edges.Tags,
		TagIds: func() []int {
			ids := make([]int, len(post.Edges.Tags))
			for i, tag := range post.Edges.Tags {
				ids[i] = tag.ID
			}
			return ids
		}(),
		CreatedAt: post.CreatedAt,
	}
	return c.JSON(model.NewSuccess("success", postResp))
}

// @Summary 删除文章
// @Description 删除指定文章
// @Tags 后台管理接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/delete/{id} [delete]
func (h *PostHandlerImpl) DeletePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	if err := h.postService.DeletePost(c.Context(), id); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// simulateAIProcessing 模拟AI处理过程
func simulateAIProcessing(targetText string, ch chan model.AIResponse) {
	defer close(ch)

	for i, text := range targetText {
		ch <- model.AIResponse{
			Content: string(text),
			Done:    i == len(targetText)-1,
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// @Summary 获取文章摘要流
// @Description 获取指定文章的摘要流
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.HttpSuccess{data=model.AIResponse}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/summary/stream/{id} [get]
func (h *PostHandlerImpl) GetSummaryForStream(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	responseChan := make(chan model.AIResponse)
	var targetStr string
	post, err := h.postService.GetRandomPost(c.Context())
	if err != nil {
		targetStr = "看来遇到了点问题，这不是你的问题"
	} else {
		if post.Summary != "" {
			targetStr = post.Summary
		} else {
			targetStr = "你好!我是AI助手,我正在处理你的请求,这是一个流式响应示例,马上就要完成了,处理完成!"
		}
	}

	go simulateAIProcessing(targetStr, responseChan)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		for response := range responseChan {
			data, err := json.Marshal(response.Content)
			if err != nil {
				continue
			}

			_, err = fmt.Fprintf(w, "%s", string(data))

			err = w.Flush()

			if err != nil {
				break
			}

			if response.Done {
				break
			}
		}

	}))

	return nil
}

// @Summary 获取文章月份统计
// @Description 获取每个月份的文章数量统计
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Param limit query int false "返回数据条数限制"
// @Success 200 {object} model.HttpSuccess{data=[]model.PostMonthStat}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/month-stats [get]
func (h *PostHandlerImpl) GetPostMonthStats(c *fiber.Ctx) error {
	var req model.PostMonthStatsReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	stats, err := h.postService.GetPostMonthStats(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	logger.Info("GetPostMonthStats: %v", stats)
	return c.JSON(model.NewSuccess("success", stats))
}

// @Summary 随机获取一篇文章
// @Description 随机获取一篇文章
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=ent.Post}
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/random [get]
func (h *PostHandlerImpl) GetRandomPost(c *fiber.Ctx) error {
	post, err := h.postService.GetRandomPost(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	if post == nil {
		return c.JSON(model.NewError(fiber.StatusNotFound, "No posts found"))
	}
	return c.JSON(model.NewSuccess("success", post))
}

// @Summary 搜索文章
// @Description 根据查询参数搜索文章
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Param query query string false "搜索查询"
// @Param limit query int false "返回数据条数限制"
// @Param offset query int false "返回数据偏移量"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.PostSearchResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/post/search [get]
func (h *PostHandlerImpl) SearchPosts(c *fiber.Ctx) error {
	var req model.PostSearchReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	results, total, err := h.postService.SearchPosts(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*model.PostSearchResp]{
		Total:   int64(total),
		Records: results,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}
