package router

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg"
)

func InitFrontendRes(app *fiber.App, frontendRes embed.FS, serviceMap pkg.ServiceMap) {

	distDir, err := fs.Sub(frontendRes, "ui/dist")
	if err != nil {
		log.Fatalln("静态资源载入失败")
	}

	app.Get("/console/*", func(c *fiber.Ctx) error {
		resp := c.Response()
		filePath := c.Params("*")
		var contentType string
		if file, err := distDir.Open(filePath); err == nil {
			defer file.Close()
			stat, err := file.Stat()
			if err == nil && !stat.IsDir() {
				ext := filepath.Ext(filePath)
				contentType = mime.TypeByExtension(ext)
				if contentType == "" {
					contentType = "application/octet-stream"
				}
				if ext == ".html" {

				}
				resp.Header.Set("Content-Type", contentType)
				_, err = io.Copy(c.Response().BodyWriter(), file)
				return err
			} else if err == nil && stat.IsDir() {
				fillePathTmp := filepath.Join(filePath, "/index.html")
				if fileSub, err := distDir.Open(fillePathTmp); err == nil {
					defer fileSub.Close()
					statSub, err := fileSub.Stat()
					if err == nil && !statSub.IsDir() {
						if filepath.Ext(filePath) == ".html" {
							resp.Header.Set("Content-Type", "text/html")
							_, err = io.Copy(c.Response().BodyWriter(), file)
							return err
						}
					}
				}
			}

		}
		if file, err := distDir.Open("index.html"); err == nil {
			defer file.Close()
			resp.Header.Set("Content-Type", "text/html")
			_, err = io.Copy(c.Response().BodyWriter(), file)
			return err
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	initFrontendRoutes(app, serviceMap)
}

func initFrontendRoutes(app *fiber.App, serviceMap pkg.ServiceMap) {
	ctx := context.Background()

	previewTheme := func(c *fiber.Ctx) (*ent.Theme, error) {
		previewName := c.Query("preview")
		if previewName == "" {
			return serviceMap.ThemeService.GetEnabledTheme(ctx)
		}

		themeName := "hoshikuzu-theme-" + previewName
		theme, err := serviceMap.ThemeService.QueryThemeByName(ctx, themeName)
		if err != nil {
			log.Printf("获取预览主题失败: %v", err)
			return serviceMap.ThemeService.GetEnabledTheme(ctx)
		}

		return theme, nil
	}

	app.Use(func(c *fiber.Ctx) error {
		theme, err := previewTheme(c)
		if err != nil {
			log.Printf("获取主题失败: %v", err)
			return c.Next()
		}

		if theme != nil && theme.Type == "external" {
			path := c.Path()
			if len(path) >= 4 && path[:4] == "/api" {
				return c.Next()
			}
			if len(path) >= 8 && path[:8] == "/console" {
				return c.Next()
			}
			if theme.ExternalURL != "" {
				return c.Redirect(theme.ExternalURL + path)
			}
			return c.Next()
		}

		c.Locals("theme", theme)
		return c.Next()
	})

	app.Get("/", renderTemplate(serviceMap, "index.html"))
	app.Get("/archives", renderTemplate(serviceMap, "archives.html"))
	app.Get("/author/:userId", renderTemplate(serviceMap, "author.html"))
	app.Get("/categories", renderTemplate(serviceMap, "categories.html"))
	app.Get("/category/:categoryName", renderTemplate(serviceMap, "category.html"))
	app.Get("/post/:slug", renderTemplate(serviceMap, "post.html"))
	app.Get("/tags", renderTemplate(serviceMap, "tags.html"))
	app.Get("/tag/:tagName", renderTemplate(serviceMap, "tag.html"))
	app.Get("/404", renderTemplate(serviceMap, "404.html"))

	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if len(path) >= 4 && path[:4] == "/api" {
			return c.Next()
		}
		if len(path) >= 8 && path[:8] == "/console" {
			return c.Next()
		}
		return c.Redirect("/404")
	})
}

func renderTemplate(serviceMap pkg.ServiceMap, templateName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		theme := c.Locals("theme")
		if theme == nil {
			ctx := context.Background()
			var err error
			theme, err = serviceMap.ThemeService.GetEnabledTheme(ctx)
			if err != nil {
				log.Printf("获取启用主题失败: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("获取主题失败")
			}
		}

		themeEntity := theme.(*ent.Theme)

		templatePath := filepath.Join(themeEntity.Path, "templates", templateName)

		log.Printf("%s", templatePath)

		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			log.Printf("模板文件不存在: %s", templatePath)
			return c.Status(fiber.StatusNotFound).SendString("模板文件不存在")
		}

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			log.Printf("解析模板失败: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("解析模板失败")
		}

		c.Set("Content-Type", "text/html; charset=utf-8")

		data := map[string]interface{}{
			"Title":  fmt.Sprintf("%s - %s", getTemplateTitle(templateName), themeEntity.DisplayName),
			"Theme":  themeEntity,
			"Path":   c.Path(),
			"Params": c.AllParams(),
		}

		if err := tmpl.Execute(c.Response().BodyWriter(), data); err != nil {
			log.Printf("渲染模板失败: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("渲染模板失败")
		}

		return nil
	}
}

func getTemplateTitle(templateName string) string {
	titles := map[string]string{
		"index.html":      "首页",
		"archives.html":   "归档",
		"author.html":     "作者",
		"categories.html": "分类",
		"category.html":   "分类",
		"post.html":       "文章",
		"tags.html":       "标签",
		"tag.html":        "标签",
		"404.html":        "404",
	}
	if title, ok := titles[templateName]; ok {
		return title
	}
	return "页面"
}
