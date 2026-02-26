package memberlevel

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	memberlevel_service "github.com/shuTwT/hoshikuzu/internal/services/mall/memberlevel"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MemberLevelHandler struct {
	memberLevelService memberlevel_service.MemberLevelService
}

func NewMemberLevelHandler(memberLevelService memberlevel_service.MemberLevelService) *MemberLevelHandler {
	return &MemberLevelHandler{
		memberLevelService: memberLevelService,
	}
}

// @Summary 查询会员等级
// @Description 查询会员等级
// @Tags 后台管理接口/会员等级
// @Accept json
// @Produce json
// @Param id path int true "会员等级 ID"
// @Success 200 {object} model.HttpSuccess{data=model.MemberLevelResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member-level/query/{id} [get]
func (h *MemberLevelHandler) QueryMemberLevel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	ml, err := h.memberLevelService.QueryMemberLevel(c, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Member level not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", ml))
}

// @Summary 查询会员等级列表
// @Description 查询会员等级列表
// @Tags 后台管理接口/会员等级
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.MemberLevelResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member-level/list [get]
func (h *MemberLevelHandler) QueryMemberLevelList(c *fiber.Ctx) error {
	memberLevels, err := h.memberLevelService.QueryMemberLevelList(c)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", memberLevels))
}

// @Summary 查询会员等级分页
// @Description 查询会员等级分页
// @Tags 后台管理接口/会员等级
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.MemberLevelResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member-level/page [get]
func (h *MemberLevelHandler) QueryMemberLevelPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, memberLevels, err := h.memberLevelService.QueryMemberLevelPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	pageResult := model.PageResult[*ent.MemberLevel]{
		Total:   int64(count),
		Records: memberLevels,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建会员等级
// @Description 创建会员等级
// @Tags 后台管理接口/会员等级
// @Accept json
// @Produce json
// @Param member_level_create_req body model.MemberLevelCreateReq true "会员等级创建请求体"
// @Success 200 {object} model.HttpSuccess{data=model.MemberLevelResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member-level/create [post]
func (h *MemberLevelHandler) CreateMemberLevel(c *fiber.Ctx) error {
	var createReq *model.MemberLevelCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	ml, err := h.memberLevelService.CreateMemberLevel(c, createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", ml))
}

// @Summary 更新会员等级
// @Description 更新会员等级
// @Tags 后台管理接口/会员等级
// @Accept json
// @Produce json
// @Param id path int true "会员等级 ID"
// @Param member_level_update_req body model.MemberLevelUpdateReq true "会员等级更新请求体"
// @Success 200 {object} model.HttpSuccess{data=model.MemberLevelResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member-level/update/{id} [put]
func (h *MemberLevelHandler) UpdateMemberLevel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var updateReq *model.MemberLevelUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	updatedMemberLevel, err := h.memberLevelService.UpdateMemberLevel(c, id, updateReq)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Member level not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", updatedMemberLevel))
}

// @Summary 删除会员等级
// @Description 删除会员等级
// @Tags 后台管理接口/会员等级
// @Accept json
// @Produce json
// @Param id path int true "会员等级 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member-level/delete/{id} [delete]
func (h *MemberLevelHandler) DeleteMemberLevel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.memberLevelService.DeleteMemberLevel(c, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Member level not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
