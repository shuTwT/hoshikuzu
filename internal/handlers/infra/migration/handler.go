package migration

import (
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	migration_service "github.com/shuTwT/hoshikuzu/internal/services/infra/migration"

	"github.com/gofiber/fiber/v2"
)

type MigrationHandler interface {
	ImportMarkdown(c *fiber.Ctx) error
	CheckDuplicate(c *fiber.Ctx) error
}

type MigrationHandlerImpl struct {
	migrationService migration_service.MigrationService
}

func NewMigrationHandlerImpl(migrationService migration_service.MigrationService) *MigrationHandlerImpl {
	return &MigrationHandlerImpl{
		migrationService: migrationService,
	}
}

// @Summary 导入Markdown文件
// @Description 批量导入Markdown文件到文章管理
// @Tags 后台管理接口/迁移
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Markdown文件"
// @Success 200 {object} model.HttpSuccess{data=model.MigrationResult}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/migration/md [post]
func (h *MigrationHandlerImpl) ImportMarkdown(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无法解析表单数据"))
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "请选择要导入的文件"))
	}

	result, err := h.migrationService.ImportMarkdownFiles(c.Context(), files)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if result.Failed > 0 {
		result.Status = "partial"
		result.Message = "部分文件导入成功"
	} else {
		result.Status = "success"
		result.Message = "所有文件导入成功"
	}

	return c.JSON(model.NewSuccess("导入完成", result))
}

// @Summary 检查重复文件
// @Description 检查待导入的文件中是否有与数据库中文章标题重复的文件
// @Tags 后台管理接口/迁移
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Markdown文件"
// @Success 200 {object} model.HttpSuccess{data=model.MigrationCheckResult}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/migration/check-duplicate [post]
func (h *MigrationHandlerImpl) CheckDuplicate(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无法解析表单数据"))
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "请选择要导入的文件"))
	}

	result, err := h.migrationService.CheckDuplicateFiles(c.Context(), files)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("检查完成", result))
}
