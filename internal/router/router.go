package router

import (
	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/handlers"
	"github.com/shuTwT/hoshikuzu/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// 注册系统路由
func initSystemRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	settingsApi := router.Group("/settings")
	{
		settingsApi.Get("/json/:key", handlerMap.SettingHandler.GetJsonSettingsMap)
		settingsApi.Post("/json/save/:key", handlerMap.SettingHandler.SaveSettings)
	}
	roleApi := router.Group("/role")
	{
		roleApi.Get("/list", handlerMap.RoleHandler.ListRole)
		roleApi.Get("/page", handlerMap.RoleHandler.ListRolePage)
		roleApi.Post("/create", handlerMap.RoleHandler.CreateRole)
		roleApi.Put("/update/:id", handlerMap.RoleHandler.UpdateRole)
		roleApi.Get("/query/:id", handlerMap.RoleHandler.QueryRole)
		roleApi.Delete("/delete/:id", handlerMap.RoleHandler.DeleteRole)
	}
	userApi := router.Group("/user")
	{
		userApi.Get("/profile", handlerMap.UserHandler.GetUserProfile)
		userApi.Get("/list", handlerMap.UserHandler.ListUser)
		userApi.Get("/page", handlerMap.UserHandler.ListUserPage)
		userApi.Post("/create", handlerMap.UserHandler.CreateUser)
		userApi.Put("/update/:id", handlerMap.UserHandler.UpdateUser)
		userApi.Get("/query/:id", handlerMap.UserHandler.QueryUser)
		userApi.Delete("/delete/:id", handlerMap.UserHandler.DeleteUser)
	}
	notificationApi := router.Group("/notifications")
	{
		notificationApi.Get("/page", handlerMap.NotificationHandler.ListNotificationPage)
		notificationApi.Get("/query/:id", handlerMap.NotificationHandler.QueryNotification)
		notificationApi.Delete("/delete/:id", handlerMap.NotificationHandler.DeleteNotification)
		notificationApi.Post("/batch/read", handlerMap.NotificationHandler.BatchMarkAsRead)
	}
}

