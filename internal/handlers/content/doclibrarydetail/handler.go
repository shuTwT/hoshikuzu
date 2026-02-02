package doclibrarydetail

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/internal/services/content/doclibrarydetail"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type DocLibraryDetailHandler interface {
	CreateDocLibraryDetail(c *fiber.Ctx) error
	UpdateDocLibraryDetail(c *fiber.Ctx) error
	GetDocLibraryDetailPage(c *fiber.Ctx) error
	GetDocLibraryDetail(c *fiber.Ctx) error
	DeleteDocLibraryDetail(c *fiber.Ctx) error
	GetDocLibraryDetailTree(c *fiber.Ctx) error
}

type DocLibraryDetailHandlerImpl struct {
	service doclibrarydetail.DocLibraryDetailService
}

func NewDocLibraryDetailHandlerImpl(service doclibrarydetail.DocLibraryDetailService) DocLibraryDetailHandler {
	return &DocLibraryDetailHandlerImpl{service: service}
}

func (h *DocLibraryDetailHandlerImpl) CreateDocLibraryDetail(c *fiber.Ctx) error {
	createReq := model.DocLibraryDetailCreateReq{}
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	_, err := h.service.CreateDocLibraryDetail(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *DocLibraryDetailHandlerImpl) UpdateDocLibraryDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	updateReq := model.DocLibraryDetailUpdateReq{}
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.UpdateDocLibraryDetail(c.Context(), id, updateReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *DocLibraryDetailHandlerImpl) GetDocLibraryDetailPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	libraryID := c.QueryInt("library_id", 0)

	details, total, err := h.service.GetDocLibraryDetailPage(c.Context(), pageQuery.Page, pageQuery.Size, libraryID)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := make([]model.DocLibraryDetailResp, 0, len(details))
	for _, d := range details {
		resp = append(resp, model.DocLibraryDetailResp{
			ID:        d.ID,
			LibraryID: d.LibraryID,
			Title:     d.Title,
			Version:   d.Version,
			Content:   d.Content,
			ParentID:  &d.ParentID,
			Path:      d.Path,
			URL:       d.URL,
			Language:  d.Language,
			CreatedAt: (*model.LocalTime)(&d.CreatedAt),
			UpdatedAt: (*model.LocalTime)(&d.UpdatedAt),
		})
	}
	return c.JSON(model.NewSuccess("success", model.PageResult[model.DocLibraryDetailResp]{
		Total:   int64(total),
		Records: resp,
	}))
}

func (h *DocLibraryDetailHandlerImpl) GetDocLibraryDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	detail, err := h.service.GetDocLibraryDetail(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	resp := model.DocLibraryDetailResp{
		ID:        detail.ID,
		LibraryID: detail.LibraryID,
		Title:     detail.Title,
		Version:   detail.Version,
		Content:   detail.Content,
		ParentID:  &detail.ParentID,
		Path:      detail.Path,
		URL:       detail.URL,
		Language:  detail.Language,
		CreatedAt: (*model.LocalTime)(&detail.CreatedAt),
		UpdatedAt: (*model.LocalTime)(&detail.UpdatedAt),
	}
	return c.JSON(model.NewSuccess("success", resp))
}

func (h *DocLibraryDetailHandlerImpl) DeleteDocLibraryDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	err = h.service.DeleteDocLibraryDetail(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *DocLibraryDetailHandlerImpl) GetDocLibraryDetailTree(c *fiber.Ctx) error {
	libraryID := c.QueryInt("library_id", 0)

	details, err := h.service.GetDocLibraryDetailTree(c.Context(), libraryID)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	tree := buildTree(details, nil)
	return c.JSON(model.NewSuccess("success", tree))
}

func buildTree(details []*ent.DocLibraryDetail, parentID *int) []model.DocLibraryDetailTreeResp {
	var tree []model.DocLibraryDetailTreeResp

	for _, detail := range details {
		if (parentID == nil && detail.ParentID == 0) || (parentID != nil && detail.ParentID != 0 && detail.ParentID == *parentID) {
			node := model.DocLibraryDetailTreeResp{
				ID:        detail.ID,
				LibraryID: detail.LibraryID,
				Title:     detail.Title,
				Version:   detail.Version,
				Content:   detail.Content,
				ParentID:  &detail.ParentID,
				Path:      detail.Path,
				URL:       detail.URL,
				Language:  detail.Language,
				CreatedAt: (*model.LocalTime)(&detail.CreatedAt),
				UpdatedAt: (*model.LocalTime)(&detail.UpdatedAt),
			}
			children := buildTree(details, &detail.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}

	return tree
}
