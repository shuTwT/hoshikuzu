package role

import (
	"strconv"

	role_service "github.com/shuTwT/hoshikuzu/internal/services/system/role"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler interface {
	ListRole(c *fiber.Ctx) error
	ListRolePage(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	QueryRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
}

type RoleHandlerImpl struct {
	roleService role_service.RoleService
}

func NewRoleHandlerImpl(roleService role_service.RoleService) *RoleHandlerImpl {
	return &RoleHandlerImpl{
		roleService: roleService,
	}
}

// @Summary 查询所有角色
// @Description 查询系统中所有角色的列表
// @Tags 后台管理接口/角色
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Role}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/role/list [get]
func (h *RoleHandlerImpl) ListRole(c *fiber.Ctx) error {
	roles, err := h.roleService.QueryRoleList(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	resps := make([]model.RoleResp, 0, len(roles))

	for _, role := range roles {
		resps = append(resps, model.RoleResp{
			ID:        role.ID,
			Name:      role.Name,
			Code:      role.Code,
			CreatedAt: model.LocalTime(role.CreatedAt),
			IsDefault: role.IsDefault,
		})
	}

	return c.JSON(model.NewSuccess("success", resps))
}

// @Summary 查询角色分页列表
// @Description 查询系统中角色的分页列表
// @Tags 后台管理接口/角色
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Role]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/role/page [get]
func (h *RoleHandlerImpl) ListRolePage(c *fiber.Ctx) error {
	var pageQuery = model.PageQuery{}
	err := c.QueryParser(&pageQuery)

	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	roles, count, err := h.roleService.QueryRolePage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	resps := make([]model.RoleResp, 0, len(roles))

	for _, role := range roles {
		resps = append(resps, model.RoleResp{
			ID:        role.ID,
			Name:      role.Name,
			Code:      role.Code,
			CreatedAt: model.LocalTime(role.CreatedAt),
			IsDefault: role.IsDefault,
		})
	}

	pageResult := model.PageResult[model.RoleResp]{
		Total:   int64(count),
		Records: resps,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建角色
// @Description 创建一个新的角色
// @Tags 后台管理接口/角色
// @Accept json
// @Produce json
// @Param createReq body model.RoleCreateReq true "角色创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Role}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/role/create [post]
func (h *RoleHandlerImpl) CreateRole(c *fiber.Ctx) error {
	var roleData model.RoleCreateReq
	if err := c.BodyParser(&roleData); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	role, err := h.roleService.CreateRole(c.Context(), roleData)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", role))
}

// @Summary 更新角色
// @Description 更新指定角色的信息
// @Tags 后台管理接口/角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param updateReq body model.RoleUpdateReq true "角色更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Role}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/role/update/{id} [put]
func (h *RoleHandlerImpl) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var roleData model.RoleUpdateReq
	if err = c.BodyParser(&roleData); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	role, err := h.roleService.UpdateRole(c.Context(), id, roleData)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", role))
}

// @Summary 查询角色
// @Description 查询指定角色的详细信息
// @Tags 后台管理接口/角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Role}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/role/query/{id} [get]
func (h *RoleHandlerImpl) QueryRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}
	role, err := h.roleService.QueryRole(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", role))
}

// @Summary 删除角色
// @Description 删除指定角色
// @Tags 后台管理接口/角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/role/delete/{id} [delete]
func (h *RoleHandlerImpl) DeleteRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}
	err = h.roleService.DeleteRole(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", nil))
}