// 注册内容路由
func initContentRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	commentApi := router.Group("/comment")
	{
		commentApi.Get("/page", handlerMap.CommentHandler.ListCommentPage)
		commentApi.Get("/query/:id", handlerMap.CommentHandler.GetComment)
	}
	albumApi := router.Group("/album")
	{
		albumApi.Get("/list", handlerMap.AlbumHandler.ListAlbum)
		albumApi.Get("/page", handlerMap.AlbumHandler.ListAlbumPage)
		albumApi.Post("/create", handlerMap.AlbumHandler.CreateAlbum)
		albumApi.Put("/update/:id", handlerMap.AlbumHandler.UpdateAlbum)
		albumApi.Get("/query/:id", handlerMap.AlbumHandler.QueryAlbum)
		albumApi.Delete("/delete/:id", handlerMap.AlbumHandler.DeleteAlbum)
	}
	albumPhotoApi := router.Group("/album-photo")
	{
		albumPhotoApi.Get("/list", handlerMap.AlbumPhotoHandler.ListAlbumPhoto)
		albumPhotoApi.Get("/page", handlerMap.AlbumPhotoHandler.ListAlbumPhotoPage)
		albumPhotoApi.Post("/create", handlerMap.AlbumPhotoHandler.CreateAlbumPhoto)
		albumPhotoApi.Put("/update/:id", handlerMap.AlbumPhotoHandler.UpdateAlbumPhoto)
		albumPhotoApi.Get("/query/:id", handlerMap.AlbumPhotoHandler.QueryAlbumPhoto)
		albumPhotoApi.Delete("/delete/:id", handlerMap.AlbumPhotoHandler.DeleteAlbumPhoto)
	}
	flinkApi := router.Group("/flink")
	{
		flinkApi.Get("/list", handlerMap.FlinkHandler.ListFlink)
		flinkApi.Get("/page", handlerMap.FlinkHandler.ListFlinkPage)
		flinkApi.Post("/create", handlerMap.FlinkHandler.CreateFlink)
		flinkApi.Put("/update/:id", handlerMap.FlinkHandler.UpdateFlink)
		flinkApi.Get("/query/:id", handlerMap.FlinkHandler.QueryFlink)
		flinkApi.Delete("/delete/:id", handlerMap.FlinkHandler.DeleteFlink)
	}
	flinkGroupApi := router.Group("/flink-group")
	{
		flinkGroupApi.Get("/list", handlerMap.FlinkGroupHandler.ListFLinkGroup)
		flinkGroupApi.Post("/create", handlerMap.FlinkGroupHandler.CreateFlinkGroup)
		flinkGroupApi.Put("/update/:id", handlerMap.FlinkGroupHandler.UpdateFlinkGroup)
		flinkGroupApi.Delete("/delete/:id", handlerMap.FlinkGroupHandler.DeleteFLinkGroup)
	}
	friendCircleRecordApi := router.Group("/friend-circle")
	{
		friendCircleRecordApi.Get("/page", handlerMap.FriendCircleHandler.ListFriendCircleRecordPage)
	}
	essayApi := router.Group("/essay")
	{
		essayApi.Get("/list", handlerMap.EssayHandler.ListEssay)
		essayApi.Get("/page", handlerMap.EssayHandler.GetEssayPage)
		essayApi.Post("/create", handlerMap.EssayHandler.CreateEssay)
		essayApi.Put("/update/:id", handlerMap.EssayHandler.UpdateEssay)
		essayApi.Delete("/delete/:id", handlerMap.EssayHandler.DeleteEssay)
	}
	postApi := router.Group("/post")
	{
		postApi.Get("/list", handlerMap.PostHandler.ListPost)
		postApi.Get("/page", handlerMap.PostHandler.ListPostPage)
		postApi.Post("/create", handlerMap.PostHandler.CreatePost)
		postApi.Put("/update/content/:id", handlerMap.PostHandler.UpdatePostContent)
		postApi.Put("/update/setting/:id", handlerMap.PostHandler.UpdatePostSetting)
		postApi.Put("/publish/:id", handlerMap.PostHandler.PublishPost)
		postApi.Put("/unpublish/:id", handlerMap.PostHandler.UnpublishPost)
		postApi.Get("/query/:id", handlerMap.PostHandler.QueryPost)
		postApi.Delete("/delete/:id", handlerMap.PostHandler.DeletePost)

	}
	categoryApi := router.Group("/category")
	{
		categoryApi.Get("/query/:id", handlerMap.CategoryHandler.QueryCategory)
		categoryApi.Get("/list", handlerMap.CategoryHandler.QueryCategoryList)
		categoryApi.Get("/page", handlerMap.CategoryHandler.QueryCategoryPage)
		categoryApi.Post("/create", handlerMap.CategoryHandler.CreateCategory)
		categoryApi.Put("/update/:id", handlerMap.CategoryHandler.UpdateCategory)
		categoryApi.Delete("/delete/:id", handlerMap.CategoryHandler.DeleteCategory)
	}
	tagApi := router.Group("/tag")
	{
		tagApi.Get("/query/:id", handlerMap.TagHandler.QueryTag)
		tagApi.Get("/list", handlerMap.TagHandler.QueryTagList)
		tagApi.Get("/page", handlerMap.TagHandler.QueryTagPage)
		tagApi.Post("/create", handlerMap.TagHandler.CreateTag)
		tagApi.Put("/update/:id", handlerMap.TagHandler.UpdateTag)
		tagApi.Delete("/delete/:id", handlerMap.TagHandler.DeleteTag)
	}
	docLibraryApi := router.Group("/doclibrary")
	{
		docLibraryApi.Get("/list", handlerMap.DocLibraryHandler.GetDocLibraryList)
		docLibraryApi.Get("/page", handlerMap.DocLibraryHandler.GetDocLibraryPage)
		docLibraryApi.Post("/create", handlerMap.DocLibraryHandler.CreateDocLibrary)
		docLibraryApi.Put("/update/:id", handlerMap.DocLibraryHandler.UpdateDocLibrary)
		docLibraryApi.Get("/query/:id", handlerMap.DocLibraryHandler.GetDocLibrary)
		docLibraryApi.Delete("/delete/:id", handlerMap.DocLibraryHandler.DeleteDocLibrary)
	}
	docLibraryDetailApi := router.Group("/doclibrarydetail")
	{
		docLibraryDetailApi.Get("/page", handlerMap.DocLibraryDetailHandler.GetDocLibraryDetailPage)
		docLibraryDetailApi.Post("/create", handlerMap.DocLibraryDetailHandler.CreateDocLibraryDetail)
		docLibraryDetailApi.Put("/update/:id", handlerMap.DocLibraryDetailHandler.UpdateDocLibraryDetail)
		docLibraryDetailApi.Get("/query/:id", handlerMap.DocLibraryDetailHandler.GetDocLibraryDetail)
		docLibraryDetailApi.Delete("/delete/:id", handlerMap.DocLibraryDetailHandler.DeleteDocLibraryDetail)
		docLibraryDetailApi.Get("/tree", handlerMap.DocLibraryDetailHandler.GetDocLibraryDetailTree)
	}
	knowledgeBaseApi := router.Group("/knowledgebase")
	{
		knowledgeBaseApi.Get("/list", handlerMap.KnowledgeBaseHandler.GetKnowledgeBaseList)
		knowledgeBaseApi.Get("/page", handlerMap.KnowledgeBaseHandler.GetKnowledgeBasePage)
		knowledgeBaseApi.Post("/create", handlerMap.KnowledgeBaseHandler.CreateKnowledgeBase)
		knowledgeBaseApi.Put("/update/:id", handlerMap.KnowledgeBaseHandler.UpdateKnowledgeBase)
		knowledgeBaseApi.Get("/query/:id", handlerMap.KnowledgeBaseHandler.GetKnowledgeBase)
		knowledgeBaseApi.Delete("/delete/:id", handlerMap.KnowledgeBaseHandler.DeleteKnowledgeBase)
	}
	flinkApplicationApi := router.Group("/flink-application")
	{
		flinkApplicationApi.Get("/page", handlerMap.FlinkApplicationHandler.ListFlinkApplicationPage)
		flinkApplicationApi.Get("/query/:id", handlerMap.FlinkApplicationHandler.QueryFlinkApplication)
		flinkApplicationApi.Put("/update/:id", handlerMap.FlinkApplicationHandler.ApproveFlinkApplication)
	}
}

