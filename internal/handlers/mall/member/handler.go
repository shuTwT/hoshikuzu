package member

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	member_service "github.com/shuTwT/hoshikuzu/internal/services/mall/member"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MemberHandler interface {
	QueryMember(c *fiber.Ctx) error
	QueryMemberPage(c *fiber.Ctx) error
	CreateMember(c *fiber.Ctx) error
	UpdateMember(c *fiber.Ctx) error
	DeleteMember(c *fiber.Ctx) error
}

type MemberHandlerImpl struct {
	userService   user_service.UserService
	memberService member_service.MemberService
}

func NewMemberHandlerImpl(userService user_service.UserService, memberService member_service.MemberService) *MemberHandlerImpl {
	return &MemberHandlerImpl{
		userService:   userService,
		memberService: memberService,
	}
}

func (h *MemberHandlerImpl) QueryMember(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid user ID format",
		))
	}

	m, err := h.memberService.QueryMember(c, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Member not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", m))
}

func (h *MemberHandlerImpl) QueryMemberPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, members, err := h.memberService.QueryMemberPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	memberResps := make([]*model.MemberResp, 0, len(members))
	for _, m := range members {
		user, err := h.userService.QueryUserById(c.Context(), m.UserID)
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError,
				err.Error(),
			))
		}
		memberResps = append(memberResps, &model.MemberResp{
			ID:          m.ID,
			UserID:      m.UserID,
			UserName:    user.Name,
			MemberLevel: m.MemberLevel,
			MemberNo:    m.MemberNo,
			JoinTime:    m.JoinTime.Format("2006-01-02 15:04:05"),
			ExpireTime:  m.ExpireTime.Format("2006-01-02 15:04:05"),
			Points:      m.Points,
			TotalSpent:  m.TotalSpent,
			OrderCount:  m.OrderCount,
			Active:      m.Active,
			Remark:      m.Remark,
		})
	}
	pageResult := model.PageResult[*model.MemberResp]{
		Total:   int64(count),
		Records: memberResps,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

func (h *MemberHandlerImpl) CreateMember(c *fiber.Ctx) error {
	var createReq *model.MemberCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	m, err := h.memberService.CreateMember(c.Context(), createReq.UserID, createReq.MemberLevel, createReq.MemberNo)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", m))
}

func (h *MemberHandlerImpl) UpdateMember(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var updateReq *model.MemberUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	updatedMember, err := h.memberService.UpdateMember(c, id, updateReq)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Member not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", updatedMember))
}

func (h *MemberHandlerImpl) DeleteMember(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.memberService.DeleteMember(c, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Member not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
