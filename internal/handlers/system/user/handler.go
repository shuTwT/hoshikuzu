package user_handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/middleware"
	role_service "github.com/shuTwT/hoshikuzu/internal/services/system/role"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type UserHandler interface {
	ListUser(c *fiber.Ctx) error
	ListUserPage(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	QueryUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetPersonalAccessTokenList(c *fiber.Ctx) error
	GetPersonalAccessToken(c *fiber.Ctx) error
	CreatePat(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
	SearchUsers(c *fiber.Ctx) error
}

type UserHandlerImpl struct {
	userService user_service.UserService
	roleService role_service.RoleService
}

func NewUserHandlerImpl(userService user_service.UserService, roleService role_service.RoleService) *UserHandlerImpl {
	return &UserHandlerImpl{
		userService: userService,
		roleService: roleService,
	}
}

// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} ent.User
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/list [get]
func (h *UserHandlerImpl) ListUser(c *fiber.Ctx) error {
	users, err := h.userService.ListUser(c)
	if err != nil {
		return c.JSON(model.NewError(-1, err.Error()))
	}
	userRespList := []model.UserResp{}
	for _, user := range users {
		userRespList = append(userRespList, model.UserResp{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: &user.RoleID,
		})
	}
	return c.JSON(model.NewSuccess("success", userRespList))
}

// @Summary 获取用户分页列表
// @Description 获取所有用户的分页列表
// @Tags user
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.UserResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/page [get]
func (h *UserHandlerImpl) ListUserPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, users, err := h.userService.ListUserPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	userRespList := []model.UserResp{}
	for _, user := range users {
		userResp := model.UserResp{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: &user.RoleID,
		}
		role, _ := h.roleService.QueryRole(c.Context(), user.RoleID)
		userResp.Role = role
		userRespList = append(userRespList, userResp)
	}
	pageResult := model.PageResult[model.UserResp]{
		Total:   int64(count),
		Records: userRespList,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建用户
// @Description 创建一个新用户
// @Tags user
// @Accept json
// @Produce json
// @Param user body ent.User true "用户信息"
// @Success 201 {object} ent.User
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/create [post]
func (h *UserHandlerImpl) CreateUser(c *fiber.Ctx) error {
	var req model.UserCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	user, err := h.userService.CreateUser(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", user))
}

// @Summary 更新用户
// @Description 更新指定用户的信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param user body ent.User true "用户信息"
// @Success 200 {object} ent.User
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/update/{id} [put]
func (h *UserHandlerImpl) UpdateUser(c *fiber.Ctx) error {
	var err error
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var req model.UserUpdateReq
	if err = c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	updatedUser, err := h.userService.UpdateUser(c.Context(), id, req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", updatedUser))
}

// @Summary 查询用户
// @Description 查询指定用户的详细信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} ent.User
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/query/{id} [get]
func (h *UserHandlerImpl) QueryUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	user, err := h.userService.QueryUserById(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	userResp := model.UserResp{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: &user.PhoneNumber,
		RoleID:      &user.RoleID,
	}

	return c.JSON(model.NewSuccess("success", userResp))
}

// @Summary 删除用户
// @Description 删除指定用户
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/delete/{id} [delete]
func (h *UserHandlerImpl) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.userService.DeleteUser(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 查询个人令牌列表
// @Description 查询当前用户的所有个人令牌
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} []model.PersonalAccessTokenListResp
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/personal-access-token/list [get]
func (h *UserHandlerImpl) GetPersonalAccessTokenList(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(
			fiber.StatusUnauthorized, "Unauthorized",
		))
	}
	userId := loginUser.ID
	tokens, err := h.userService.GetPersonalAccessTokenList(c.Context(), userId)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	result := []model.PersonalAccessTokenListResp{}

	for _, token := range tokens {
		result = append(result, model.PersonalAccessTokenListResp{
			ID:          token.ID,
			Name:        token.Name,
			Expires:     model.ParseTime(token.Expires),
			Description: token.Name,
		})
	}

	return c.JSON(model.NewSuccess("success", result))
}

// @Summary 查询个人令牌
// @Description 查询指定个人令牌的详细信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "个人令牌ID"
// @Success 200 {object} model.PersonalAccessTokenResp
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/personal-access-token/query/{id} [get]
func (h *UserHandlerImpl) GetPersonalAccessToken(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(
			fiber.StatusInternalServerError, "Unauthorized",
		))
	}
	userId := loginUser.ID
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(
			fiber.StatusBadRequest, err.Error(),
		))
	}
	token, err := h.userService.GetPersonalAccessToken(c.Context(), userId, id)
	if err != nil {
		return c.JSON(model.NewError(
			fiber.StatusBadRequest, err.Error(),
		))
	}
	result := model.PersonalAccessTokenResp{
		ID:          token.ID,
		Name:        token.Name,
		Expires:     model.ParseTime(token.Expires),
		Description: token.Description,
		Token:       token.Token,
	}
	return c.JSON(model.NewSuccess("success", result))
}

// @Summary 创建 personalAccessToken 个人令牌
// @Description 创建一个新的个人令牌
// @Tags user
// @Accept json
// @Produce json
// @Param createReq body model.PersonalAccessTokenCreateReq true "个人令牌创建请求"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/personal-access-token/create [post]
func (h *UserHandlerImpl) CreatePat(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(
			fiber.StatusInternalServerError, "Unauthorized",
		))
	}
	userId := loginUser.ID
	var createReq model.PersonalAccessTokenCreateReq
	err := c.BodyParser(&createReq)
	if err != nil {
		return c.JSON(model.NewError(
			fiber.StatusBadRequest, err.Error(),
		))
	}
	// 查找用户
	_, err = h.userService.CreatePersonalAccessToken(c.Context(), userId, createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 查询用户个人信息
// @Description 查询指定用户的个人信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} model.UserProfileResp
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/user/profile [get]
func (h *UserHandlerImpl) GetUserProfile(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(
			fiber.StatusInternalServerError, "Unauthorized",
		))
	}
	userId := loginUser.ID

	user, err := h.userService.QueryUserById(c.Context(), userId)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"User not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	role, err := h.roleService.QueryRole(c.Context(), user.RoleID)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	result := model.UserProfileResp{
		UserID:              userId,
		Email:               user.Email,
		EmailVerified:       user.EmailVerified,
		Name:                user.Name,
		PhoneNumber:         &user.PhoneNumber,
		PhoneNumberVerified: user.PhoneNumberVerified,
		Role:                role,
	}

	return c.JSON(model.NewSuccess("success", result))
}

func (h *UserHandlerImpl) SearchUsers(c *fiber.Ctx) error {
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
