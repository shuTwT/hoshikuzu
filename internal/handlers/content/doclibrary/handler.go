package doclibrary

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/internal/services/content/doclibrary"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type DocLibraryHandler interface {
	CreateDocLibrary(c *fiber.Ctx) error
	UpdateDocLibrary(c *fiber.Ctx) error
	GetDocLibraryPage(c *fiber.Ctx) error
	GetDocLibrary(c *fiber.Ctx) error
	DeleteDocLibrary(c *fiber.Ctx) error
	GetDocLibraryList(c *fiber.Ctx) error
}

type DocLibraryHandlerImpl struct {
	service doclibrary.DocLibraryService
}

func NewDocLibraryHandlerImpl(service doclibrary.DocLibraryService) DocLibraryHandler {
	return &DocLibraryHandlerImpl{service: service}
}

func (h *DocLibraryHandlerImpl) CreateDocLibrary(c *fiber.Ctx) error {
	createReq := &model.DocLibraryCreateReq{}
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	_, err := h.service.CreateDocLibrary(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *DocLibraryHandlerImpl) UpdateDocLibrary(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	updateReq := &model.DocLibraryUpdateReq{}
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.UpdateDocLibrary(c.Context(), id, updateReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *DocLibraryHandlerImpl) GetDocLibraryPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	libraries, total, err := h.service.GetDocLibraryPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.DocLibraryResp, 0, len(libraries))
	for _, l := range libraries {
		resp = append(resp, model.DocLibraryResp{
			ID:          l.ID,
			Name:        l.Name,
			Alias:       l.Alias,
			Description: l.Description,
			Source:      l.Source.String(),
			URL:         l.URL,
			CreatedAt:   (*model.LocalTime)(&l.CreatedAt),
			UpdatedAt:   (*model.LocalTime)(&l.UpdatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", model.PageResult[model.DocLibraryResp]{
		Total:   int64(total),
		Records: resp,
	}))
}

func (h *DocLibraryHandlerImpl) GetDocLibrary(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	library, err := h.service.GetDocLibrary(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := model.DocLibraryResp{
		ID:          library.ID,
		Name:        library.Name,
		Alias:       library.Alias,
		Description: library.Description,
		Source:      library.Source.String(),
		URL:         library.URL,
		CreatedAt:   (*model.LocalTime)(&library.CreatedAt),
		UpdatedAt:   (*model.LocalTime)(&library.UpdatedAt),
	}
	return c.JSON(model.NewSuccess("success", resp))
}

func (h *DocLibraryHandlerImpl) DeleteDocLibrary(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.DeleteDocLibrary(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *DocLibraryHandlerImpl) GetDocLibraryList(c *fiber.Ctx) error {
	libraries, err := h.service.GetDocLibraryList(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.DocLibraryResp, 0, len(libraries))
	for _, l := range libraries {
		resp = append(resp, model.DocLibraryResp{
			ID:          l.ID,
			Name:        l.Name,
			Alias:       l.Alias,
			Description: l.Description,
			Source:      l.Source.String(),
			URL:         l.URL,
			CreatedAt:   (*model.LocalTime)(&l.CreatedAt),
			UpdatedAt:   (*model.LocalTime)(&l.UpdatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", resp))
}
