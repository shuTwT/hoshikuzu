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
		settingsApi.Get("/json/:key", handlerMap.SettingHandler.GetJsonSettingsMap).Name("settingsMap")
		settingsApi.Post("/json/save/:key", handlerMap.SettingHandler.SaveSettings).Name("jsonSettingsSave")
	}
	roleApi := router.Group("/role")
	{
		roleApi.Get("/list", handlerMap.RoleHandler.ListRole).Name("roleList")
		roleApi.Get("/page", handlerMap.RoleHandler.ListRolePage).Name("rolePage")
		roleApi.Post("/create", handlerMap.RoleHandler.CreateRole).Name("roleCreate")
		roleApi.Put("/update/:id", handlerMap.RoleHandler.UpdateRole).Name("roleUpdate")
		roleApi.Get("/query/:id", handlerMap.RoleHandler.QueryRole).Name("roleQuery")
		roleApi.Delete("/delete/:id", handlerMap.RoleHandler.DeleteRole).Name("roleDelete")
	}
	userApi := router.Group("/user")
	{
		userApi.Get("/profile", handlerMap.UserHandler.GetUserProfile).Name("userProfile")
		userApi.Get("/list", handlerMap.UserHandler.ListUser).Name("userList")
		userApi.Get("/page", handlerMap.UserHandler.ListUserPage).Name("userPage")
		userApi.Post("/create", handlerMap.UserHandler.CreateUser).Name("userCreate")
		userApi.Put("/update/:id", handlerMap.UserHandler.UpdateUser).Name("userUpdate")
		userApi.Get("/query/:id", handlerMap.UserHandler.QueryUser).Name("userQuery")
		userApi.Delete("/delete/:id", handlerMap.UserHandler.DeleteUser).Name("userDelete")
	}
	notificationApi := router.Group("/notifications")
	{
		notificationApi.Get("/page", handlerMap.NotificationHandler.ListNotificationPage).Name("notificationPage")
		notificationApi.Get("/query/:id", handlerMap.NotificationHandler.QueryNotification).Name("notificationQuery")
		notificationApi.Delete("/delete/:id", handlerMap.NotificationHandler.DeleteNotification).Name("notificationDelete")
		notificationApi.Post("/batch/read", handlerMap.NotificationHandler.BatchMarkAsRead).Name("notificationBatchRead")
	}
}

