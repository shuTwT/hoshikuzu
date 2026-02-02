package albumphoto

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	albumphotoservice "github.com/shuTwT/hoshikuzu/internal/services/content/albumphoto"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type AlbumPhotoHandler interface {
	ListAlbumPhoto(c *fiber.Ctx) error
	ListAlbumPhotoPage(c *fiber.Ctx) error
	CreateAlbumPhoto(c *fiber.Ctx) error
	UpdateAlbumPhoto(c *fiber.Ctx) error
	QueryAlbumPhoto(c *fiber.Ctx) error
	DeleteAlbumPhoto(c *fiber.Ctx) error
}

type AlbumPhotoHandlerImpl struct {
	albumPhotoService albumphotoservice.AlbumPhotoService
}

func NewAlbumPhotoHandlerImpl(albumPhotoService albumphotoservice.AlbumPhotoService) *AlbumPhotoHandlerImpl {
	return &AlbumPhotoHandlerImpl{
		albumPhotoService: albumPhotoService,
	}
}

// @Summary 查询相册照片列表
// @Description 查询所有相册照片
// @Tags albumPhotos
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.AlbumPhoto}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album-photo/list [get]
func (h *AlbumPhotoHandlerImpl) ListAlbumPhoto(c *fiber.Ctx) error {
	albumPhotos, err := h.albumPhotoService.ListAlbumPhoto(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", albumPhotos))
}

// @Summary 查询相册照片分页列表
// @Description 查询所有相册照片分页列表
// @Tags albumPhotos
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.AlbumPhoto]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album-photo/page [get]
func (h *AlbumPhotoHandlerImpl) ListAlbumPhotoPage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, albumPhotos, err := h.albumPhotoService.ListAlbumPhotoPage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	pageResult := model.PageResult[*ent.AlbumPhoto]{
		Total:   int64(count),
		Records: albumPhotos,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建相册照片
// @Description 创建一个新的相册照片
// @Tags albumPhotos
// @Accept json
// @Produce json
// @Param album_photo body model.AlbumPhotoCreateReq true "相册照片信息"
// @Success 200 {object} model.HttpSuccess{data=model.AlbumPhotoCreateReq}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album-photo/create [post]
func (h *AlbumPhotoHandlerImpl) CreateAlbumPhoto(c *fiber.Ctx) error {
	var albumPhoto *model.AlbumPhotoCreateReq
	if err := c.BodyParser(&albumPhoto); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	newAlbumPhoto, err := h.albumPhotoService.CreateAlbumPhoto(c.Context(), albumPhoto.AlbumID, albumPhoto.Name, albumPhoto.ImageURL, albumPhoto.Description)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", newAlbumPhoto))
}

// @Summary 更新相册照片
// @Description 更新指定相册照片的信息
// @Tags albumPhotos
// @Accept json
// @Produce json
// @Param id path string true "相册照片ID"
// @Param album_photo body model.AlbumPhotoUpdateReq true "相册照片信息"
// @Success 200 {object} model.HttpSuccess{data=model.AlbumPhotoUpdateReq}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album-photo/update/{id} [put]
func (h *AlbumPhotoHandlerImpl) UpdateAlbumPhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	var albumPhoto *model.AlbumPhotoUpdateReq
	if err := c.BodyParser(&albumPhoto); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedAlbumPhoto, err := h.albumPhotoService.UpdateAlbumPhoto(c.Context(), id, albumPhoto.Name, albumPhoto.ImageURL, albumPhoto.Description)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", updatedAlbumPhoto))
}

// @Summary 查询相册照片
// @Description 查询指定相册照片的信息
// @Tags albumPhotos
// @Accept json
// @Produce json
// @Param id path string true "相册照片ID"
// @Success 200 {object} model.HttpSuccess{data=ent.AlbumPhoto}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album-photo/query/{id} [get]
func (h *AlbumPhotoHandlerImpl) QueryAlbumPhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	albumPhoto, err := h.albumPhotoService.QueryAlbumPhoto(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", albumPhoto))
}

// @Summary 删除相册照片
// @Description 删除指定相册照片
// @Tags albumPhotos
// @Accept json
// @Produce json
// @Param id path string true "相册照片ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/album-photo/delete/{id} [delete]
func (h *AlbumPhotoHandlerImpl) DeleteAlbumPhoto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	err = h.albumPhotoService.DeleteAlbumPhoto(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}
