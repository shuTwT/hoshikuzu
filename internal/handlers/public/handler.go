package public

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/aws/smithy-go/ptr"
	"github.com/gofiber/fiber/v2"
	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/infra/logger"
	"github.com/shuTwT/hoshikuzu/internal/services/content/album"
	"github.com/shuTwT/hoshikuzu/internal/services/content/albumphoto"
	"github.com/shuTwT/hoshikuzu/internal/services/content/category"
	"github.com/shuTwT/hoshikuzu/internal/services/content/comment"
	"github.com/shuTwT/hoshikuzu/internal/services/content/essay"
	"github.com/shuTwT/hoshikuzu/internal/services/content/flink"
	"github.com/shuTwT/hoshikuzu/internal/services/content/flinkapplication"
	"github.com/shuTwT/hoshikuzu/internal/services/content/friendcircle"
	"github.com/shuTwT/hoshikuzu/internal/services/content/menu"
	"github.com/shuTwT/hoshikuzu/internal/services/content/post"
	"github.com/shuTwT/hoshikuzu/internal/services/content/tag"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/plugin"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/visit"
	"github.com/shuTwT/hoshikuzu/internal/services/mall/product"
	user "github.com/shuTwT/hoshikuzu/internal/services/system/user"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
	"github.com/valyala/fasthttp"
)

type PublicHandler struct {
	visitService            visit.VisitService
	commentService          comment.CommentService
	albumService            album.AlbumService
	albumPhotoService       albumphoto.AlbumPhotoService
	flinkService            flink.FlinkService
	client                  *ent.Client
	friendCircleService     friendcircle.FriendCircleService
	essayService            essay.EssayService
	postService             post.PostService
	categoryService         category.CategoryService
	tagService              tag.TagService
	userService             user.UserService
	productService          product.ProductService
	flinkApplicationService flinkapplication.FlinkApplicationService
	pluginService           plugin.PluginService
	menuService             menu.MenuService
}

func NewPublicHandler(visitService visit.VisitService, commentService comment.CommentService, albumService album.AlbumService, albumPhotoService albumphoto.AlbumPhotoService, flinkService flink.FlinkService, client *ent.Client, friendCircleService friendcircle.FriendCircleService, essayService essay.EssayService, postService post.PostService, categoryService category.CategoryService, tagService tag.TagService, userService user.UserService, productService product.ProductService, flinkApplicationService flinkapplication.FlinkApplicationService, pluginService plugin.PluginService, menuService menu.MenuService) *PublicHandler {
	return &PublicHandler{visitService: visitService, commentService: commentService, albumService: albumService, albumPhotoService: albumPhotoService, flinkService: flinkService, client: client, friendCircleService: friendCircleService, essayService: essayService, postService: postService, categoryService: categoryService, tagService: tagService, userService: userService, productService: productService, flinkApplicationService: flinkApplicationService, pluginService: pluginService, menuService: menuService}
}

// @Summary 处理访客访问
// @Description 处理访客访问请求，记录访问日志
// @Tags 公开接口/访客访问
// @Accept json
// @Produce json
// @Param req body model.VisitLogReq true "访客访问请求"
// @Success 200 {object} model.HttpSuccess{data=ent.VisitLog}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/visit [post]
func (h *PublicHandler) HandleVisitor(c *fiber.Ctx) error {
	ip := c.IP()
	userAgent := c.Context().UserAgent()
	var req model.VisitLogReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	return h.visitService.CreateVisitLog(c.Context(), ip, userAgent, req)
}

