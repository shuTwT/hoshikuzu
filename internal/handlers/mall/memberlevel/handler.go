package memberlevel

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	memberlevel_service "github.com/shuTwT/hoshikuzu/internal/services/mall/memberlevel"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MemberLevelHandler interface {
	QueryMemberLevel(c *fiber.Ctx) error
	QueryMemberLevelList(c *fiber.Ctx) error
	QueryMemberLevelPage(c *fiber.Ctx) error
	CreateMemberLevel(c *fiber.Ctx) error
	UpdateMemberLevel(c *fiber.Ctx) error
	DeleteMemberLevel(c *fiber.Ctx) error
}

type MemberLevelHandlerImpl struct {
	memberLevelService memberlevel_service.MemberLevelService
}

func NewMemberLevelHandlerImpl(memberLevelService memberlevel_service.MemberLevelService) *MemberLevelHandlerImpl {
	return &MemberLevelHandlerImpl{
		memberLevelService: memberLevelService,
	}
}

func (h *MemberLevelHandlerImpl) QueryMemberLevel(c *fiber.Ctx) error {
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

func (h *MemberLevelHandlerImpl) QueryMemberLevelList(c *fiber.Ctx) error {
	memberLevels, err := h.memberLevelService.QueryMemberLevelList(c)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", memberLevels))
}

func (h *MemberLevelHandlerImpl) QueryMemberLevelPage(c *fiber.Ctx) error {
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

func (h *MemberLevelHandlerImpl) CreateMemberLevel(c *fiber.Ctx) error {
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

func (h *MemberLevelHandlerImpl) UpdateMemberLevel(c *fiber.Ctx) error {
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

func (h *MemberLevelHandlerImpl) DeleteMemberLevel(c *fiber.Ctx) error {
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
