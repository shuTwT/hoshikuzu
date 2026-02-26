package category

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	category_service "github.com/shuTwT/hoshikuzu/internal/services/content/category"
	post_service "github.com/shuTwT/hoshikuzu/internal/services/content/post"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService category_service.CategoryService
	postService     post_service.PostService
}

func NewCategoryHandler(categoryService category_service.CategoryService, postService post_service.PostService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		postService:     postService,
	}
}

// @Summary 查询分类
// @Description 查询指定分类的信息
// @Tags 后台管理接口/分类
// @Accept json
// @Produce json
// @Param id path int true "分类 ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Category}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/category/query/{id} [get]
func (h *CategoryHandler) QueryCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	category, err := h.categoryService.QueryCategory(c.Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Category not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", category))
}

// @Summary 查询分类列表
// @Description 查询所有分类的列表
// @Tags 后台管理接口/分类
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Category}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/category/list [get]
func (h *CategoryHandler) QueryCategoryList(c *fiber.Ctx) error {
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
// @Tags 后台管理接口/分类
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Category]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/category/page [get]
func (h *CategoryHandler) QueryCategoryPage(c *fiber.Ctx) error {
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

// @Summary 创建分类
// @Description 创建一个新的分类
// @Tags 后台管理接口/分类
// @Accept json
// @Produce json
// @Param category body model.CategoryCreateReq true "分类创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Category}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/category/create [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var createReq model.CategoryCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	category, err := h.categoryService.CreateCategory(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", category))
}

// @Summary 更新分类
// @Description 更新指定分类的信息
// @Tags 后台管理接口/分类
// @Accept json
// @Produce json
// @Param id path int true "分类 ID"
// @Param category body model.CategoryUpdateReq true "分类更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.Category}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/category/update/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	var updateReq *model.CategoryUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	updatedCategory, err := h.categoryService.UpdateCategory(c.Context(), id, updateReq)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Category not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", updatedCategory))
}

// @Summary 删除分类
// @Description 删除指定分类
// @Tags 后台管理接口/分类
// @Accept json
// @Produce json
// @Param id path int true "分类 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/category/delete/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"Category not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
