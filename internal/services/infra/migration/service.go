package migration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/post"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type MigrationService interface {
	ImportMarkdownFiles(ctx context.Context, files []*multipart.FileHeader) (*model.MigrationResult, error)
	CheckDuplicateFiles(ctx context.Context, files []*multipart.FileHeader) (*model.MigrationCheckResult, error)
}

type MigrationServiceImpl struct {
	client        *ent.Client
	markdown      goldmark.Markdown
	htmlSanitizer *bluemonday.Policy
}

func NewMigrationServiceImpl(client *ent.Client) *MigrationServiceImpl {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	policy := bluemonday.UGCPolicy()
	policy.AllowElements("a", "abbr", "b", "blockquote", "br", "caption", "cite", "code", "col", "colgroup",
		"dd", "del", "details", "div", "dl", "dt", "em", "figcaption", "figure", "h1", "h2", "h3", "h4", "h5", "h6",
		"hr", "i", "img", "ins", "kbd", "li", "ol", "p", "pre", "q", "rp", "rt", "ruby", "s", "samp",
		"section", "small", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead",
		"tr", "u", "ul", "var", "wbr")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`[\w\s\-_]+`)).OnElements("div", "span", "p", "td", "th", "li")
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowAttrs("src", "alt", "title", "width", "height").OnElements("img")
	policy.AllowAttrs("colspan", "rowspan").OnElements("td", "th")

	return &MigrationServiceImpl{
		client:        client,
		markdown:      md,
		htmlSanitizer: policy,
	}
}

func (s *MigrationServiceImpl) ImportMarkdownFiles(ctx context.Context, files []*multipart.FileHeader) (*model.MigrationResult, error) {
	result := &model.MigrationResult{
		Total:   len(files),
		Success: 0,
		Failed:  0,
		Errors:  []string{},
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: 无法打开文件 - %v", fileHeader.Filename, err))
			continue
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: 读取文件失败 - %v", fileHeader.Filename, err))
			continue
		}

		title := extractTitleFromMarkdown(fileHeader.Filename, string(content))

		var buf bytes.Buffer
		if err = s.markdown.Convert(content, &buf); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Markdown转HTML失败 - %v", fileHeader.Filename, err))
			continue
		}

		htmlContent := s.htmlSanitizer.Sanitize(buf.String())

		existingPost, _ := s.client.Post.Query().
			Where(post.Title(title)).
			First(ctx)

		if existingPost != nil {
			_, updateErr := s.client.Post.UpdateOneID(existingPost.ID).
				SetContent(htmlContent).
				SetMdContent(string(content)).
				SetHTMLContent(htmlContent).
				Save(ctx)

			if updateErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: 更新文章失败 - %v", fileHeader.Filename, updateErr))
				continue
			}
		} else {
			_, err = s.client.Post.Create().
				SetTitle(title).
				SetContent(htmlContent).
				SetMdContent(string(content)).
				SetHTMLContent(htmlContent).
				SetContentType("markdown").
				SetStatus("draft").
				SetAuthor("匿名作者").
				Save(ctx)
		}

		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: 保存到数据库失败 - %v", fileHeader.Filename, err))
			continue
		}

		result.Success++
	}

	return result, nil
}

func (s *MigrationServiceImpl) CheckDuplicateFiles(ctx context.Context, files []*multipart.FileHeader) (*model.MigrationCheckResult, error) {
	duplicates := []model.DuplicateFile{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			continue
		}

		title := extractTitleFromMarkdown(fileHeader.Filename, string(content))

		existingPost, err := s.client.Post.Query().
			Where(post.Title(title)).
			First(ctx)

		if err == nil && existingPost != nil {
			duplicates = append(duplicates, model.DuplicateFile{
				Filename: fileHeader.Filename,
				Title:    title,
				PostID:   &existingPost.ID,
			})
		}
	}

	return &model.MigrationCheckResult{
		HasDuplicates: len(duplicates) > 0,
		Duplicates:    duplicates,
	}, nil
}

func extractTitleFromMarkdown(filename, content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