// 注册内容路由
func initContentRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	commentApi := router.Group("/comment")
	{
		commentApi.Get("/page", handlerMap.CommentHandler.ListCommentPage).Name("commentPage")
		commentApi.Get("/query/:id", handlerMap.CommentHandler.GetComment).Name("commentQuery")
	}
	albumApi := router.Group("/album")
	{
		albumApi.Get("/list", handlerMap.AlbumHandler.ListAlbum).Name("albumList")
		albumApi.Get("/page", handlerMap.AlbumHandler.ListAlbumPage).Name("albumPage")
		albumApi.Post("/create", handlerMap.AlbumHandler.CreateAlbum).Name("albumCreate")
		albumApi.Put("/update/:id", handlerMap.AlbumHandler.UpdateAlbum).Name("albumUpdate")
		albumApi.Get("/query/:id", handlerMap.AlbumHandler.QueryAlbum).Name("albumQuery")
		albumApi.Delete("/delete/:id", handlerMap.AlbumHandler.DeleteAlbum).Name("albumDelete")
	}
	albumPhotoApi := router.Group("/album-photo")
	{
		albumPhotoApi.Get("/list", handlerMap.AlbumPhotoHandler.ListAlbumPhoto).Name("albumPhotoList")
		albumPhotoApi.Get("/page", handlerMap.AlbumPhotoHandler.ListAlbumPhotoPage).Name("albumPhotoPage")
		albumPhotoApi.Post("/create", handlerMap.AlbumPhotoHandler.CreateAlbumPhoto).Name("albumPhotoCreate")
		albumPhotoApi.Put("/update/:id", handlerMap.AlbumPhotoHandler.UpdateAlbumPhoto).Name("albumPhotoUpdate")
		albumPhotoApi.Get("/query/:id", handlerMap.AlbumPhotoHandler.QueryAlbumPhoto).Name("albumPhotoQuery")
		albumPhotoApi.Delete("/delete/:id", handlerMap.AlbumPhotoHandler.DeleteAlbumPhoto).Name("albumPhotoDelete")
	}
	flinkApi := router.Group("/flink")
	{
		flinkApi.Get("/list", handlerMap.FlinkHandler.ListFlink).Name("flinkList")
		flinkApi.Get("/page", handlerMap.FlinkHandler.ListFlinkPage).Name("flinkPage")
		flinkApi.Post("/create", handlerMap.FlinkHandler.CreateFlink).Name("flinkCreate")
		flinkApi.Put("/update/:id", handlerMap.FlinkHandler.UpdateFlink).Name("flinkUpdate")
		flinkApi.Get("/query/:id", handlerMap.FlinkHandler.QueryFlink).Name("flinkQuery")
		flinkApi.Delete("/delete/:id", handlerMap.FlinkHandler.DeleteFlink).Name("flinkDelete")
	}
	flinkGroupApi := router.Group("/flink-group")
	{
		flinkGroupApi.Get("/list", handlerMap.FlinkGroupHandler.ListFLinkGroup).Name("flinkGroupList")
		flinkGroupApi.Post("/create", handlerMap.FlinkGroupHandler.CreateFlinkGroup).Name("flinkGroupCreate")
		flinkGroupApi.Put("/update/:id", handlerMap.FlinkGroupHandler.UpdateFlinkGroup).Name("flinkGroupUpdate")
		flinkGroupApi.Delete("/delete/:id", handlerMap.FlinkGroupHandler.DeleteFLinkGroup).Name("flinkGroupDelete")
	}
	friendCircleRecordApi := router.Group("/friend-circle")
	{
		friendCircleRecordApi.Get("/page", handlerMap.FriendCircleHandler.ListFriendCircleRecordPage).Name("friendCircleRecordPage")
	}
	essayApi := router.Group("/essay")
	{
		essayApi.Get("/list", handlerMap.EssayHandler.ListEssay).Name("essayList")
		essayApi.Get("/page", handlerMap.EssayHandler.GetEssayPage).Name("essayPage")
		essayApi.Post("/create", handlerMap.EssayHandler.CreateEssay).Name("essayCreate")
		essayApi.Put("/update/:id", handlerMap.EssayHandler.UpdateEssay).Name("essayUpdate")
		essayApi.Delete("/delete/:id", handlerMap.EssayHandler.DeleteEssay).Name("essayDelete")
	}
	postApi := router.Group("/post")
	{
		postApi.Get("/list", handlerMap.PostHandler.ListPost).Name("postList")
		postApi.Get("/page", handlerMap.PostHandler.ListPostPage).Name("postPage")
		postApi.Get("/month-stats", handlerMap.PostHandler.GetPostMonthStats).Name("postMonthStats")
		postApi.Get("/random", handlerMap.PostHandler.GetRandomPost).Name("postRandom")
		postApi.Get("/slug/:slug", handlerMap.PostHandler.QueryPostBySlug).Name("postQueryBySlug")
		postApi.Post("/create", handlerMap.PostHandler.CreatePost).Name("postCreate")
		postApi.Put("/update/content/:id", handlerMap.PostHandler.UpdatePostContent).Name("postUpdateContent")
		postApi.Put("/update/setting/:id", handlerMap.PostHandler.UpdatePostSetting).Name("postUpdateSetting")
		postApi.Put("/publish/:id", handlerMap.PostHandler.PublishPost).Name("postPublish")
		postApi.Put("/unpublish/:id", handlerMap.PostHandler.UnpublishPost).Name("postUnpublish")
		postApi.Get("/query/:id", handlerMap.PostHandler.QueryPost).Name("postQuery")
		postApi.Delete("/delete/:id", handlerMap.PostHandler.DeletePost).Name("postDelete")

	}
	categoryApi := router.Group("/category")
	{
		categoryApi.Get("/query/:id", handlerMap.CategoryHandler.QueryCategory).Name("categoryQuery")
		categoryApi.Get("/list", handlerMap.CategoryHandler.QueryCategoryList).Name("categoryList")
		categoryApi.Get("/page", handlerMap.CategoryHandler.QueryCategoryPage).Name("categoryPage")
		categoryApi.Post("/create", handlerMap.CategoryHandler.CreateCategory).Name("categoryCreate")
		categoryApi.Put("/update/:id", handlerMap.CategoryHandler.UpdateCategory).Name("categoryUpdate")
		categoryApi.Delete("/delete/:id", handlerMap.CategoryHandler.DeleteCategory).Name("categoryDelete")
	}
	tagApi := router.Group("/tag")
	{
		tagApi.Get("/query/:id", handlerMap.TagHandler.QueryTag).Name("tagQuery")
		tagApi.Get("/list", handlerMap.TagHandler.QueryTagList).Name("tagList")
		tagApi.Get("/page", handlerMap.TagHandler.QueryTagPage).Name("tagPage")
		tagApi.Post("/create", handlerMap.TagHandler.CreateTag).Name("tagCreate")
		tagApi.Put("/update/:id", handlerMap.TagHandler.UpdateTag).Name("tagUpdate")
		tagApi.Delete("/delete/:id", handlerMap.TagHandler.DeleteTag).Name("tagDelete")
	}
	docLibraryApi := router.Group("/doclibrary")
	{
		docLibraryApi.Get("/list", handlerMap.DocLibraryHandler.GetDocLibraryList).Name("docLibraryList")
		docLibraryApi.Get("/page", handlerMap.DocLibraryHandler.GetDocLibraryPage).Name("docLibraryPage")
		docLibraryApi.Post("/create", handlerMap.DocLibraryHandler.CreateDocLibrary).Name("docLibraryCreate")
		docLibraryApi.Put("/update/:id", handlerMap.DocLibraryHandler.UpdateDocLibrary).Name("docLibraryUpdate")
		docLibraryApi.Get("/query/:id", handlerMap.DocLibraryHandler.GetDocLibrary).Name("docLibraryQuery")
		docLibraryApi.Delete("/delete/:id", handlerMap.DocLibraryHandler.DeleteDocLibrary).Name("docLibraryDelete")
	}
	docLibraryDetailApi := router.Group("/doclibrarydetail")
	{
		docLibraryDetailApi.Get("/page", handlerMap.DocLibraryDetailHandler.GetDocLibraryDetailPage).Name("docLibraryDetailPage")
		docLibraryDetailApi.Post("/create", handlerMap.DocLibraryDetailHandler.CreateDocLibraryDetail).Name("docLibraryDetailCreate")
		docLibraryDetailApi.Put("/update/:id", handlerMap.DocLibraryDetailHandler.UpdateDocLibraryDetail).Name("docLibraryDetailUpdate")
		docLibraryDetailApi.Get("/query/:id", handlerMap.DocLibraryDetailHandler.GetDocLibraryDetail).Name("docLibraryDetailQuery")
		docLibraryDetailApi.Delete("/delete/:id", handlerMap.DocLibraryDetailHandler.DeleteDocLibraryDetail).Name("docLibraryDetailDelete")
		docLibraryDetailApi.Get("/tree", handlerMap.DocLibraryDetailHandler.GetDocLibraryDetailTree).Name("docLibraryDetailTree")
	}
	knowledgeBaseApi := router.Group("/knowledgebase")
	{
		knowledgeBaseApi.Get("/list", handlerMap.KnowledgeBaseHandler.GetKnowledgeBaseList).Name("knowledgeBaseList")
		knowledgeBaseApi.Get("/page", handlerMap.KnowledgeBaseHandler.GetKnowledgeBasePage).Name("knowledgeBasePage")
		knowledgeBaseApi.Post("/create", handlerMap.KnowledgeBaseHandler.CreateKnowledgeBase).Name("knowledgeBaseCreate")
		knowledgeBaseApi.Put("/update/:id", handlerMap.KnowledgeBaseHandler.UpdateKnowledgeBase).Name("knowledgeBaseUpdate")
		knowledgeBaseApi.Get("/query/:id", handlerMap.KnowledgeBaseHandler.GetKnowledgeBase).Name("knowledgeBaseQuery")
		knowledgeBaseApi.Delete("/delete/:id", handlerMap.KnowledgeBaseHandler.DeleteKnowledgeBase).Name("knowledgeBaseDelete")
	}
	flinkApplicationApi := router.Group("/flink-application")
	{
		flinkApplicationApi.Get("/page", handlerMap.FlinkApplicationHandler.ListFlinkApplicationPage).Name("flinkApplicationPage")
		flinkApplicationApi.Get("/query/:id", handlerMap.FlinkApplicationHandler.QueryFlinkApplication).Name("flinkApplicationQuery")
		flinkApplicationApi.Put("/update/:id", handlerMap.FlinkApplicationHandler.ApproveFlinkApplication).Name("flinkApplicationUpdate")
	}
}