// 注册基础设施路由
func initInfraRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	storageStrategyApi := router.Group("/storage-strategy")
	{
		storageStrategyApi.Get("/list", handlerMap.StorageStrategyHandler.ListStorageStrategy)
		storageStrategyApi.Get("/list-all", handlerMap.StorageStrategyHandler.ListStorageStrategyAll)
		storageStrategyApi.Post("/create", handlerMap.StorageStrategyHandler.CreateStorageStrategy)
		storageStrategyApi.Put("/update/:id", handlerMap.StorageStrategyHandler.UpdateStorageStrategy)
		storageStrategyApi.Get("/query/:id", handlerMap.StorageStrategyHandler.QueryStorageStrategy)
		storageStrategyApi.Delete("/delete/:id", handlerMap.StorageStrategyHandler.DeleteStorageStrategy)
		storageStrategyApi.Put("/default/:id", handlerMap.StorageStrategyHandler.SetDefaultStorageStrategy)
	}
	fileApi := router.Group("/file")
	{
		fileApi.Get("/list", handlerMap.FileHandler.ListFile)
		fileApi.Get("/page", handlerMap.FileHandler.ListFilePage)
		fileApi.Get("/query/:id", handlerMap.FileHandler.QueryFile)
		fileApi.Delete("/delete/:id", handlerMap.FileHandler.DeleteFile)
		fileApi.Post("/upload", handlerMap.FileHandler.Upload)
	}
	scheduleJobApi := router.Group("/schedule-job")
	{
		scheduleJobApi.Post("/create", handlerMap.ScheduleJobHandler.CreateScheduleJob)
		scheduleJobApi.Get("/page", handlerMap.ScheduleJobHandler.ListScheduleJobPage)
		scheduleJobApi.Get("/query/:id", handlerMap.ScheduleJobHandler.QueryScheduleJob)
		scheduleJobApi.Put("/update/:id", handlerMap.ScheduleJobHandler.UpdateScheduleJob)
		scheduleJobApi.Delete("/delete/:id", handlerMap.ScheduleJobHandler.DeleteScheduleJob)
		scheduleJobApi.Post("/execute/:id", handlerMap.ScheduleJobHandler.ExecuteScheduleJobNow)
	}
	migrationApi := router.Group("/migration")
	{
		migrationApi.Post("/md", handlerMap.MigrationHandler.ImportMarkdown)
		migrationApi.Post("/check-duplicate", handlerMap.MigrationHandler.CheckDuplicate)
	}
	visitLogApi := router.Group("/visit-log")
	{
		visitLogApi.Get("/page", handlerMap.VisitHandler.ListVisitLogPage)
		visitLogApi.Get("/query/:id", handlerMap.VisitHandler.QueryVisitLog)
		visitLogApi.Delete("/delete/:id", handlerMap.VisitHandler.DeleteVisitLog)
		visitLogApi.Post("/batch/delete", handlerMap.VisitHandler.BatchDeleteVisitLog)
	}
	pluginApi := router.Group("/plugin")
	{
		pluginApi.Post("/create", handlerMap.PluginHandler.CreatePlugin)
		pluginApi.Get("/page", handlerMap.PluginHandler.ListPluginPage)
		pluginApi.Get("/query/:id", handlerMap.PluginHandler.QueryPlugin)
		pluginApi.Delete("/delete/:id", handlerMap.PluginHandler.DeletePlugin)
		pluginApi.Post("/:id/start", handlerMap.PluginHandler.StartPlugin)
		pluginApi.Post("/:id/stop", handlerMap.PluginHandler.StopPlugin)
		pluginApi.Post("/:id/restart", handlerMap.PluginHandler.RestartPlugin)
	}
	themeApi := router.Group("/theme")
	{
		themeApi.Post("/upload", handlerMap.ThemeHandler.UploadThemeFile)
		themeApi.Post("/create", handlerMap.ThemeHandler.CreateTheme)
		themeApi.Get("/page", handlerMap.ThemeHandler.ListThemePage)
		themeApi.Get("/query/:id", handlerMap.ThemeHandler.QueryTheme)
		themeApi.Delete("/delete/:id", handlerMap.ThemeHandler.DeleteTheme)
		themeApi.Post("/:id/enable", handlerMap.ThemeHandler.EnableTheme)
		themeApi.Post("/:id/disable", handlerMap.ThemeHandler.DisableTheme)
	}
	licenseApi := router.Group("/license")
	{
		licenseApi.Get("/page", handlerMap.LicenseHandler.ListLicensePage)
		licenseApi.Get("/query/:id", handlerMap.LicenseHandler.QueryLicense)
		licenseApi.Post("/create", handlerMap.LicenseHandler.CreateLicense)
		licenseApi.Put("/update/:id", handlerMap.LicenseHandler.UpdateLicense)
		licenseApi.Delete("/delete/:id", handlerMap.LicenseHandler.DeleteLicense)
		licenseApi.Post("/verify", handlerMap.LicenseHandler.VerifyLicense)
	}
}

