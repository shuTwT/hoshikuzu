package handlers

import (
	"github.com/shuTwT/hoshikuzu/ent"
	album_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/album"
	albumphoto_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/albumphoto"
	category_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/category"
	comment_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/comment"
	doclibrary_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/doclibrary"
	doclibrarydetail_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/doclibrarydetail"
	essay_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/essay"
	flink_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/flink"
	flinkapplication_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/flinkapplication"
	flinkgroup_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/flinkgroup"
	friendcircle_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/friendcircle"
	knowledgebase_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/knowledgebase"
	post_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/post"
	tag_handler "github.com/shuTwT/hoshikuzu/internal/handlers/content/tag"
	file_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/file"
	license_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/license"
	migration_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/migration"
	plugin_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/plugin"
	schedulejob_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/schedulejob"
	storagestrategy "github.com/shuTwT/hoshikuzu/internal/handlers/infra/storagestrategy"
	theme_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/theme"
	visit_handler "github.com/shuTwT/hoshikuzu/internal/handlers/infra/visit"
	coupon_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/coupon"
	couponusage_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/couponusage"
	member_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/member"
	memberlevel_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/memberlevel"
	payorder_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/payorder"
	product_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/product"
	wallet_handler "github.com/shuTwT/hoshikuzu/internal/handlers/mall/wallet"
	public_handler "github.com/shuTwT/hoshikuzu/internal/handlers/public"
	apiinterface_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/apiinterface"
	auth_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/auth"
	common_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/common"
	initialize_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/initialize"
	notification_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/notification"
	role_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/role"
	route_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/route"
	setting_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/setting"
	user_handler "github.com/shuTwT/hoshikuzu/internal/handlers/system/user"
	"github.com/shuTwT/hoshikuzu/pkg"
)

type HandlerMap struct {
	AlbumHandler            album_handler.AlbumHandler
	AlbumPhotoHandler       albumphoto_handler.AlbumPhotoHandler
	ApiInterfaceHandler     apiinterface_handler.ApiInterfaceHandler
	AuthHandler             auth_handler.AuthHandler
	CategoryHandler         category_handler.CategoryHandler
	CommentHandler          comment_handler.CommentHandler
	CommonHandler           common_handler.CommonHandler
	CouponHandler           coupon_handler.CouponHandler
	CouponUsageHandler      couponusage_handler.CouponUsageHandler
	DocLibraryHandler       doclibrary_handler.DocLibraryHandler
	DocLibraryDetailHandler doclibrarydetail_handler.DocLibraryDetailHandler
	FileHandler             file_handler.FileHandler
	LicenseHandler          license_handler.LicenseHandler
	FlinkHandler            flink_handler.FlinkHandler
	FlinkApplicationHandler flinkapplication_handler.FlinkApplicationHandler
	FlinkGroupHandler       flinkgroup_handler.FlinkGroupHandler
	FriendCircleHandler     friendcircle_handler.FriendCircleHandler
	InitializeHandler       initialize_handler.InitializeHandler
	KnowledgeBaseHandler    knowledgebase_handler.KnowledgeBaseHandler
	MemberHandler           member_handler.MemberHandler
	MemberLevelHandler      memberlevel_handler.MemberLevelHandler
	MigrationHandler        migration_handler.MigrationHandler
	NotificationHandler     notification_handler.NotificationHandler
	PayOrderHandler         payorder_handler.PayOrderHandler
	PluginHandler           plugin_handler.PluginHandler
	PostHandler             post_handler.PostHandler
	ProductHandler          product_handler.ProductHandler
	RoleHandler             role_handler.RoleHandler
	RouteHandler            route_handler.RouteHandler
	ScheduleJobHandler      schedulejob_handler.ScheduleJobHandler
	SettingHandler          setting_handler.SettingHandler
	TagHandler              tag_handler.TagHandler
	ThemeHandler            theme_handler.ThemeHandler
	UserHandler             user_handler.UserHandler
	EssayHandler            essay_handler.EssayHandler
	StorageStrategyHandler  storagestrategy.StorageStrategyHandler
	VisitHandler            visit_handler.VisitHandler
	WalletHandler           wallet_handler.WalletHandler
	PublicHandler           public_handler.PublicHandler
}

