package cmd

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/shuTwT/hoshikuzu/internal/handlers"
	"github.com/shuTwT/hoshikuzu/internal/infra/database"
	"github.com/shuTwT/hoshikuzu/internal/infra/schedule"
	"github.com/shuTwT/hoshikuzu/internal/infra/schedule/manager"
	"github.com/shuTwT/hoshikuzu/internal/router"
	"github.com/shuTwT/hoshikuzu/pkg"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

type (
	User struct {
		Name string `validate:"required,min=5,max=20"` // Required field, min 5 char long max 20
		Age  int    `validate:"required,teener"`       // Required field, and client needs to implement our 'teener' tag format which we'll see later
	}

	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

// This is the validator instance
// for more information see: https://github.com/go-playground/validator
var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func InitializeApp(assetsRes embed.FS, frontendRes embed.FS) (*fiber.App, func()) {
	config.Init()
	dbType := config.GetString(config.DATABASE_TYPE)
	dbConfig := database.DBConfig{
		DBType: dbType,
		DBUrl:  config.GetString(config.DATABASE_URL),
	}
	db, err := database.InitializeDB(dbConfig, true)
	if err != nil {
		panic(err)
	}

	// 初始化 Redis 客户端
	rdb, err := database.NewRedisClient(context.Background())

	cleanup := func() {
		err := db.Close()
		if err != nil {
			slog.Error("Failed to close database connection", "error", err.Error())
		}
		if rdb != nil {
			if err := rdb.Close(); err != nil {
				slog.Error("Failed to close Redis connection", "error", err.Error())
			}
		}
	}

	scheduleManager := manager.NewScheduleManager()

	pkg.ExtractDefaultTheme(assetsRes)

	serviceMap := pkg.InitializeServices(assetsRes, db, scheduleManager)

	if err := serviceMap.ThemeService.RegisterDefaultTheme(context.Background()); err != nil {
		slog.Error("Failed to register default theme", "error", err.Error())
	}

	if !fiber.IsChild() {
		// 主进程程初始化定时任务
		err := schedule.InitializeSchedule(db, scheduleManager, serviceMap.FriendCircleService)
		if err != nil {
			defer scheduleManager.Shutdown()
		}
	}

	// myValidator := &XValidator{
	// 	validator: validate,
	// }

	app := fiber.New(fiber.Config{
		AppName:       "Fiber HTML Template Demo",
		Prefork:       false,   // 启用多进程（Prefork 模式）
		CaseSensitive: true,    // 路由大小写敏感
		StrictRouting: true,    // 严格匹配带 / 和不带 / 的路由
		ServerHeader:  "Fiber", // 返回响应头中的 Server 字段
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(model.HttpError{
				Code: fiber.ErrBadRequest.Code,
				Msg:  err.Error(),
			})
		},
	})

	if config.GetBool(config.SWAGGER_ENABLE) {
		app.Get("/swagger/*", swagger.HandlerDefault) // default
	}

	// app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
	// 	URL:         "http://example.com/doc.json",
	// 	DeepLinking: false,
	// 	// Expand ("list") or Collapse ("none") tag groups by default
	// 	DocExpansion: "none",
	// 	// Prefill OAuth ClientId on Authorize popup
	// 	OAuth: &swagger.OAuthConfig{
	// 		AppName:  "OAuth Provider",
	// 		ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
	// 	},
	// 	// Ability to change OAuth2 redirect uri location
	// 	OAuth2RedirectUrl: "http://localhost:3000/swagger/oauth2-redirect.html",
	// }))
	app.Use(logger.New())
	app.Use(cors.New())

	handlerMap := handlers.InitHandler(serviceMap, db)
	router.InitFrontendRes(app, frontendRes, serviceMap)
	router.Initialize(app, handlerMap, db)

	go func() {
		if fiber.IsChild() {
			return
		}
		if err := serviceMap.PluginService.AutoStartPlugins(context.Background()); err != nil {
			fmt.Printf("自动启动插件失败: %v\n", err)
		}
		slog.Info("自动启动插件成功")
	}()

	return app, cleanup
}
