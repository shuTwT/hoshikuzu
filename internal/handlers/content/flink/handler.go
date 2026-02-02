package flink

import (
	"log"
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/flink"
	flink_service "github.com/shuTwT/hoshikuzu/internal/services/content/flink"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"
)

type FlinkHandler interface {
	ListFlink(c *fiber.Ctx) error
	ListFlinkPage(c *fiber.Ctx) error
	CreateFlink(c *fiber.Ctx) error
	UpdateFlink(c *fiber.Ctx) error
	QueryFlink(c *fiber.Ctx) error
	DeleteFlink(c *fiber.Ctx) error
	RandomFlink(c *fiber.Ctx) error
}

type FlinkHandlerImpl struct {
	client       *ent.Client
	flinkService flink_service.FlinkService
}

func NewFlinkHandlerImpl(client *ent.Client, flinkService flink_service.FlinkService) *FlinkHandlerImpl {
	return &FlinkHandlerImpl{
		client:       client,
		flinkService: flinkService,
	}
}

// @Summary 获取所有Flink
// @Description 获取所有Flink
// @Tags flink
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.FlinkResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/list [get]
func (h *FlinkHandlerImpl) ListFlink(c *fiber.Ctx) error {
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
			CreatedAt:          &flink.CreatedAt,
			UpdatedAt:          &flink.UpdatedAt,
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
// @Tags flink
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.FlinkResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/page [get]
func (h *FlinkHandlerImpl) ListFlinkPage(c *fiber.Ctx) error {
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
			CreatedAt:          &flink.CreatedAt,
			UpdatedAt:          &flink.UpdatedAt,
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

// @Summary 创建Flink
// @Description 创建Flink
// @Tags flink
// @Accept json
// @Produce json
// @Param flink_create_req body model.FlinkCreateReq true "Flink创建请求体"
// @Success 200 {object} model.HttpSuccess{data=ent.FLink}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/create [post]
func (h *FlinkHandlerImpl) CreateFlink(c *fiber.Ctx) error {
	var createReq *model.FlinkCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	flink, err := h.client.FLink.Create().
		SetName(createReq.Name).
		SetURL(createReq.URL).
		SetAvatarURL(createReq.AvatarURL).
		SetDescription(createReq.Description).
		SetCoverURL(createReq.CoverURL).
		SetSnapshotURL(createReq.SnapshotURL).
		SetEmail(createReq.Email).
		SetEnableFriendCircle(createReq.EnableFriendCircle).
		SetNillableFriendCircleRuleID(createReq.FriendCircleRuleID).
		SetGroupID(createReq.GroupID).
		Save(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", flink))
}

// @Summary 更新Flink
// @Description 更新Flink
// @Tags flink
// @Accept json
// @Produce json
// @Param id path int true "Flink ID"
// @Param flink_update_req body model.FlinkUpdateReq true "Flink更新请求体"
// @Success 200 {object} model.HttpSuccess{data=ent.FLink}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/update/{id} [put]
func (h *FlinkHandlerImpl) UpdateFlink(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	var updateReq *model.FlinkUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	flink, err := h.client.FLink.UpdateOneID(id).
		SetName(updateReq.Name).
		SetURL(updateReq.URL).
		SetAvatarURL(updateReq.AvatarURL).
		SetDescription(updateReq.Description).
		SetCoverURL(updateReq.CoverURL).
		SetSnapshotURL(updateReq.SnapshotURL).
		SetEmail(updateReq.Email).
		SetEnableFriendCircle(updateReq.EnableFriendCircle).
		SetFriendCircleRuleID(updateReq.FriendCircleRuleID).
		SetGroupID(updateReq.GroupID).
		Save(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", flink))
}

// @Summary 查询Flink
// @Description 查询Flink
// @Tags flink
// @Accept json
// @Produce json
// @Param id path int true "Flink ID"
// @Success 200 {object} model.HttpSuccess{data=ent.FLink}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/query/{id} [get]
func (h *FlinkHandlerImpl) QueryFlink(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	flink, err := h.client.FLink.Query().
		Where(flink.ID(id)).
		First(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", flink))
}

// @Summary 删除Flink
// @Description 删除Flink
// @Tags flink
// @Accept json
// @Produce json
// @Param id path int true "Flink ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/{id} [delete]
func (h *FlinkHandlerImpl) DeleteFlink(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}
	if err := h.client.FLink.DeleteOneID(id).Exec(c.Context()); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 随机查询Flink
// @Description 随机查询Flink
// @Tags flink
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=ent.FLink}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink/random [get]
func (h *FlinkHandlerImpl) RandomFlink(c *fiber.Ctx) error {
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
