package pkg

import (
	"embed"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/infra/schedule/manager"
	album_service "github.com/shuTwT/hoshikuzu/internal/services/content/album"
	albumphoto_service "github.com/shuTwT/hoshikuzu/internal/services/content/albumphoto"
	category_service "github.com/shuTwT/hoshikuzu/internal/services/content/category"
	comment_service "github.com/shuTwT/hoshikuzu/internal/services/content/comment"
	doclibrary_service "github.com/shuTwT/hoshikuzu/internal/services/content/doclibrary"
	doclibrarydetail_service "github.com/shuTwT/hoshikuzu/internal/services/content/doclibrarydetail"
	essay_service "github.com/shuTwT/hoshikuzu/internal/services/content/essay"
	flink_service "github.com/shuTwT/hoshikuzu/internal/services/content/flink"
	flinkapplication_service "github.com/shuTwT/hoshikuzu/internal/services/content/flinkapplication"
	friend_circle_service "github.com/shuTwT/hoshikuzu/internal/services/content/friendcircle"
	knowledgebase_service "github.com/shuTwT/hoshikuzu/internal/services/content/knowledgebase"
	post_service "github.com/shuTwT/hoshikuzu/internal/services/content/post"
	tag_service "github.com/shuTwT/hoshikuzu/internal/services/content/tag"
	file_service "github.com/shuTwT/hoshikuzu/internal/services/infra/file"
	license_service "github.com/shuTwT/hoshikuzu/internal/services/infra/license"
	migration_service "github.com/shuTwT/hoshikuzu/internal/services/infra/migration"
	permission_service "github.com/shuTwT/hoshikuzu/internal/services/infra/permission"
	plugin_service "github.com/shuTwT/hoshikuzu/internal/services/infra/plugin"
	schedulejob_service "github.com/shuTwT/hoshikuzu/internal/services/infra/schedulejob"
	storagestrategy_service "github.com/shuTwT/hoshikuzu/internal/services/infra/storagestrategy"
	theme_service "github.com/shuTwT/hoshikuzu/internal/services/infra/theme"
	visit_service "github.com/shuTwT/hoshikuzu/internal/services/infra/visit"
	coupon_service "github.com/shuTwT/hoshikuzu/internal/services/mall/coupon"
	couponusage_service "github.com/shuTwT/hoshikuzu/internal/services/mall/couponusage"
	member_service "github.com/shuTwT/hoshikuzu/internal/services/mall/member"
	memberlevel_service "github.com/shuTwT/hoshikuzu/internal/services/mall/memberlevel"
	payorder_service "github.com/shuTwT/hoshikuzu/internal/services/mall/payorder"
	product_service "github.com/shuTwT/hoshikuzu/internal/services/mall/product"
	wallet_service "github.com/shuTwT/hoshikuzu/internal/services/mall/wallet"
	api_interface_service "github.com/shuTwT/hoshikuzu/internal/services/system/apiinterface"
	auth_service "github.com/shuTwT/hoshikuzu/internal/services/system/auth"
	common_service "github.com/shuTwT/hoshikuzu/internal/services/system/common"
	notification_service "github.com/shuTwT/hoshikuzu/internal/services/system/notification"
	role_service "github.com/shuTwT/hoshikuzu/internal/services/system/role"
	setting_service "github.com/shuTwT/hoshikuzu/internal/services/system/setting"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"
)

