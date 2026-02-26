package file

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/infra/storage"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/file"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/storagestrategy"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileService    file.FileService
	storageService storagestrategy.StorageStrategyService
}

func NewFileHandler(fileService file.FileService, storageService storagestrategy.StorageStrategyService) *FileHandler {
	return &FileHandler{fileService: fileService, storageService: storageService}
}

// @Summary 获取文件列表
// @Description 获取所有文件的列表
// @Tags 后台管理接口/文件
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]model.FileResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/file/list [get]
func (h *FileHandler) ListFile(c *fiber.Ctx) error {
	files, err := h.fileService.ListFile(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fileResps := make([]*model.FileResp, 0, len(files))
	for _, file := range files {
		fileResps = append(fileResps, &model.FileResp{
			ID:                file.ID,
			CreatedAt:         model.LocalTime(file.CreatedAt),
			Name:              file.Name,
			Path:              file.Path,
			URL:               file.URL,
			Type:              file.Type,
			Size:              file.Size,
			StorageStrategyID: file.StorageStrategyID,
			StorageStrategy: func() *string {
				if file.Edges.StorageStrategy != nil {
					return &file.Edges.StorageStrategy.Name
				}
				return nil
			}(),
		})
	}
	return c.JSON(model.NewSuccess("文件列表获取成功", fileResps))
}

// @Summary 获取文件分页列表
// @Description 获取所有文件的分页列表
// @Tags 后台管理接口/文件
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param name query string false "文件名称"
// @Param type query string false "文件类型"
// @Param storage_strategy_id query int false "存储策略ID"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.FileResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/file/page [get]
func (h *FileHandler) ListFilePage(c *fiber.Ctx) error {
	var pageReq model.FilePageReq
	if err := c.QueryParser(&pageReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, files, err := h.fileService.ListFilePageWithQuery(c.Context(), pageReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fileResps := make([]*model.FileResp, 0, len(files))
	for _, file := range files {
		fileResps = append(fileResps, &model.FileResp{
			ID:                file.ID,
			CreatedAt:         model.LocalTime(file.CreatedAt),
			Name:              file.Name,
			Path:              file.Path,
			URL:               file.URL,
			Type:              file.Type,
			Size:              file.Size,
			StorageStrategyID: file.StorageStrategyID,
			StorageStrategy: func() *string {
				if file.Edges.StorageStrategy != nil {
					return &file.Edges.StorageStrategy.Name
				}
				return nil
			}(),
		})
	}

	pageResult := model.PageResult[*model.FileResp]{
		Total:   int64(count),
		Records: fileResps,
	}
	return c.JSON(model.NewSuccess("文件列表获取成功", pageResult))
}

// @Summary 查询文件
// @Description 查询指定文件的信息
// @Tags 后台管理接口/文件
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} model.HttpSuccess{data=model.FileResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/file/query/{id} [get]
func (h *FileHandler) QueryFile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	file, err := h.fileService.QueryFile(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fileResp := &model.FileResp{
		ID:                file.ID,
		CreatedAt:         model.LocalTime(file.CreatedAt),
		Name:              file.Name,
		Path:              file.Path,
		URL:               file.URL,
		Type:              file.Type,
		Size:              file.Size,
		StorageStrategyID: file.StorageStrategyID,
		StorageStrategy: func() *string {
			if file.Edges.StorageStrategy != nil {
				return &file.Edges.StorageStrategy.Name
			}
			return nil
		}(),
	}
	return c.JSON(model.NewSuccess("文件查询成功", fileResp))
}

// @Summary 删除文件
// @Description 删除指定文件
// @Tags 后台管理接口/文件
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/file/delete/{id} [delete]
func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	err = h.fileService.DeleteFile(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("文件删除成功", nil))
}

// @Summary 上传文件
// @Description 上传文件到指定的存储策略
// @Tags 后台管理接口/文件
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Param storage_strategy formData int false "存储策略ID" default(1)
// @Success 200 {object} model.HttpSuccess{data=[]ent.File}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/file/upload [post]
func (h *FileHandler) Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	files := form.File["file"]

	if len(files) == 0 {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "请选择要上传的文件"))
	}

	var storageStrategyID int

	// 检查存储器策略
	if formValue, ok := form.Value["storage_strategy"]; ok {
		if len(formValue) > 0 {
			storageStrategyID, err = strconv.Atoi(formValue[0])
			if err != nil {
				return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的存储策略ID"))
			}
		}
	}

	var results []*ent.File

	for _, file := range files {

		// 获取存储策略
		strategy, err := h.storageService.GetStorageStrategyByID(c.Context(), storageStrategyID)
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
		}
		uploader, err := storage.GetUploader(strategy)
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
		}
		// 上传文件
		// 打开文件获取 io.Reader
		f, err := file.Open()
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		defer f.Close()

		url, err := uploader.Upload(file.Filename, f, file.Size, file.Header.Get("Content-Type"))
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		// 保存文件
		var fullUrl string

		if strategy.Domain != "" {
			fullUrl = strategy.Domain + "/" + url
		} else {
			fullUrl = strategy.Endpoint + "/" + url
		}

		// 保存到数据库
		newFile, err := h.fileService.CreateFile(c.Context(), strategy.ID, file.Filename, strategy.BasePath, fullUrl, file.Header.Get("Content-Type"), strconv.FormatInt(file.Size, 10))
		if err != nil {
			return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		results = append(results, newFile)
	}

	return c.JSON(model.NewSuccess("文件上传成功", results))
}
