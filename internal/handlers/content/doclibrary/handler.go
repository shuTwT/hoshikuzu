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

// @Summary 创建文档库
// @Description 创建新的文档库
// @Tags 文档库
// @Accept json
// @Produce json
// @Param doc_library_create_req body model.DocLibraryCreateReq true "文档库创建请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library/create [post]
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

// @Summary 更新文档库
// @Description 更新指定文档库的信息
// @Tags 文档库
// @Accept json
// @Produce json
// @Param id path int true "文档库 ID"
// @Param doc_library_update_req body model.DocLibraryUpdateReq true "文档库更新请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library/update/{id} [put]
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

// @Summary 获取文档库分页列表
// @Description 获取文档库的分页列表
// @Tags 文档库
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.DocLibraryResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library/page [get]
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

// @Summary 获取文档库
// @Description 获取指定文档库的信息
// @Tags 文档库
// @Accept json
// @Produce json
// @Param id path int true "文档库 ID"
// @Success 200 {object} model.HttpSuccess{data=model.DocLibraryResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library/query/{id} [get]
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

// @Summary 删除文档库
// @Description 删除指定文档库
// @Tags 文档库
// @Accept json
// @Produce json
// @Param id path int true "文档库 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library/delete/{id} [delete]
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

// @Summary 获取文档库列表
// @Description 获取所有文档库的列表
// @Tags 文档库
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.DocLibraryResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/doc-library/list [get]
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