// 注册基础设施路由
func initInfraRouter(router fiber.Router, handlerMap handlers.HandlerMap) {
	storageStrategyApi := router.Group("/storage-strategy")
	{
		storageStrategyApi.Get("/list", handlerMap.StorageStrategyHandler.ListStorageStrategy).Name("storageStrategyList")
		storageStrategyApi.Get("/list-all", handlerMap.StorageStrategyHandler.ListStorageStrategyAll).Name("storageStrategyListAll")
		storageStrategyApi.Post("/create", handlerMap.StorageStrategyHandler.CreateStorageStrategy).Name("storageStrategyCreate")
		storageStrategyApi.Put("/update/:id", handlerMap.StorageStrategyHandler.UpdateStorageStrategy).Name("storageStrategyUpdate")
		storageStrategyApi.Get("/query/:id", handlerMap.StorageStrategyHandler.QueryStorageStrategy).Name("storageStrategyQuery")
		storageStrategyApi.Delete("/delete/:id", handlerMap.StorageStrategyHandler.DeleteStorageStrategy).Name("storageStrategyDelete")
		storageStrategyApi.Put("/default/:id", handlerMap.StorageStrategyHandler.SetDefaultStorageStrategy).Name("storageStrategyDefault")
	}
	fileApi := router.Group("/file")
	{
		fileApi.Get("/list", handlerMap.FileHandler.ListFile).Name("fileList")
		fileApi.Get("/page", handlerMap.FileHandler.ListFilePage).Name("filePage")
		fileApi.Get("/query/:id", handlerMap.FileHandler.QueryFile).Name("fileQuery")
		fileApi.Delete("/delete/:id", handlerMap.FileHandler.DeleteFile).Name("fileDelete")
		fileApi.Post("/upload", handlerMap.FileHandler.Upload).Name("fileUpload")
	}
	scheduleJobApi := router.Group("/schedule-job")
	{
		scheduleJobApi.Post("/create", handlerMap.ScheduleJobHandler.CreateScheduleJob).Name("scheduleJobCreate")
		scheduleJobApi.Get("/page", handlerMap.ScheduleJobHandler.ListScheduleJobPage).Name("scheduleJobPage")
		scheduleJobApi.Get("/query/:id", handlerMap.ScheduleJobHandler.QueryScheduleJob).Name("scheduleJobQuery")
		scheduleJobApi.Put("/update/:id", handlerMap.ScheduleJobHandler.UpdateScheduleJob).Name("scheduleJobUpdate")
		scheduleJobApi.Delete("/delete/:id", handlerMap.ScheduleJobHandler.DeleteScheduleJob).Name("scheduleJobDelete")
		scheduleJobApi.Post("/execute/:id", handlerMap.ScheduleJobHandler.ExecuteScheduleJobNow).Name("scheduleJobExecute")
	}
	migrationApi := router.Group("/migration")
	{
		migrationApi.Post("/md", handlerMap.MigrationHandler.ImportMarkdown).Name("migrationMd")
		migrationApi.Post("/check-duplicate", handlerMap.MigrationHandler.CheckDuplicate).Name("migrationCheckDuplicate")
	}
	visitLogApi := router.Group("/visit-log")
	{
		visitLogApi.Get("/page", handlerMap.VisitHandler.ListVisitLogPage).Name("visitLogPage")
		visitLogApi.Get("/query/:id", handlerMap.VisitHandler.QueryVisitLog).Name("visitLogQuery")
		visitLogApi.Delete("/delete/:id", handlerMap.VisitHandler.DeleteVisitLog).Name("visitLogDelete")
		visitLogApi.Post("/batch/delete", handlerMap.VisitHandler.BatchDeleteVisitLog).Name("visitLogBatchDelete")
	}
	pluginApi := router.Group("/plugin")
	{
		pluginApi.Post("/create", handlerMap.PluginHandler.CreatePlugin).Name("pluginCreate")
		pluginApi.Get("/page", handlerMap.PluginHandler.ListPluginPage).Name("pluginPage")
		pluginApi.Get("/query/:id", handlerMap.PluginHandler.QueryPlugin).Name("pluginQuery")
		pluginApi.Delete("/delete/:id", handlerMap.PluginHandler.DeletePlugin).Name("pluginDelete")
		pluginApi.Post("/:id/start", handlerMap.PluginHandler.StartPlugin).Name("pluginStart")
		pluginApi.Post("/:id/stop", handlerMap.PluginHandler.StopPlugin).Name("pluginStop")
		pluginApi.Post("/:id/restart", handlerMap.PluginHandler.RestartPlugin).Name("pluginRestart")
	}
	themeApi := router.Group("/theme")
	{
		themeApi.Post("/upload", handlerMap.ThemeHandler.UploadThemeFile).Name("themeUpload")
		themeApi.Post("/create", handlerMap.ThemeHandler.CreateTheme).Name("themeCreate")
		themeApi.Get("/page", handlerMap.ThemeHandler.ListThemePage).Name("themePage")
		themeApi.Get("/query/:id", handlerMap.ThemeHandler.QueryTheme).Name("themeQuery")
		themeApi.Delete("/delete/:id", handlerMap.ThemeHandler.DeleteTheme).Name("themeDelete")
		themeApi.Post("/:id/enable", handlerMap.ThemeHandler.EnableTheme).Name("themeEnable")
		themeApi.Post("/:id/disable", handlerMap.ThemeHandler.DisableTheme).Name("themeDisable")
	}
	licenseApi := router.Group("/license")
	{
		licenseApi.Get("/page", handlerMap.LicenseHandler.ListLicensePage).Name("licensePage")
		licenseApi.Get("/query/:id", handlerMap.LicenseHandler.QueryLicense).Name("licenseQuery")
		licenseApi.Post("/create", handlerMap.LicenseHandler.CreateLicense).Name("licenseCreate")
		licenseApi.Put("/update/:id", handlerMap.LicenseHandler.UpdateLicense).Name("licenseUpdate")
		licenseApi.Delete("/delete/:id", handlerMap.LicenseHandler.DeleteLicense).Name("licenseDelete")
		licenseApi.Post("/verify", handlerMap.LicenseHandler.VerifyLicense).Name("licenseVerify")
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

func Initialize(router *fiber.App, handlerMap handlers.HandlerMap, dbClient *ent.Client) {
	router.Use(middleware.Security)
	router.Get("/api/preinit", handlerMap.InitializeHandler.PreInit)
	router.Post("/api/initialize", handlerMap.InitializeHandler.Initialize)

	auth := router.Group("/api/auth")
	{
		auth.Post("/login/password", handlerMap.AuthHandler.Login)
	}

	api := router.Group("/api")
	{
		apiV1 := api.Group("/v1")
		{
			apiV1.Get("/visit", handlerMap.VisitHandler.HandleVisitor)
			apiV1.Get("/routes", handlerMap.RouteHandler.GetRoutes)

			apiV1.Get("/api-interface/page", handlerMap.ApiInterfaceHandler.ListApiRoutesPage)

			apiV1.Get("/settings", handlerMap.SettingHandler.GetSettings)
			apiV1.All("/twikoo", handlerMap.CommentHandler.HandleTwikoo).Name("twikoo")
			apiV1.Get("/comment/recent", handlerMap.CommentHandler.RecentComment).Name("recentComment")
			apiV1.Get("/flink/random", handlerMap.FlinkHandler.RandomFlink).Name("randomFlink")
			apiV1.Get("/posts/:id/summary/stream", handlerMap.PostHandler.GetSummaryForStream).Name("postSummaryStream")

			apiV1.Get("/post/search", handlerMap.PostHandler.SearchPosts).Name("postSearch")
			apiV1.Get("/user/search", handlerMap.UserHandler.SearchUsers).Name("userSearch")
			apiV1.Get("/product/search", handlerMap.ProductHandler.SearchProducts).Name("productSearch")

			apiV1.Post("/flink-application/create", handlerMap.FlinkApplicationHandler.CreateFlinkApplication).Name("flinkApplicationCreate")

			apiV1.Use(middleware.FlexibleAuth(dbClient))

			// 首页统计信息接口
			apiV1.Get("/common/statistic", handlerMap.CommonHandler.GetHomeStatistics).Name("homeStatistic")

			apiV1.Get("/user/personal-access-token/list", handlerMap.UserHandler.GetPersonalAccessTokenList).Name("patList")
			apiV1.Get("/user/personal-access-token/query/:id", handlerMap.UserHandler.GetPersonalAccessToken).Name("patInfo")
			apiV1.Post("/user/personal-access-token/create", handlerMap.UserHandler.CreatePat).Name("patCreate")

			initContentRouter(apiV1, handlerMap)
			initInfraRouter(apiV1, handlerMap)
			initMallRouter(apiV1, handlerMap)
			initSystemRouter(apiV1, handlerMap)
		}
	}
}
