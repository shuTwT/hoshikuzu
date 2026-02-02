package theme

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	theme_ent "github.com/shuTwT/hoshikuzu/ent/theme"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"gopkg.in/yaml.v3"
)

type ThemeService interface {
	ListThemePage(ctx context.Context, page, size int) (int, []*ent.Theme, error)
	QueryTheme(ctx context.Context, id int) (*ent.Theme, error)
	QueryThemeByName(ctx context.Context, name string) (*ent.Theme, error)
	UploadThemeFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
	CreateTheme(ctx context.Context, req *model.CreateThemeReq) (*ent.Theme, error)
	DeleteTheme(ctx context.Context, id int) error
	EnableTheme(ctx context.Context, id int) error
	DisableTheme(ctx context.Context, id int) error
	RegisterDefaultTheme(ctx context.Context) error
	GetEnabledTheme(ctx context.Context) (*ent.Theme, error)
}

type ThemeServiceImpl struct {
	client *ent.Client
}

func NewThemeServiceImpl(client *ent.Client) *ThemeServiceImpl {
	return &ThemeServiceImpl{
		client: client,
	}
}

func (s *ThemeServiceImpl) ListThemePage(ctx context.Context, page, size int) (int, []*ent.Theme, error) {
	count, err := s.client.Theme.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	themes, err := s.client.Theme.Query().
		Order(ent.Desc(theme_ent.FieldID)).
		Offset((page - 1) * size).
		Limit(size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, themes, nil
}

func (s *ThemeServiceImpl) QueryTheme(ctx context.Context, id int) (*ent.Theme, error) {
	themeEntity, err := s.client.Theme.Query().
		Where(theme_ent.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return themeEntity, nil
}

func (s *ThemeServiceImpl) QueryThemeByName(ctx context.Context, name string) (*ent.Theme, error) {
	themeEntity, err := s.client.Theme.Query().
		Where(theme_ent.Name(name)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return themeEntity, nil
}

func (s *ThemeServiceImpl) UploadThemeFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", errors.New("文件不能为空")
	}

	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer srcFile.Close()

	tmpDir := "./data/tmp/themes"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	fileName := fmt.Sprintf("theme-%d.zip", time.Now().UnixNano())
	filePath := filepath.Join(tmpDir, fileName)

	dstFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("复制文件失败: %w", err)
	}

	return filePath, nil
}

func (s *ThemeServiceImpl) CreateTheme(ctx context.Context, req *model.CreateThemeReq) (*ent.Theme, error) {
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}

	if req.Type == "internal" {
		return s.createInternalTheme(ctx, req)
	} else if req.Type == "external" {
		return s.createExternalTheme(ctx, req)
	}

	return nil, errors.New("无效的主题类型")
}