func InitHandler(serviceMap pkg.ServiceMap, db *ent.Client) HandlerMap {
	albumHandler := album_handler.NewAlbumHandlerImpl(serviceMap.AlbumService)
	albnumPhotoHandler := albumphoto_handler.NewAlbumPhotoHandlerImpl(serviceMap.AlbumPhotoService)
	apiInterfaceHandler := apiinterface_handler.NewApiInterfaceHandlerImpl(serviceMap.ApiInterfaceService)
	authHandler := auth_handler.NewAuthHandlerImpl(serviceMap.AuthService)
	categoryHandler := category_handler.NewCategoryHandlerImpl(serviceMap.CategoryService, serviceMap.PostService)
	commentHandler := comment_handler.NewCommentHandlerImpl(serviceMap.CommentService)
	commonHandler := common_handler.NewCommonHandlerImpl(serviceMap.CommonService)
	couponHandler := coupon_handler.NewCouponHandlerImpl(serviceMap.CouponService)
	couponUsageHandler := couponusage_handler.NewCouponUsageHandlerImpl(serviceMap.CouponUsageService)
	doclibraryHandler := doclibrary_handler.NewDocLibraryHandlerImpl(serviceMap.DocLibraryService)
	doclibrarydetailHandler := doclibrarydetail_handler.NewDocLibraryDetailHandlerImpl(serviceMap.DocLibraryDetailService)
	fileHandler := file_handler.NewFileHandlerImpl(serviceMap.FileService, serviceMap.StorageStrategyService)
	licenseHandler := license_handler.NewLicenseHandlerImpl(serviceMap.LicenseService)
	flinkHandler := flink_handler.NewFlinkHandlerImpl(db, serviceMap.FlinkService)
	flinkApplicationHandler := flinkapplication_handler.NewFlinkApplicationHandlerImpl(db, serviceMap.FlinkApplicationService)
	flinkGroupHandler := flinkgroup_handler.NewFlinkGroupHandlerImpl(db, serviceMap.FlinkService)
	friendCircleHandler := friendcircle_handler.NewFriendCircleHandlerImpl(serviceMap.FriendCircleService)
	initializeHandler := initialize_handler.NewInitializeHandlerImpl(db, serviceMap.UserService, serviceMap.SettingService)
	knowledgeBaseHandler := knowledgebase_handler.NewKnowledgeBaseHandlerImpl(serviceMap.KnowledgeBaseService)
	payOrderHandler := payorder_handler.NewPayOrderHandlerImpl(db, serviceMap.PayOrderService)
	postHandler := post_handler.NewPostHandlerImpl(serviceMap.PostService)
	productHandler := product_handler.NewProductHandlerImpl(serviceMap.ProductService)
	roleHandler := role_handler.NewRoleHandlerImpl(serviceMap.RoleService)
	routeHandler := route_handler.NewRouteHandlerImpl()
	settingHandler := setting_handler.NewSettingHandlerImpl(serviceMap.SettingService)
	tagHandler := tag_handler.NewTagHandlerImpl(serviceMap.TagService)
	userHandler := user_handler.NewUserHandlerImpl(serviceMap.UserService, serviceMap.RoleService)
	essayHandler := essay_handler.NewEssayHandlerImpl(serviceMap.EssayService)
	storageStrategyHandler := storagestrategy.NewStorageStrategyHandlerImpl(serviceMap.StorageStrategyService)
	visitHandler := visit_handler.NewVisitHandlerImpl(serviceMap.VisitService)
	walletHandler := wallet_handler.NewWalletHandlerImpl(serviceMap.WalletService)
	memberHandler := member_handler.NewMemberHandlerImpl(serviceMap.UserService, serviceMap.MemberService)
	memberLevelHandler := memberlevel_handler.NewMemberLevelHandlerImpl(serviceMap.MemberLevelService)
	migrationHandler := migration_handler.NewMigrationHandlerImpl(serviceMap.MigrationService)
	notificationHandler := notification_handler.NewNotificationHandlerImpl(serviceMap.NotificationService)
	pluginHandler := plugin_handler.NewPluginHandlerImpl(serviceMap.PluginService)
	scheduleJobHandler := schedulejob_handler.NewScheduleJobHandlerImpl(serviceMap.ScheduleJobService)
	themeHandler := theme_handler.NewThemeHandlerImpl(serviceMap.ThemeService)
	publicHandler := public_handler.NewPublicHandlerImpl()

	handlerMap := HandlerMap{
		AlbumHandler:            albumHandler,
		AlbumPhotoHandler:       albnumPhotoHandler,
		ApiInterfaceHandler:     apiInterfaceHandler,
		AuthHandler:             authHandler,
		CategoryHandler:         categoryHandler,
		CommentHandler:          commentHandler,
		CommonHandler:           commonHandler,
		CouponHandler:           couponHandler,
		CouponUsageHandler:      couponUsageHandler,
		DocLibraryHandler:       doclibraryHandler,
		DocLibraryDetailHandler: doclibrarydetailHandler,
		FileHandler:             fileHandler,
		FlinkHandler:            flinkHandler,
		FlinkApplicationHandler: flinkApplicationHandler,
		FlinkGroupHandler:       flinkGroupHandler,
		FriendCircleHandler:     friendCircleHandler,
		InitializeHandler:       initializeHandler,
		KnowledgeBaseHandler:    knowledgeBaseHandler,
		LicenseHandler:          licenseHandler,
		MemberHandler:           memberHandler,
		MemberLevelHandler:      memberLevelHandler,
		MigrationHandler:        migrationHandler,
		NotificationHandler:     notificationHandler,
		PayOrderHandler:         payOrderHandler,
		PluginHandler:           pluginHandler,
		PostHandler:             postHandler,
		ProductHandler:          productHandler,
		RoleHandler:             roleHandler,
		RouteHandler:            routeHandler,
		ScheduleJobHandler:      scheduleJobHandler,
		SettingHandler:          settingHandler,
		TagHandler:              tagHandler,
		UserHandler:             userHandler,
		EssayHandler:            essayHandler,
		StorageStrategyHandler:  storageStrategyHandler,
		VisitHandler:            visitHandler,
		WalletHandler:           walletHandler,
		ThemeHandler:            themeHandler,
		PublicHandler:           publicHandler,
	}

	return handlerMap

}
