package category

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	category_service "github.com/shuTwT/hoshikuzu/internal/services/content/category"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler interface {
	QueryCategory(c *fiber.Ctx) error
	QueryCategoryList(c *fiber.Ctx) error
	QueryCategoryPage(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type CategoryHandlerImpl struct {
	categoryService category_service.CategoryService
}

func NewCategoryHandlerImpl(categoryService category_service.CategoryService) *CategoryHandlerImpl {
	return &CategoryHandlerImpl{
		categoryService: categoryService,
	}
}

func (h *CategoryHandlerImpl) QueryCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	category, err := h.categoryService.QueryCategory(c, id)
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

func (h *CategoryHandlerImpl) QueryCategoryList(c *fiber.Ctx) error {
	categories, err := h.categoryService.QueryCategoryList(c)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.JSON(model.NewSuccess("success", categories))
}

func (h *CategoryHandlerImpl) QueryCategoryPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	err := c.QueryParser(&pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	count, categories, err := h.categoryService.QueryCategoryPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	pageResult := model.PageResult[*ent.Category]{
		Total:   int64(count),
		Records: categories,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

func (h *CategoryHandlerImpl) CreateCategory(c *fiber.Ctx) error {
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

func (h *CategoryHandlerImpl) UpdateCategory(c *fiber.Ctx) error {
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

	updatedCategory, err := h.categoryService.UpdateCategory(c, id, updateReq)
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

func (h *CategoryHandlerImpl) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.categoryService.DeleteCategory(c, id)
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
