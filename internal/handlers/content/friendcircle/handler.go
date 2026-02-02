package friendcircle

import (
	"strconv"

	friendcircle "github.com/shuTwT/hoshikuzu/internal/services/content/friendcircle"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type FriendCircleHandler interface {
	ListFriendCircleRecordPage(c *fiber.Ctx) error
	CreateFriendCircleRecord(c *fiber.Ctx) error
	UpdateFriendCircleRecord(c *fiber.Ctx) error
	DeleteFriendCircleRecord(c *fiber.Ctx) error
}

type FriendCircleHandlerImpl struct {
	friendCircleService friendcircle.FriendCircleService
}

func NewFriendCircleHandlerImpl(friendCircleService friendcircle.FriendCircleService) *FriendCircleHandlerImpl {
	return &FriendCircleHandlerImpl{
		friendCircleService: friendCircleService,
	}
}

// @Summary 获取朋友圈记录分页列表
// @Description 获取朋友圈记录分页列表
// @Tags friend_circle_records
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.FriendCircleRecordResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/friend_circle_records [get]
func (h *FriendCircleHandlerImpl) ListFriendCircleRecordPage(c *fiber.Ctx) error {
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

// @Summary 创建朋友圈记录
// @Description 创建朋友圈记录
// @Tags friend_circle_records
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=ent.FriendCircleRecord}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/friend_circle_records [post]
func (h *FriendCircleHandlerImpl) CreateFriendCircleRecord(c *fiber.Ctx) error {
	var req *model.FriendCircleRecordSaveReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	newRecord, err := h.friendCircleService.CreateFriendCircleRecord(c.Context(), req.Author, req.Title, req.LinkURL, req.AvatarURL, "", "")
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", newRecord))
}

// @Summary 更新朋友圈记录
// @Description 更新朋友圈记录
// @Tags friend_circle_records
// @Accept json
// @Produce json
// @Param id path int true "朋友圈记录 ID"
// @Success 200 {object} model.HttpSuccess{data=ent.FriendCircleRecord}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/friend_circle_records/{id} [put]
func (h *FriendCircleHandlerImpl) UpdateFriendCircleRecord(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var req *model.FriendCircleRecordSaveReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedRecord, err := h.friendCircleService.UpdateFriendCircleRecord(c.Context(), id, req.Author, req.Title, req.LinkURL, req.AvatarURL, "", "")
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", updatedRecord))
}

// @Summary 删除朋友圈记录
// @Description 删除朋友圈记录
// @Tags friend_circle_records
// @Accept json
// @Produce json
// @Param id path int true "朋友圈记录 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/friend_circle_records/{id} [delete]
func (h *FriendCircleHandlerImpl) DeleteFriendCircleRecord(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	err = h.friendCircleService.DeleteFriendCircleRecord(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