// 注册商城路由
func initMallRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	productApi := router.Group("/product")
	{
		productApi.Get("/list", handlerMap.ProductHandler.ListProducts).Name("productList")
		productApi.Get("/page", handlerMap.ProductHandler.ListProductsPage).Name("productPage")
		productApi.Post("/create", handlerMap.ProductHandler.CreateProduct).Name("productCreate")
		productApi.Put("/update/:id", handlerMap.ProductHandler.UpdateProduct).Name("productUpdate")
		productApi.Get("/query/:id", handlerMap.ProductHandler.QueryProduct).Name("productQuery")
		productApi.Delete("/delete/:id", handlerMap.ProductHandler.DeleteProduct).Name("productDelete")
		productApi.Put("/batch", handlerMap.ProductHandler.BatchUpdateProducts).Name("productBatchUpdate")
		productApi.Post("/batch/delete", handlerMap.ProductHandler.BatchDeleteProducts).Name("productBatchDelete")
	}
	memberApi := router.Group("/member")
	{
		memberApi.Get("/query/:user_id", handlerMap.MemberHandler.QueryMember).Name("memberQuery")
		memberApi.Get("/page", handlerMap.MemberHandler.QueryMemberPage).Name("memberPage")
		memberApi.Post("/create", handlerMap.MemberHandler.CreateMember).Name("memberCreate")
		memberApi.Put("/update/:id", handlerMap.MemberHandler.UpdateMember).Name("memberUpdate")
		memberApi.Delete("/delete/:id", handlerMap.MemberHandler.DeleteMember).Name("memberDelete")
	}
	memberLevelApi := router.Group("/member-level")
	{
		memberLevelApi.Get("/query/:id", handlerMap.MemberLevelHandler.QueryMemberLevel).Name("memberLevelQuery")
		memberLevelApi.Get("/list", handlerMap.MemberLevelHandler.QueryMemberLevelList).Name("memberLevelList")
		memberLevelApi.Get("/page", handlerMap.MemberLevelHandler.QueryMemberLevelPage).Name("memberLevelPage")
		memberLevelApi.Post("/create", handlerMap.MemberLevelHandler.CreateMemberLevel).Name("memberLevelCreate")
		memberLevelApi.Put("/update/:id", handlerMap.MemberLevelHandler.UpdateMemberLevel).Name("memberLevelUpdate")
		memberLevelApi.Delete("/delete/:id", handlerMap.MemberLevelHandler.DeleteMemberLevel).Name("memberLevelDelete")
	}
	payOrderApi := router.Group("/pay-order")
	{
		payOrderApi.Get("/page", handlerMap.PayOrderHandler.ListPayOrderPage).Name("payOrderPage")
		payOrderApi.Put("/update/:id", handlerMap.PayOrderHandler.UpdatePayOrder).Name("payOrderUpdate")
		payOrderApi.Get("/query/:id", handlerMap.PayOrderHandler.QueryPayOrder).Name("payOrderQuery")
		payOrderApi.Delete("/delete/:id", handlerMap.PayOrderHandler.DeletePayOrder).Name("payOrderDelete")
		payOrderApi.Post("/submit", handlerMap.PayOrderHandler.SubmitPayOrder).Name("payOrderSubmit")
		payOrderApi.Get("/today-stats", handlerMap.PayOrderHandler.GetTodayStats).Name("payOrderTodayStats")
	}

	walletApi := router.Group("/wallet")
	{
		walletApi.Get("/query/:user_id", handlerMap.WalletHandler.QueryWallet).Name("walletQuery")
		walletApi.Get("/page", handlerMap.WalletHandler.QueryWalletPage).Name("walletPage")
		walletApi.Put("/update/:id", handlerMap.WalletHandler.UpdateWallet).Name("walletUpdate")
	}
	couponApi := router.Group("/coupon")
	{
		couponApi.Get("/list", handlerMap.CouponHandler.ListCoupons).Name("couponList")
		couponApi.Get("/page", handlerMap.CouponHandler.ListCouponsPage).Name("couponPage")
		couponApi.Post("/create", handlerMap.CouponHandler.CreateCoupon).Name("couponCreate")
		couponApi.Put("/update/:id", handlerMap.CouponHandler.UpdateCoupon).Name("couponUpdate")
		couponApi.Get("/query/:id", handlerMap.CouponHandler.QueryCoupon).Name("couponQuery")
		couponApi.Delete("/delete/:id", handlerMap.CouponHandler.DeleteCoupon).Name("couponDelete")
		couponApi.Put("/batch", handlerMap.CouponHandler.BatchUpdateCoupons).Name("couponBatchUpdate")
		couponApi.Post("/batch/delete", handlerMap.CouponHandler.BatchDeleteCoupons).Name("couponBatchDelete")
		couponApi.Get("/search", handlerMap.CouponHandler.SearchCoupons).Name("couponSearch")
	}
	couponUsageApi := router.Group("/coupon-usage")
	{
		couponUsageApi.Get("/list", handlerMap.CouponUsageHandler.ListCouponUsages).Name("couponUsageList")
		couponUsageApi.Get("/page", handlerMap.CouponUsageHandler.ListCouponUsagesPage).Name("couponUsagePage")
		couponUsageApi.Post("/create", handlerMap.CouponUsageHandler.CreateCouponUsage).Name("couponUsageCreate")
		couponUsageApi.Put("/update/:id", handlerMap.CouponUsageHandler.UpdateCouponUsage).Name("couponUsageUpdate")
		couponUsageApi.Get("/query/:id", handlerMap.CouponUsageHandler.QueryCouponUsage).Name("couponUsageQuery")
		couponUsageApi.Delete("/delete/:id", handlerMap.CouponUsageHandler.DeleteCouponUsage).Name("couponUsageDelete")
		couponUsageApi.Put("/batch", handlerMap.CouponUsageHandler.BatchUpdateCouponUsages).Name("couponUsageBatchUpdate")
		couponUsageApi.Post("/batch/delete", handlerMap.CouponUsageHandler.BatchDeleteCouponUsages).Name("couponUsageBatchDelete")
		couponUsageApi.Get("/search", handlerMap.CouponUsageHandler.SearchCouponUsages).Name("couponUsageSearch")
	}
}