func ExtractDefaultTheme(assetsRes embed.FS) {
	themeDir := "./data/themes"
	defaultThemePath := "./data/themes/hoshikuzu-theme-ace"
	sourceThemePath := "assets/themes/hoshikuzu-theme-ace"

	if _, err := os.Stat(defaultThemePath); err == nil {
		slog.Info("Default theme already exists", "path", defaultThemePath)
		return
	}

	slog.Info("Extracting default theme", "from", sourceThemePath, "to", defaultThemePath)

	if err := os.MkdirAll(themeDir, 0755); err != nil {
		slog.Error("Failed to create themes directory", "error", err.Error())
		return
	}

	sourceDir, err := fs.Sub(assetsRes, sourceThemePath)
	if err != nil {
		slog.Error("Failed to get source directory", "error", err.Error())
		return
	}

	err = fs.WalkDir(sourceDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(defaultThemePath, path)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		sourceFile, err := sourceDir.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = destFile.ReadFrom(sourceFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		slog.Error("Failed to extract default theme", "error", err.Error())
	} else {
		slog.Info("Default theme extracted successfully", "path", defaultThemePath)
	}
}

type ServiceMap struct {
	AlbumService            album_service.AlbumService
	AlbumPhotoService       albumphoto_service.AlbumPhotoService
	ApiInterfaceService     api_interface_service.ApiInterfaceService
	AuthService             auth_service.AuthService
	CategoryService         category_service.CategoryService
	CommentService          comment_service.CommentService
	CommonService           common_service.CommonService
	CouponService           coupon_service.CouponService
	CouponUsageService      couponusage_service.CouponUsageService
	DocLibraryService       doclibrary_service.DocLibraryService
	DocLibraryDetailService doclibrarydetail_service.DocLibraryDetailService
	EssayService            essay_service.EssayService
	FileService             file_service.FileService
	FlinkService            flink_service.FlinkService
	FlinkApplicationService flinkapplication_service.FlinkApplicationService
	LicenseService          license_service.LicenseService
	FriendCircleService     friend_circle_service.FriendCircleService
	KnowledgeBaseService    knowledgebase_service.KnowledgeBaseService
	MemberLevelService      memberlevel_service.MemberLevelService
	MemberService           member_service.MemberService
	MigrationService        migration_service.MigrationService
	NotificationService     notification_service.NotificationService
	PayOrderService         payorder_service.PayOrderService
	PermissionService       permission_service.PermissionService
	PluginService           plugin_service.PluginService
	PostService             post_service.PostService
	ProductService          product_service.ProductService
	RoleService             role_service.RoleService
	ScheduleJobService      schedulejob_service.ScheduleJobService
	SettingService          setting_service.SettingService
	StorageStrategyService  storagestrategy_service.StorageStrategyService
	TagService              tag_service.TagService
	ThemeService            theme_service.ThemeService
	UserService             user_service.UserService
	VisitService            visit_service.VisitService
	WalletService           wallet_service.WalletService
}

func InitializeServices(assetsRes embed.FS, db *ent.Client, scheduleManager *manager.ScheduleManager) ServiceMap {

	albumService := album_service.NewAlbumServiceImpl(db)
	albumPhotoService := albumphoto_service.NewAlbumPhotoServiceImpl(db)
	apiInterfaceService := api_interface_service.NewApiInterfaceServiceImpl(db)
	authService := auth_service.NewAuthServiceImpl(db)
	categoryService := category_service.NewCategoryServiceImpl(db)
	commentService := comment_service.NewCommentServiceImpl(db)
	commonService := common_service.NewCommonServiceImpl(db, user_service.NewUserServiceImpl(db), post_service.NewPostServiceImpl(db), comment_service.NewCommentServiceImpl(db))
	couponService := coupon_service.NewCouponServiceImpl(db)
	couponUsageService := couponusage_service.NewCouponUsageServiceImpl(db)
	doclibraryService := doclibrary_service.NewDocLibraryServiceImpl(db)
	doclibrarydetailService := doclibrarydetail_service.NewDocLibraryDetailServiceImpl(db)
	essayService := essay_service.NewEssayServiceImpl(db)
	fileService := file_service.NewFileServiceImpl(db)
	licenseService := license_service.NewLicenseServiceImpl(db)
	flinkService := flink_service.NewFlinkServiceImpl(db)
	flinkApplicationService := flinkapplication_service.NewFlinkApplicationServiceImpl(db)
	friendCircleService := friend_circle_service.NewFriendCircleServiceImpl(db)
	knowledgeBaseService := knowledgebase_service.NewKnowledgeBaseServiceImpl(db)
	memberLevelService := memberlevel_service.NewMemberLevelServiceImpl(db)
	memberService := member_service.NewMemberServiceImpl(db)
	payOderService := payorder_service.NewPayOrderServiceImpl(db)
	permissionService := permission_service.NewPermissionServiceImpl(db)
	pluginService := plugin_service.NewPluginServiceImpl(db)
	postService := post_service.NewPostServiceImpl(db)
	productService := product_service.NewProductServiceImpl(db)
	roleService := role_service.NewRoleServiceImpl(db)
	settingService := setting_service.NewSettingServiceImpl(db)
	storageStrategyService := storagestrategy_service.NewStorageStrategyServiceImpl(db)
	themeService := theme_service.NewThemeServiceImpl(db)
	tagService := tag_service.NewTagServiceImpl(db)
	userService := user_service.NewUserServiceImpl(db)
	visitService := visit_service.NewVisitServiceImpl(db)
	walletService := wallet_service.NewWalletServiceImpl(db)
	migrationService := migration_service.NewMigrationServiceImpl(db)
	notificationService := notification_service.NewNotificationServiceImpl(db)
	scheduleJobService := schedulejob_service.NewScheduleJobServiceImpl(db, scheduleManager)

	permissionService.LoadPermissionsFromDef(assetsRes)

	serviceMap := ServiceMap{
		AlbumService:            albumService,
		AlbumPhotoService:       albumPhotoService,
		ApiInterfaceService:     apiInterfaceService,
		AuthService:             authService,
		CategoryService:         categoryService,
		CommentService:          commentService,
		CommonService:           commonService,
		CouponService:           couponService,
		CouponUsageService:      couponUsageService,
		DocLibraryService:       doclibraryService,
		DocLibraryDetailService: doclibrarydetailService,
		EssayService:            essayService,
		FileService:             fileService,
		FlinkService:            flinkService,
		FlinkApplicationService: flinkApplicationService,
		FriendCircleService:     friendCircleService,
		KnowledgeBaseService:    knowledgeBaseService,
		LicenseService:          licenseService,
		MemberLevelService:      memberLevelService,
		MemberService:           memberService,
		MigrationService:        migrationService,
		NotificationService:     notificationService,
		PayOrderService:         payOderService,
		PermissionService:       permissionService,
		PluginService:           pluginService,
		PostService:             postService,
		ProductService:          productService,
		RoleService:             roleService,
		ScheduleJobService:      scheduleJobService,
		SettingService:          settingService,
		StorageStrategyService:  storageStrategyService,
		TagService:              tagService,
		ThemeService:            themeService,
		UserService:             userService,
		VisitService:            visitService,
		WalletService:           walletService,
	}

	return serviceMap
}
