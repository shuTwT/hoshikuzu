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

// @Summary 创建文档库详情
// @Description 创建新的文档库详情
// @Tags 文档库详情
// @Accept json
// @Produce json
// @Param docLibraryDetail body model.DocLibraryDetailCreateReq true "文档库详情创建请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library-detail/create [post]
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

// @Summary 更新文档库详情
// @Description 更新指定文档库详情的信息
// @Tags 文档库详情
// @Accept json
// @Produce json
// @Param id path int true "文档库详情 ID"
// @Param docLibraryDetail body model.DocLibraryDetailUpdateReq true "文档库详情更新请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library-detail/update/{id} [put]
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

// @Summary 获取文档库详情分页列表
// @Description 获取指定文档库下的文档库详情分页列表
// @Tags 文档库详情
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param library_id query int true "文档库 ID"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.DocLibraryDetailResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library-detail/page [get]
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

// @Summary 获取文档库详情
// @Description 获取指定文档库详情的信息
// @Tags 文档库详情
// @Accept json
// @Produce json
// @Param id path int true "文档库详情 ID"
// @Success 200 {object} model.HttpSuccess{data=model.DocLibraryDetailResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library-detail/query/{id} [get]
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

// @Summary 删除文档库详情
// @Description 删除指定文档库详情
// @Tags 文档库详情
// @Accept json
// @Produce json
// @Param id path int true "文档库详情 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library-detail/delete/{id} [delete]
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

// @Summary 获取文档库详情树
// @Description 获取指定文档库下的文档库详情树结构
// @Tags 文档库详情
// @Accept json
// @Produce json
// @Param library_id query int true "文档库 ID"
// @Success 200 {object} model.HttpSuccess{data=[]model.DocLibraryDetailTreeResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library-detail/tree [get]
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
