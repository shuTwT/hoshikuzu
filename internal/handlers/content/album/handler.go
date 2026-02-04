package album

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/content/album"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type AlbumHandler interface {
	ListAlbum(c *fiber.Ctx) error
	ListAlbumPage(c *fiber.Ctx) error
	CreateAlbum(c *fiber.Ctx) error
	UpdateAlbum(c *fiber.Ctx) error
	QueryAlbum(c *fiber.Ctx) error
	DeleteAlbum(c *fiber.Ctx) error
}

type AlbumHandlerImpl struct {
	albumService album.AlbumService
}

func NewAlbumHandlerImpl(albumService album.AlbumService) *AlbumHandlerImpl {
	return &AlbumHandlerImpl{
		albumService: albumService,
	}
}

// @Summary 查询相册列表
// @Description 查询所有相册
// @Tags 后台管理接口/相册
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Album}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album/list [get]
func (h *AlbumHandlerImpl) ListAlbum(c *fiber.Ctx) error {
	albums, err := h.albumService.ListAlbum(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", albums))
}

// @Summary 查询相册列表分页
// @Description 查询相册列表分页
// @Tags 后台管理接口/相册
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Album]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album/page [get]
func (h *AlbumHandlerImpl) ListAlbumPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, albums, err := h.albumService.ListAlbumPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	pageResult := model.PageResult[*ent.Album]{
		Total:   int64(count),
		Records: albums,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建相册
// @Description 创建一个新相册
// @Tags 后台管理接口/相册
// @Accept json
// @Produce json
// @Param req body model.AlbumCreateReq true "相册创建请求"
// @Success 200 {object} model.HttpSuccess{data=model.AlbumCreateReq}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album/create [post]
func (h *AlbumHandlerImpl) CreateAlbum(c *fiber.Ctx) error {
	var album *model.AlbumCreateReq
	if err := c.BodyParser(&album); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	newAlbum, err := h.albumService.CreateAlbum(c.Context(), album)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", newAlbum))
}

// @Summary 更新相册
// @Description 更新指定相册的信息
// @Tags 后台管理接口/相册
// @Accept json
// @Produce json
// @Param id path string true "相册ID"
// @Param req body model.AlbumUpdateReq true "相册更新请求"
// @Success 200 {object} model.HttpSuccess{data=model.AlbumUpdateReq}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album/update/{id} [put]
func (h *AlbumHandlerImpl) UpdateAlbum(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	var album *model.AlbumUpdateReq
	if err := c.BodyParser(&album); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedAlbum, err := h.albumService.UpdateAlbum(c.Context(), id, album)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", updatedAlbum))
}

// @Summary 查询相册
// @Description 查询指定相册的信息
// @Tags 后台管理接口/相册
// @Accept json
// @Produce json
// @Param id path string true "相册ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Album}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album/query/{id} [get]
func (h *AlbumHandlerImpl) QueryAlbum(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	album, err := h.albumService.QueryAlbum(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", album))
}

// @Summary 删除相册
// @Description 删除指定相册
// @Tags 后台管理接口/相册
// @Accept json
// @Produce json
// @Param id path string true "相册ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album/delete/{id} [delete]
func (h *AlbumHandlerImpl) DeleteAlbum(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	err = h.albumService.DeleteAlbum(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