func (s *ThemeServiceImpl) createInternalTheme(ctx context.Context, req *model.CreateThemeReq) (*ent.Theme, error) {
	if req.FilePath == "" {
		return nil, errors.New("文件路径不能为空")
	}

	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		return nil, errors.New("文件不存在")
	}

	zipReader, err := zip.OpenReader(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("打开压缩包失败: %w", err)
	}
	defer zipReader.Close()

	var themeConfigContent []byte
	var settingConfigContent []byte
	themeDir := ""

	for _, f := range zipReader.File {
		if f.Name == "theme.yaml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("打开theme.yaml文件失败: %w", err)
			}
			themeConfigContent, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("读取theme.yaml文件失败: %w", err)
			}
		} else if f.Name == "setting.yaml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("打开setting.yaml文件失败: %w", err)
			}
			settingConfigContent, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("读取setting.yaml文件失败: %w", err)
			}
		} else if strings.HasSuffix(f.Name, "/") {
			if themeDir == "" {
				themeDir = strings.TrimSuffix(f.Name, "/")
			}
		}
	}

	if themeConfigContent == nil {
		return nil, errors.New("压缩包中未找到 theme.yaml 文件")
	}

	if settingConfigContent == nil {
		return nil, errors.New("压缩包中未找到 setting.yaml 文件")
	}

	var themeConfig model.ThemeConfig
	if err := yaml.Unmarshal(themeConfigContent, &themeConfig); err != nil {
		return nil, fmt.Errorf("解析theme.yaml文件失败: %w", err)
	}

	if err := validateThemeConfig(&themeConfig); err != nil {
		return nil, err
	}

	exists, err := s.client.Theme.Query().Where(theme_ent.Name(themeConfig.Name)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("检查主题是否存在失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("主题 '%s' 已存在", themeConfig.Name)
	}

	themesDir := "./data/themes"
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return nil, fmt.Errorf("创建主题目录失败: %w", err)
	}

	targetDir := filepath.Join(themesDir, themeConfig.Name)
	if err := os.RemoveAll(targetDir); err != nil {
		return nil, fmt.Errorf("清理旧主题目录失败: %w", err)
	}

	for _, f := range zipReader.File {
		targetPath := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return nil, fmt.Errorf("创建目录失败: %w", err)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return nil, fmt.Errorf("创建父目录失败: %w", err)
			}
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("打开压缩文件失败: %w", err)
			}
			defer rc.Close()

			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return nil, fmt.Errorf("创建文件失败: %w", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, rc); err != nil {
				return nil, fmt.Errorf("解压文件失败: %w", err)
			}
		}
	}

	os.Remove(req.FilePath)

	builder := s.client.Theme.Create().
		SetType("internal").
		SetName(themeConfig.Name).
		SetDisplayName(themeConfig.DisplayName).
		SetVersion(themeConfig.Version).
		SetRequire(themeConfig.Require).
		SetPath(targetDir).
		SetEnabled(false)

	if themeConfig.Description != "" {
		builder.SetDescription(themeConfig.Description)
	}
	if themeConfig.Author != nil {
		if themeConfig.Author.Name != "" {
			builder.SetAuthorName(themeConfig.Author.Name)
		}
		if themeConfig.Author.Email != "" {
			builder.SetAuthorEmail(themeConfig.Author.Email)
		}
	}
	if themeConfig.Logo != "" {
		builder.SetLogo(themeConfig.Logo)
	}
	if themeConfig.Homepage != "" {
		builder.SetHomepage(themeConfig.Homepage)
	}
	if themeConfig.Repo != "" {
		builder.SetRepo(themeConfig.Repo)
	}
	if themeConfig.Issue != "" {
		builder.SetIssue(themeConfig.Issue)
	}
	if themeConfig.SettingName != "" {
		builder.SetSettingName(themeConfig.SettingName)
	}
	if themeConfig.ConfigMapName != "" {
		builder.SetConfigMapName(themeConfig.ConfigMapName)
	}
	if themeConfig.License != "" {
		builder.SetLicense(themeConfig.License)
	}

	themeEntity, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("保存主题信息失败: %w", err)
	}

	return themeEntity, nil
}

func (s *ThemeServiceImpl) createExternalTheme(ctx context.Context, req *model.CreateThemeReq) (*ent.Theme, error) {
	if req.Name == "" {
		return nil, errors.New("主题名称不能为空")
	}
	if req.DisplayName == "" {
		return nil, errors.New("显示名称不能为空")
	}
	if req.ExternalURL == "" {
		return nil, errors.New("外部主题URL不能为空")
	}
	if req.Version == "" {
		return nil, errors.New("版本号不能为空")
	}

	exists, err := s.client.Theme.Query().Where(theme_ent.Name(req.Name)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("检查主题是否存在失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("主题 '%s' 已存在", req.Name)
	}

	builder := s.client.Theme.Create().
		SetType("external").
		SetName(req.Name).
		SetDisplayName(req.DisplayName).
		SetExternalURL(req.ExternalURL).
		SetVersion(req.Version).
		SetEnabled(false)

	if req.Description != "" {
		builder.SetDescription(req.Description)
	}

	themeEntity, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("保存外部主题信息失败: %w", err)
	}

	return themeEntity, nil
}