// 公开路由无需认证
func initPublicRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	publicApi := router.Group("/public")
	{
		// 访问统计接口
		publicApi.Post("/visit", handlerMap.VisitHandler.HandleVisitor)
		// twikoo 接口
		publicApi.All("/twikoo", handlerMap.CommentHandler.HandleTwikoo)
		// 最近评论接口
		publicApi.Get("/comment/recent", handlerMap.CommentHandler.RecentComment)
		// 相册列表接口
		publicApi.Get("/album/list", handlerMap.AlbumHandler.ListAlbum)
		// 相册分页接口
		publicApi.Get("/album/page", handlerMap.AlbumHandler.ListAlbumPage)
		// 相册照片列表接口
		publicApi.Get("/album-photo/list", handlerMap.AlbumPhotoHandler.ListAlbumPhoto)
		// 相册照片分页接口
		publicApi.Get("/album-photo/page", handlerMap.AlbumPhotoHandler.ListAlbumPhotoPage)
		// 友链列表接口
		publicApi.Get("/flink/list", handlerMap.FlinkHandler.ListFlink)
		// 友链分页接口
		publicApi.Get("/flink/page", handlerMap.FlinkHandler.ListFlinkPage)
		// 友链分组列表接口
		publicApi.Get("/flink-group/list", handlerMap.FlinkGroupHandler.ListFLinkGroup)
		// 朋友圈记录分页接口
		publicApi.Get("/friend-circle-record/page", handlerMap.FriendCircleHandler.ListFriendCircleRecordPage)
		// 说说(随笔、瞬间)列表接口
		publicApi.Get("/essay/list", handlerMap.EssayHandler.ListEssay)
		// 说说(随笔、瞬间)分页接口
		publicApi.Get("/essay/page", handlerMap.EssayHandler.GetEssayPage)
		// 文章列表接口
		publicApi.Get("/post/list", handlerMap.PostHandler.ListPost)
		// 文章分页接口
		publicApi.Get("/post/page", handlerMap.PostHandler.ListPostPage)
		// 文章搜索接口
		publicApi.Get("/post/search", handlerMap.PostHandler.SearchPosts)
		// 文章月统计接口
		publicApi.Get("/post/month-stats", handlerMap.PostHandler.GetPostMonthStats)
		// 随机文章接口
		publicApi.Get("/post/random", handlerMap.PostHandler.GetRandomPost)
		// 文章详情接口
		publicApi.Get("/post/slug/:slug", handlerMap.PostHandler.QueryPostBySlug)
		// 分类列表接口
		publicApi.Get("/category/list", handlerMap.CategoryHandler.QueryCategoryList)
		// 分类分页接口
		publicApi.Get("/category/page", handlerMap.CategoryHandler.QueryCategoryPage)
		// 标签列表接口
		publicApi.Get("/tag/list", handlerMap.TagHandler.QueryTagList)
		// 标签分页接口
		publicApi.Get("/tag/page", handlerMap.TagHandler.QueryTagPage)
		// 用户搜索接口
		publicApi.Get("/user/search", handlerMap.UserHandler.SearchUsers)
		// 随机友链接口
		publicApi.Get("/flink/random", handlerMap.FlinkHandler.RandomFlink)
		// 文章摘要接口
		publicApi.Get("/post/:id/summary/stream", handlerMap.PostHandler.GetSummaryForStream)
		// 商品搜索接口
		publicApi.Get("/product/search", handlerMap.ProductHandler.SearchProducts)
		// 友链申请接口
		publicApi.Post("/flink-application/create", handlerMap.FlinkApplicationHandler.CreateFlinkApplication)
		// 插件注册接口，仅在debug模式下生效，不需要认证
		publicApi.Post("/plugin/register", handlerMap.PluginHandler.RegisterPlugin)
		// 插件心跳接口，仅在debug模式下生效，不需要认证
		publicApi.Post("/plugin/heartbeat", handlerMap.PluginHandler.HeartbeatPlugin)
	}
}

