package theme

import (
	"log/slog"
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/theme"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type ThemeHandler interface {
	UploadThemeFile(c *fiber.Ctx) error
	CreateTheme(c *fiber.Ctx) error
	ListThemePage(c *fiber.Ctx) error
	QueryTheme(c *fiber.Ctx) error
	DeleteTheme(c *fiber.Ctx) error
	EnableTheme(c *fiber.Ctx) error
	DisableTheme(c *fiber.Ctx) error
}

type ThemeHandlerImpl struct {
	themeService theme.ThemeService
}

func NewThemeHandlerImpl(themeService theme.ThemeService) *ThemeHandlerImpl {
	return &ThemeHandlerImpl{themeService: themeService}
}

// @Summary 上传主题文件
// @Description 上传一个新的主题文件
// @Tags 主题
// @Accept json
// @Produce json
// @Param file formData file true "主题文件"
// @Success 200 {object} model.HttpSuccess{data=map[string]string}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/upload [post]
func (h *ThemeHandlerImpl) UploadThemeFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		slog.Error("Failed to get uploaded file", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "获取上传文件失败"))
	}

	filePath, err := h.themeService.UploadThemeFile(c.Context(), file)
	if err != nil {
		slog.Error("Failed to upload theme file", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Theme file uploaded successfully", "file_path", filePath)
	return c.JSON(model.NewSuccess("主题文件上传成功", map[string]string{"file_path": filePath}))
}

// @Summary 创建主题
// @Description 创建一个新的主题
// @Tags 主题
// @Accept json
// @Produce json
// @Param req body model.CreateThemeReq true "主题创建请求"
// @Success 200 {object} model.HttpSuccess{data=model.ThemeResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/create [post]
func (h *ThemeHandlerImpl) CreateTheme(c *fiber.Ctx) error {
	var req model.CreateThemeReq
	if err := c.BodyParser(&req); err != nil {
		slog.Error("Failed to parse request body", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "请求参数解析失败"))
	}

	themeEntity, err := h.themeService.CreateTheme(c.Context(), &req)
	if err != nil {
		slog.Error("Failed to create theme", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildThemeResp(themeEntity)
	slog.Info("Theme created successfully", "theme_id", themeEntity.ID, "theme_name", themeEntity.Name)
	return c.JSON(model.NewSuccess("主题创建成功", resp))
}

// @Summary 查询主题列表
// @Description 查询所有主题的分页列表
// @Tags 主题
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.ThemeResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/page [get]
func (h *ThemeHandlerImpl) ListThemePage(c *fiber.Ctx) error {
	var pageQuery model.PageQuery
	if err := c.QueryParser(&pageQuery); err != nil {
		slog.Error("Failed to parse query parameters", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "查询参数解析失败"))
	}

	count, themes, err := h.themeService.ListThemePage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		slog.Error("Failed to list themes", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	themeResps := make([]*model.ThemeResp, 0, len(themes))
	for _, t := range themes {
		themeResps = append(themeResps, h.buildThemeResp(t))
	}

	pageResult := model.PageResult[*model.ThemeResp]{
		Total:   int64(count),
		Records: themeResps,
	}
	return c.JSON(model.NewSuccess("主题列表获取成功", pageResult))
}

// @Summary 查询主题
// @Description 查询指定主题的信息
// @Tags 主题
// @Accept json
// @Produce json
// @Param id path int true "主题ID"
// @Success 200 {object} model.HttpSuccess{data=model.ThemeResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/query/{id} [get]
func (h *ThemeHandlerImpl) QueryTheme(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid theme ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的主题ID"))
	}

	themeEntity, err := h.themeService.QueryTheme(c.Context(), id)
	if err != nil {
		slog.Error("Failed to query theme", "theme_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	resp := h.buildThemeResp(themeEntity)
	return c.JSON(model.NewSuccess("主题查询成功", resp))
}

// @Summary 删除主题
// @Description 删除指定主题
// @Tags 主题
// @Accept json
// @Produce json
// @Param id path int true "主题ID"
// @Success 200 {object} model.HttpSuccess{data=model.ThemeResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/delete/{id} [delete]
func (h *ThemeHandlerImpl) DeleteTheme(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid theme ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的主题ID"))
	}

	err = h.themeService.DeleteTheme(c.Context(), id)
	if err != nil {
		slog.Error("Failed to delete theme", "theme_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Theme deleted successfully", "theme_id", id)
	return c.JSON(model.NewSuccess("主题删除成功", nil))
}

// @Summary 启用主题
// @Description 启用指定主题
// @Tags 主题
// @Accept json
// @Produce json
// @Param id path int true "主题ID"
// @Success 200 {object} model.HttpSuccess{data=model.ThemeResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/enable/{id} [put]
func (h *ThemeHandlerImpl) EnableTheme(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid theme ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的主题ID"))
	}

	err = h.themeService.EnableTheme(c.Context(), id)
	if err != nil {
		slog.Error("Failed to enable theme", "theme_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Theme enabled successfully", "theme_id", id)
	return c.JSON(model.NewSuccess("主题启用成功", nil))
}

// @Summary 禁用主题
// @Description 禁用指定主题
// @Tags 主题
// @Accept json
// @Produce json
// @Param id path int true "主题ID"
// @Success 200 {object} model.HttpSuccess{data=model.ThemeResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/theme/disable/{id} [put]
func (h *ThemeHandlerImpl) DisableTheme(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		slog.Error("Invalid theme ID", "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的主题ID"))
	}

	err = h.themeService.DisableTheme(c.Context(), id)
	if err != nil {
		slog.Error("Failed to disable theme", "theme_id", id, "error", err.Error())
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	slog.Info("Theme disabled successfully", "theme_id", id)
	return c.JSON(model.NewSuccess("主题禁用成功", nil))
}

func (h *ThemeHandlerImpl) buildThemeResp(t *ent.Theme) *model.ThemeResp {
	return &model.ThemeResp{
		ID:            t.ID,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		Type:          t.Type,
		Name:          t.Name,
		DisplayName:   t.DisplayName,
		Description:   t.Description,
		AuthorName:    t.AuthorName,
		AuthorEmail:   t.AuthorEmail,
		Logo:          t.Logo,
		Homepage:      t.Homepage,
		Repo:          t.Repo,
		Issue:         t.Issue,
		SettingName:   t.SettingName,
		ConfigMapName: t.ConfigMapName,
		Version:       t.Version,
		Require:       t.Require,
		License:       t.License,
		Path:          t.Path,
		ExternalURL:   t.ExternalURL,
		Enabled:       t.Enabled,
	}
}