func (s *ThemeServiceImpl) DeleteTheme(ctx context.Context, id int) error {
	themeEntity, err := s.client.Theme.Query().Where(theme_ent.ID(id)).First(ctx)
	if err != nil {
		return err
	}

	if themeEntity.Enabled {
		return errors.New("主题已启用，无法删除")
	}

	themeDir := themeEntity.Path
	if err := os.RemoveAll(themeDir); err != nil {
		return fmt.Errorf("删除主题目录失败: %w", err)
	}

	err = s.client.Theme.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ThemeServiceImpl) EnableTheme(ctx context.Context, id int) error {
	themeEntity, err := s.client.Theme.Query().Where(theme_ent.ID(id)).First(ctx)
	if err != nil {
		return err
	}

	if themeEntity.Enabled {
		return errors.New("主题已启用")
	}

	err = s.client.Theme.Update().
		SetEnabled(false).
		Where(theme_ent.Enabled(true)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("禁用其他主题失败: %w", err)
	}

	err = s.client.Theme.UpdateOneID(id).
		SetEnabled(true).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ThemeServiceImpl) DisableTheme(ctx context.Context, id int) error {
	themeEntity, err := s.client.Theme.Query().Where(theme_ent.ID(id)).First(ctx)
	if err != nil {
		return err
	}

	if !themeEntity.Enabled {
		return errors.New("主题未启用")
	}

	err = s.client.Theme.UpdateOneID(id).
		SetEnabled(false).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ThemeServiceImpl) GetEnabledTheme(ctx context.Context) (*ent.Theme, error) {
	themeEntity, err := s.client.Theme.Query().
		Where(theme_ent.Enabled(true)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return themeEntity, nil
}

func validateThemeConfig(config *model.ThemeConfig) error {
	if config.Name == "" {
		return errors.New("主题名称不能为空")
	}
	if config.DisplayName == "" {
		return errors.New("显示名称不能为空")
	}
	if config.Version == "" {
		return errors.New("版本号不能为空")
	}
	if config.Type != "theme" {
		return errors.New("类型必须为 theme")
	}
	return nil
}

func (s *ThemeServiceImpl) RegisterDefaultTheme(ctx context.Context) error {
	themeDir := "./data/themes/hoshikuzu-theme-ace"
	themeConfigPath := filepath.Join(themeDir, "theme.yaml")

	// Check if theme.yaml exists
	if _, err := os.Stat(themeConfigPath); os.IsNotExist(err) {
		// If file not found, it might mean extraction failed or path is wrong.
		// Since extraction is done before this call, we should log error but not panic?
		// But here we return error.
		return fmt.Errorf("default theme config not found at %s", themeConfigPath)
	}

	// Read theme.yaml
	themeConfigContent, err := os.ReadFile(themeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read theme config: %w", err)
	}

	var themeConfig model.ThemeConfig
	if err := yaml.Unmarshal(themeConfigContent, &themeConfig); err != nil {
		return fmt.Errorf("failed to parse theme config: %w", err)
	}

	// Check if theme exists in DB
	exists, err := s.client.Theme.Query().Where(theme_ent.Name(themeConfig.Name)).Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if theme exists: %w", err)
	}

	if !exists {
		builder := s.client.Theme.Create().
			SetType("internal").
			SetName(themeConfig.Name).
			SetDisplayName(themeConfig.DisplayName).
			SetVersion(themeConfig.Version).
			SetRequire(themeConfig.Require).
			SetPath(themeDir).
			SetEnabled(false)

		if themeConfig.Description != "" {
			builder.SetDescription(themeConfig.Description)
		}
		if themeConfig.Author != nil {
			if themeConfig.Author.Name != "" {
				builder.SetAuthorName(themeConfig.Author.Name)
			}
			if themeConfig.Author.Email != "" {
				builder.SetAuthorEmail(themeConfig.Author.Email)
			}
		}
		if themeConfig.Logo != "" {
			builder.SetLogo(themeConfig.Logo)
		}
		if themeConfig.Homepage != "" {
			builder.SetHomepage(themeConfig.Homepage)
		}
		if themeConfig.Repo != "" {
			builder.SetRepo(themeConfig.Repo)
		}
		if themeConfig.Issue != "" {
			builder.SetIssue(themeConfig.Issue)
		}
		if themeConfig.SettingName != "" {
			builder.SetSettingName(themeConfig.SettingName)
		}
		if themeConfig.ConfigMapName != "" {
			builder.SetConfigMapName(themeConfig.ConfigMapName)
		}
		if themeConfig.License != "" {
			builder.SetLicense(themeConfig.License)
		}

		_, err := builder.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to save default theme to DB: %w", err)
		}
	}

	// Check if we need to enable it
	// If it's the only theme, enable it.
	// Or if no theme is enabled, enable it.

	count, err := s.client.Theme.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count themes: %w", err)
	}

	enabledCount, err := s.client.Theme.Query().Where(theme_ent.Enabled(true)).Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count enabled themes: %w", err)
	}

	if count == 1 || enabledCount == 0 {
		// Enable this theme
		// We need the ID first.
		themeEntity, err := s.client.Theme.Query().Where(theme_ent.Name(themeConfig.Name)).First(ctx)
		if err != nil {
			return fmt.Errorf("failed to query default theme: %w", err)
		}

		// Only enable if not already enabled (though enabledCount==0 implies it's not)
		if !themeEntity.Enabled {
			// Use internal logic to enable to ensure consistency (e.g. disable others)
			// But here we know others are disabled or don't exist.
			// Still, using EnableTheme is safer if we change logic later.
			return s.EnableTheme(ctx, themeEntity.ID)
		}
	}

	return nil
}