func Initialize(router *fiber.App, handlerMap handlers.HandlerMap, dbClient *ent.Client) {
	router.Use(middleware.Security)
	router.Get("/api/preinit", handlerMap.InitializeHandler.PreInit)
	router.Post("/api/initialize", handlerMap.InitializeHandler.Initialize)

	auth := router.Group("/api/auth")
	{
		auth.Post("/login/password", handlerMap.AuthHandler.Login)
	}

	api := router.Group("/api")
	initPublicRouter(api, handlerMap)
	{

		apiV1 := api.Group("/v1")
		{

			// 路由列表接口
			apiV1.Get("/routes", handlerMap.RouteHandler.GetRoutes)
			apiV1.Get("/api-interface/page", handlerMap.ApiInterfaceHandler.ListApiRoutesPage)
			apiV1.Get("/settings", handlerMap.SettingHandler.GetSettings)

			apiV1.Use(middleware.FlexibleAuth(dbClient))

			// 首页统计信息接口
			apiV1.Get("/common/statistic", handlerMap.CommonHandler.GetHomeStatistics)

			apiV1.Get("/user/personal-access-token/list", handlerMap.UserHandler.GetPersonalAccessTokenList)
			apiV1.Get("/user/personal-access-token/query/:id", handlerMap.UserHandler.GetPersonalAccessToken)
			apiV1.Post("/user/personal-access-token/create", handlerMap.UserHandler.CreatePat)

			initContentRouter(apiV1, handlerMap)
			initInfraRouter(apiV1, handlerMap)
			initMallRouter(apiV1, handlerMap)
			initSystemRouter(apiV1, handlerMap)
		}
	}
}
