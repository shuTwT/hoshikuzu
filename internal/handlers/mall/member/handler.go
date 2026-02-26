package member

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	member_service "github.com/shuTwT/hoshikuzu/internal/services/mall/member"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MemberHandler struct {
	userService   user_service.UserService
	memberService member_service.MemberService
}

func NewMemberHandler(userService user_service.UserService, memberService member_service.MemberService) *MemberHandler {
	return &MemberHandler{
		userService:   userService,
		memberService: memberService,
	}
}

// @Summary 查询会员
// @Description 查询会员
// @Tags 后台管理接口/会员
// @Accept json
// @Produce json
// @Param user_id path int true "用户 ID"
// @Success 200 {object} model.HttpSuccess{data=model.MemberResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member/query/{user_id} [get]
func (h *MemberHandler) QueryMember(c *fiber.Ctx) error {
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

// @Summary 查询会员列表分页
// @Description 查询会员列表分页
// @Tags 后台管理接口/会员
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.MemberResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member/page [get]
func (h *MemberHandler) QueryMemberPage(c *fiber.Ctx) error {
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

// @Summary 创建会员
// @Description 创建会员
// @Tags 后台管理接口/会员
// @Accept json
// @Produce json
// @Param member_create_req body model.MemberCreateReq true "会员创建请求体"
// @Success 200 {object} model.HttpSuccess{data=model.MemberResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member/create [post]
func (h *MemberHandler) CreateMember(c *fiber.Ctx) error {
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

// @Summary 更新会员
// @Description 更新会员
// @Tags 后台管理接口/会员
// @Accept json
// @Produce json
// @Param id path int true "会员 ID"
// @Param member_update_req body model.MemberUpdateReq true "会员更新请求体"
// @Success 200 {object} model.HttpSuccess{data=model.MemberResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member/update/{id} [put]
func (h *MemberHandler) UpdateMember(c *fiber.Ctx) error {
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

// @Summary 删除会员
// @Description 删除会员
// @Tags 后台管理接口/会员
// @Accept json
// @Produce json
// @Param id path int true "会员 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/member/delete/{id} [delete]
func (h *MemberHandler) DeleteMember(c *fiber.Ctx) error {
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