// @Summary 处理Twikoo评论事件
// @Description 处理Twikoo评论事件，包括获取配置、获取评论、获取评论数量、提交评论
// @Tags 公开接口/评论
// @Accept json
// @Produce json
// @Param twikoo_req_body body model.TwikooReqBody true "Twikoo请求体"
// @Success 200 {object} model.HttpSuccess{data=interface{}}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/twikoo [get]
// @Router /api/v1/public/twikoo [post]
// @Router /api/v1/public/twikoo [put]
// @Router /api/v1/public/twikoo [delete]
func (h *PublicHandler) HandleTwikoo(c *fiber.Ctx) error {

	var reqBody model.TwikooReqBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	// 配置
	if reqBody.Event == model.TWIKOO_EVENT.GetConfig {
		return c.JSON(fiber.Map{
			"code":   0,
			"config": fiber.Map{},
		})
	}
	// 获取评论
	if reqBody.Event == model.TWIKOO_EVENT.CommentGET {
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
	if reqBody.Event == model.TWIKOO_EVENT.GetCommentsCount {
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
	if reqBody.Event == model.TWIKOO_EVENT.CommentSubmit {
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
// @Tags 公开接口/评论
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]interface{}}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/comment/recent [get]
func (h *PublicHandler) RecentComment(c *fiber.Ctx) error {
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

// @Summary 查询相册列表
// @Description 查询所有相册
// @Tags 公开接口/相册
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Album}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/album/list [get]
func (h *PublicHandler) ListAlbum(c *fiber.Ctx) error {
	albums, err := h.albumService.ListAlbum(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", albums))
}

// @Summary 查询相册列表分页
// @Description 查询相册列表分页
// @Tags 公开接口/相册
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Album]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/album/page [get]
func (h *PublicHandler) ListAlbumPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, albums, err := h.albumService.ListAlbumPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	pageResult := model.PageResult[*ent.Album]{
		Total:   int64(count),
		Records: albums,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 查询相册照片列表
// @Description 查询所有相册照片
// @Tags 公开接口/相册照片
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.AlbumPhoto}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/album-photo/list [get]
func (h *PublicHandler) ListAlbumPhoto(c *fiber.Ctx) error {
	albumPhotos, err := h.albumPhotoService.ListAlbumPhoto(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", albumPhotos))
}

// @Summary 查询相册照片分页列表
// @Description 查询所有相册照片分页列表
// @Tags 公开接口/相册照片
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.AlbumPhoto]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/album-photo/page [get]
func (h *PublicHandler) ListAlbumPhotoPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, albumPhotos, err := h.albumPhotoService.ListAlbumPhotoPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	pageResult := model.PageResult[*ent.AlbumPhoto]{
		Total:   int64(count),
		Records: albumPhotos,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 获取所有Flink
// @Description 获取所有Flink
// @Tags 公开接口/友链
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.FlinkResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/flink/list [get]
func (h *PublicHandler) ListFlink(c *fiber.Ctx) error {
	var listPage model.FlinkListReq
	if err := c.QueryParser(&listPage); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	flinks, err := h.flinkService.ListFlink(c.Context(), listPage)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	result := []model.FlinkResp{}
	for _, flink := range flinks {
		var groupResp *model.FlinkGroupResp
		if flink.Edges.Group != nil {
			groupResp = &model.FlinkGroupResp{
				ID:   flink.Edges.Group.ID,
				Name: flink.Edges.Group.Name,
			}
		}
		result = append(result, model.FlinkResp{
			ID:                 flink.ID,
			CreatedAt:          (model.LocalTime)(flink.CreatedAt),
			UpdatedAt:          (model.LocalTime)(flink.UpdatedAt),
			Name:               flink.Name,
			URL:                flink.URL,
			AvatarURL:          flink.AvatarURL,
			Description:        flink.Description,
			CoverURL:           flink.CoverURL,
			Status:             flink.Status,
			SnapshotUrl:        flink.SnapshotURL,
			Email:              flink.Email,
			EnableFriendCircle: flink.EnableFriendCircle,
			FriendCircleRuleID: flink.FriendCircleRuleID,
			Group:              groupResp,
		})
	}
	return c.JSON(model.NewSuccess("success", result))
}

// @Summary 获取Flink分页列表
// @Description 获取Flink分页列表
// @Tags 公开接口/友链
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.FlinkResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/flink/page [get]
func (h *PublicHandler) ListFlinkPage(c *fiber.Ctx) error {
	var pageQuery model.FlinkPageReq
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	log.Printf("pageQuery: %+v", pageQuery)
	flinks, count, err := h.flinkService.ListFlinkPage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	records := []model.FlinkResp{}
	for _, flink := range flinks {
		records = append(records, model.FlinkResp{
			ID:                 flink.ID,
			CreatedAt:          (model.LocalTime)(flink.CreatedAt),
			UpdatedAt:          (model.LocalTime)(flink.UpdatedAt),
			Name:               flink.Name,
			URL:                flink.URL,
			AvatarURL:          flink.AvatarURL,
			Description:        flink.Description,
			CoverURL:           flink.CoverURL,
			Status:             flink.Status,
			SnapshotUrl:        flink.SnapshotURL,
			Email:              flink.Email,
			EnableFriendCircle: flink.EnableFriendCircle,
			FriendCircleRuleID: flink.FriendCircleRuleID,
		})
	}
	pageResult := model.PageResult[model.FlinkResp]{
		Total:   int64(count),
		Records: records,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 查询友链组列表
// @Description 查询友链组列表
// @Tags 公开接口/友链组
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.FlinkGroupResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/flink-group/list [get]
func (h *PublicHandler) ListFLinkGroup(c *fiber.Ctx) error {
	flinkGroups, err := h.client.FLinkGroup.Query().All(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]*model.FlinkGroupResp, 0)
	for _, flinkGroup := range flinkGroups {
		count, _ := h.flinkService.CountFlinkByGroupID(c.Context(), flinkGroup.ID)
		resp = append(resp, &model.FlinkGroupResp{
			ID:          flinkGroup.ID,
			Name:        flinkGroup.Name,
			Description: flinkGroup.Description,
			Count:       count,
		})
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 获取朋友圈记录分页列表
// @Description 获取朋友圈记录分页列表
// @Tags 公开接口/友链朋友圈
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.FriendCircleRecordResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/friend-circle-record/page [get]
func (h *PublicHandler) ListFriendCircleRecordPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, records, err := h.friendCircleService.ListFriendCircleRecordPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	resps := []model.FriendCircleRecordResp{}
	for _, record := range records {
		resps = append(resps, model.FriendCircleRecordResp{
			ID:          record.ID,
			Title:       record.Title,
			Author:      record.Author,
			LinkURL:     record.LinkURL,
			AvatarURL:   record.AvatarURL,
			PublishedAt: record.PublishedAt,
		})
	}

	pageResult := model.PageResult[model.FriendCircleRecordResp]{
		Total:   int64(count),
		Records: resps,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 获取说说分页列表
// @Description 获取说说分页列表
// @Tags 公开接口/说说
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.EssayResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/essay/page [get]
func (h *PublicHandler) GetEssayPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	essays, total, err := h.essayService.GetEssayPage(c.Context(), pageQuery.Page, pageQuery.Size)
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
// @Tags 公开接口/说说
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} model.HttpSuccess{data=[]model.EssayResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/essay/list [get]
func (h *PublicHandler) ListEssay(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	essays, err := h.essayService.GetEssayList(c.Context(), limit)
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

// @Summary 查询所有文章
// @Description 查询所有文章
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.PostResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/post/list [get]
func (h *PublicHandler) ListPost(c *fiber.Ctx) error {
	var req model.PostListReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	status := "published"
	req.Status = &status
	req.IsVisible = ptr.Bool(true)
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
			PublishedAt:           (*model.LocalTime)(post.PublishedAt),
			ViewCount:             post.ViewCount,
			CommentCount:          post.CommentCount,
			Cover:                 post.Cover,
			Keywords:              post.Keywords,
			Copyright:             post.Copyright,
			Author:                post.Author,
			Summary:               post.Summary,
			CreatedAt:             (model.LocalTime)(post.CreatedAt),
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
// @Tags 公开接口/文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.PostResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/post/page [get]
func (h *PublicHandler) ListPostPage(c *fiber.Ctx) error {
	var req model.PostPageReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	status := "published"
	req.Status = &status
	req.IsVisible = ptr.Bool(true)
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
			PublishedAt:           (*model.LocalTime)(post.PublishedAt),
			ViewCount:             post.ViewCount,
			CommentCount:          post.CommentCount,
			Cover:                 post.Cover,
			Keywords:              post.Keywords,
			Copyright:             post.Copyright,
			Author:                post.Author,
			Summary:               post.Summary,
			CreatedAt:             (model.LocalTime)(post.CreatedAt),
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
// @Router /api/v1/public/post/{id}/summary/stream [get]
func (h *PublicHandler) GetSummaryForStream(c *fiber.Ctx) error {
	// postId, err := strconv.Atoi(c.Params("id"))
	// if err != nil {
	// 	return c.JSON(model.NewError(fiber.StatusBadRequest,
	// 		"Invalid ID format"))
	// }
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	responseChan := make(chan model.AIResponse)
	var targetStr string
	post, err := h.postService.GetRandomPost(c.Context())
	if err != nil {
		targetStr = "看来遇到了点问题，这不是你的问题" + err.Error()
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

			_, err = fmt.Fprintf(w, "data: %s\n\n", string(data))

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
// @Router /api/v1/public/post/month-stats [get]
func (h *PublicHandler) GetPostMonthStats(c *fiber.Ctx) error {
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
// @Router /api/v1/public/post/random [get]
func (h *PublicHandler) GetRandomPost(c *fiber.Ctx) error {
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
// @Router /api/v1/public/post/search [get]
func (h *PublicHandler) SearchPosts(c *fiber.Ctx) error {
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
// @Router /api/v1/public/post/slug/{slug} [get]
func (h *PublicHandler) QueryPostBySlug(c *fiber.Ctx) error {
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
		PublishedAt:           (*model.LocalTime)(post.PublishedAt),
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
		CreatedAt: (model.LocalTime)(post.CreatedAt),
	}
	return c.JSON(model.NewSuccess("success", postResp))
}

// @Summary 查询分类列表
// @Description 查询所有分类的列表
// @Tags 公开接口/分类
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Category}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/category/list [get]
func (h *PublicHandler) QueryCategoryList(c *fiber.Ctx) error {
	categories, err := h.categoryService.QueryCategoryList(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	resps := []model.CategoryResp{}
	for _, cat := range categories {
		var postCount int
		postCount, err = h.postService.PostCountByCategory(c.Context(), cat.ID)
		if err != nil {
			postCount = 0
		}
		resps = append(resps, model.CategoryResp{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Slug:        cat.Slug,
			SortOrder:   cat.SortOrder,
			Active:      cat.Active,
			PostCount:   postCount,
		})
	}

	return c.JSON(model.NewSuccess("success", resps))
}

// @Summary 查询分类分页列表
// @Description 查询所有分类的分页列表
// @Tags 公开接口/分类
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Category]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/category/page [get]
func (h *PublicHandler) QueryCategoryPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, categories, err := h.categoryService.QueryCategoryPage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	resps := []model.CategoryResp{}
	for _, cat := range categories {
		var postCount int
		postCount, err = h.postService.PostCountByCategory(c.Context(), cat.ID)
		if err != nil {
			postCount = 0
		}
		resps = append(resps, model.CategoryResp{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Slug:        cat.Slug,
			SortOrder:   cat.SortOrder,
			Active:      cat.Active,
			PostCount:   postCount,
		})
	}

	pageResult := model.PageResult[model.CategoryResp]{
		Total:   int64(count),
		Records: resps,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 查询标签列表
// @Description 查询所有标签
// @Tags 公开接口/标签
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Tag}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/tag/list [get]
func (h *PublicHandler) QueryTagList(c *fiber.Ctx) error {
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
// @Tags 公开接口/标签
// @Accept json
// @Produce json
// @Param limit query int false "返回数据条数限制"
// @Param offset query int false "返回数据偏移量"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Tag]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/tag/page [get]
func (h *PublicHandler) QueryTagPage(c *fiber.Ctx) error {
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

// @Summary 查询用户分页列表
// @Description 查询所有用户的分页列表
// @Tags 公开接口/用户
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.UserSearchResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/user/search [get]
func (h *PublicHandler) SearchUsers(c *fiber.Ctx) error {
	var req model.UserSearchReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	results, total, err := h.userService.SearchUsers(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*model.UserSearchResp]{
		Total:   int64(total),
		Records: results,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 随机查询Flink
// @Description 随机查询Flink
// @Tags 公开接口/友链
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=ent.FLink}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/flink/random [get]
func (h *PublicHandler) RandomFlink(c *fiber.Ctx) error {
	var req model.FlinkRandomReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	flink, err := h.client.FLink.Query().
		Order(sql.OrderByRand()).
		Limit(req.Limit).
		All(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", flink))
}

// @Summary 搜索商品
// @Description 根据搜索条件查询商品
// @Tags 公开接口/商品
// @Accept json
// @Produce json
// @Param name query string false "商品名称"
// @Param category_id query int false "商品分类ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.ProductSearchResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/product/search [get]
func (h *PublicHandler) SearchProducts(c *fiber.Ctx) error {
	var req model.ProductSearchReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	results, total, err := h.productService.SearchProducts(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*model.ProductSearchResp]{
		Total:   int64(total),
		Records: results,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建友链申请
// @Description 创建一个新的友链申请
// @Tags 公开接口/友链申请
// @Accept json
// @Produce json
// @Param req body model.FlinkApplicationCreateReq true "友链申请创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.FLinkApplication}
// @Failure 400 {object} model.HttpError
// @Router /api/v1/public/flink-application/create [post]
func (h *PublicHandler) CreateFlinkApplication(c *fiber.Ctx) error {
	var createReq *model.FlinkApplicationCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	application, err := h.flinkApplicationService.CreateFlinkApplication(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", application))
}

// @Summary 注册插件
// @Description 注册新插件到系统
// @Tags 公开接口/插件
// @Accept json
// @Produce json
// @Param pluginInfo body model.PluginRegisterReq true "插件注册信息"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/plugin/register [post]
func (h *PublicHandler) RegisterPlugin(c *fiber.Ctx) error {
	// 检查debug模式是否开启
	if !config.GetBool(config.SERVER_DEBUG) {
		slog.Warn("RegisterPlugin called but debug mode is not enabled")
		return c.JSON(model.NewError(fiber.StatusForbidden, "此接口仅在debug模式下可用"))
	}

	// 接收插件注册信息
	var pluginInfo model.PluginRegisterReq
	if err := c.BodyParser(&pluginInfo); err != nil {
		slog.Error("Failed to parse plugin registration info", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "解析插件注册信息失败"))
	}

	// 调用服务层方法存储插件注册信息
	err := h.pluginService.RegisterPlugin(c.Context(), &pluginInfo)
	if err != nil {
		slog.Error("Failed to register plugin", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin registered successfully", "plugin_name", pluginInfo.Name)
	return c.JSON(model.NewSuccess("插件注册成功", nil))
}

// @Summary 插件心跳
// @Description 更新插件的心跳时间
// @Tags 公开接口/插件
// @Accept json
// @Produce json
// @Param heartbeatInfo body model.PluginHeartbeatReq true "插件心跳信息"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/plugin/heartbeat [post]
func (h *PublicHandler) HeartbeatPlugin(c *fiber.Ctx) error {
	// 检查debug模式是否开启
	if !config.GetBool(config.SERVER_DEBUG) {
		slog.Warn("HeartbeatPlugin called but debug mode is not enabled")
		return c.JSON(model.NewError(fiber.StatusForbidden, "此接口仅在debug模式下可用"))
	}

	// 接收插件心跳信息
	var heartbeatInfo model.PluginHeartbeatReq
	if err := c.BodyParser(&heartbeatInfo); err != nil {
		slog.Error("Failed to parse plugin heartbeat info", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "解析插件心跳信息失败"))
	}

	// 调用服务层方法更新插件的心跳时间
	err := h.pluginService.HeartbeatPlugin(c.Context(), &heartbeatInfo)
	if err != nil {
		slog.Error("Failed to update plugin heartbeat", "plugin_name", heartbeatInfo.Name, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Plugin heartbeat updated successfully", "plugin_name", heartbeatInfo.Name)
	return c.JSON(model.NewSuccess("插件心跳更新成功", nil))
}

// @Summary 获取前台菜单列表
// @Description 获取所有可见的前台菜单，用于前端导航栏展示
// @Tags 公开接口/菜单
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.MenuResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/public/menu/list [get]
func (h *PublicHandler) GetMenuList(c *fiber.Ctx) error {
	menus, err := h.menuService.QueryMenuList(c)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", menus))
}
