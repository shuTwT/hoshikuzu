package flinkgroup

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/content/flink"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type FlinkGroupHandler interface {
	ListFLinkGroup(c *fiber.Ctx) error
	CreateFlinkGroup(c *fiber.Ctx) error
	UpdateFlinkGroup(c *fiber.Ctx) error
	DeleteFLinkGroup(c *fiber.Ctx) error
}

type FlinkGroupHandlerImpl struct {
	client       *ent.Client
	flinkService flink.FlinkService
}

func NewFlinkGroupHandlerImpl(client *ent.Client, flinkService flink.FlinkService) *FlinkGroupHandlerImpl {
	return &FlinkGroupHandlerImpl{
		client:       client,
		flinkService: flinkService,
	}
}

// @Summary 查询友链组列表
// @Description 查询友链组列表
// @Tags 后台管理接口/友链组
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.FlinkGroupResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink-group/list [get]
func (h *FlinkGroupHandlerImpl) ListFLinkGroup(c *fiber.Ctx) error {
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

// @Summary 创建友链组
// @Description 创建友链组
// @Tags 后台管理接口/友链组
// @Accept json
// @Produce json
// @Param flink_group_create_req body ent.FLinkGroup true "FlinkGroup创建请求体"
// @Success 200 {object} model.HttpSuccess{data=ent.FLinkGroup}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink-group/create [post]
func (h *FlinkGroupHandlerImpl) CreateFlinkGroup(c *fiber.Ctx) error {
	var createReq *ent.FLinkGroup
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	flinkGroup, err := h.client.FLinkGroup.Create().
		SetName(createReq.Name).
		SetDescription(createReq.Description).
		Save(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", flinkGroup))
}

// @Summary 更新友链组
// @Description 更新友链组
// @Tags 后台管理接口/友链组
// @Accept json
// @Produce json
// @Param id path int true "友链组 ID"
// @Param flink_group_update_req body ent.FLinkGroup true "友链组更新请求体"
// @Success 200 {object} model.HttpSuccess{data=ent.FLinkGroup}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink-group/update/{id} [put]
func (h *FlinkGroupHandlerImpl) UpdateFlinkGroup(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	var updateReq *ent.FLinkGroup
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	flinkGroup, err := h.client.FLinkGroup.UpdateOneID(id).
		SetName(updateReq.Name).
		SetDescription(updateReq.Description).
		Save(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", flinkGroup))
}

// @Summary 删除友链组
// @Description 删除指定友链组
// @Tags 后台管理接口/友链组
// @Accept json
// @Produce json
// @Param id path int true "友链组 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/flink-group/delete/{id} [delete]
func (h *FlinkGroupHandlerImpl) DeleteFLinkGroup(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.client.FLinkGroup.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}
