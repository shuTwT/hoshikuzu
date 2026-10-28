package pkg

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/shuTwT/hoshikuzu/ent"
	plugin_infra "github.com/shuTwT/hoshikuzu/internal/infra/plugin"
	"github.com/shuTwT/hoshikuzu/internal/infra/schedule/manager"
	ai_service "github.com/shuTwT/hoshikuzu/internal/services/ai/chat"
	album_service "github.com/shuTwT/hoshikuzu/internal/services/content/album"
	albumphoto_service "github.com/shuTwT/hoshikuzu/internal/services/content/albumphoto"
	category_service "github.com/shuTwT/hoshikuzu/internal/services/content/category"
	comment_service "github.com/shuTwT/hoshikuzu/internal/services/content/comment"
	essay_service "github.com/shuTwT/hoshikuzu/internal/services/content/essay"
	flink_service "github.com/shuTwT/hoshikuzu/internal/services/content/flink"
	flinkapplication_service "github.com/shuTwT/hoshikuzu/internal/services/content/flinkapplication"
	friend_circle_service "github.com/shuTwT/hoshikuzu/internal/services/content/friendcircle"
	menu_service "github.com/shuTwT/hoshikuzu/internal/services/content/menu"
	post_service "github.com/shuTwT/hoshikuzu/internal/services/content/post"
	tag_service "github.com/shuTwT/hoshikuzu/internal/services/content/tag"
	file_service "github.com/shuTwT/hoshikuzu/internal/services/infra/file"
	license_service "github.com/shuTwT/hoshikuzu/internal/services/infra/license"
	mail_service "github.com/shuTwT/hoshikuzu/internal/services/infra/mail"
	mailtemplate_service "github.com/shuTwT/hoshikuzu/internal/services/infra/mailtemplate"
	migration_service "github.com/shuTwT/hoshikuzu/internal/services/infra/migration"
	permission_service "github.com/shuTwT/hoshikuzu/internal/services/infra/permission"
	plugin_service "github.com/shuTwT/hoshikuzu/internal/services/infra/plugin"
	schedulejob_service "github.com/shuTwT/hoshikuzu/internal/services/infra/schedulejob"
	smstemplate_service "github.com/shuTwT/hoshikuzu/internal/services/infra/smstemplate"
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
	AIService               ai_service.AIService
	AlbumService            album_service.AlbumService
	AlbumPhotoService       albumphoto_service.AlbumPhotoService
	AuthService             auth_service.AuthService
	CategoryService         category_service.CategoryService
	CommentService          comment_service.CommentService
	CommonService           common_service.CommonService
	CouponService           coupon_service.CouponService
	CouponUsageService      couponusage_service.CouponUsageService
	EssayService            essay_service.EssayService
	FileService             file_service.FileService
	FlinkService            flink_service.FlinkService
	FlinkApplicationService flinkapplication_service.FlinkApplicationService
	LicenseService          license_service.LicenseService
	FriendCircleService     friend_circle_service.FriendCircleService
	MailService             mail_service.MailService
	MailTemplateService     mailtemplate_service.MailTemplateService
	MenuService             menu_service.MenuService
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
	SmsTemplateService      smstemplate_service.SmsTemplateService
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
	authService := auth_service.NewAuthServiceImpl(db)
	categoryService := category_service.NewCategoryServiceImpl(db)
	commentService := comment_service.NewCommentServiceImpl(db)
	essayService := essay_service.NewEssayServiceImpl(db)
	fileService := file_service.NewFileServiceImpl(db)
	licenseService := license_service.NewLicenseServiceImpl(db)
	flinkService := flink_service.NewFlinkServiceImpl(db)
	flinkApplicationService := flinkapplication_service.NewFlinkApplicationServiceImpl(db)
	friendCircleService := friend_circle_service.NewFriendCircleServiceImpl(db)
	menuService := menu_service.NewMenuServiceImpl(db)
	settingService := setting_service.NewSettingServiceImpl(db)
	mailService := mail_service.NewMailServiceImpl(settingService)
	mailTemplateService := mailtemplate_service.NewMailTemplateServiceImpl(db)
	memberLevelService := memberlevel_service.NewMemberLevelServiceImpl(db)
	memberService := member_service.NewMemberServiceImpl(db)
	payOderService := payorder_service.NewPayOrderServiceImpl(db, settingService)
	permissionService := permission_service.NewPermissionServiceImpl(db)
	pluginManager := plugin_infra.NewPluginManager(db)
	pluginService := plugin_service.NewPluginServiceImpl(db, pluginManager)
	productService := product_service.NewProductServiceImpl(db)
	roleService := role_service.NewRoleServiceImpl(db)
	aiService := ai_service.NewAIServiceImpl(db)
	if err := aiService.CleanupLegacySettings(context.Background()); err != nil {
		panic("failed cleaning legacy AI settings: " + err.Error())
	}
	if err := aiService.MigrateLegacyConfig(context.Background()); err != nil {
		panic("failed migrating legacy AI config: " + err.Error())
	}
	postService := post_service.NewPostServiceImpl(db, aiService)
	commonService := common_service.NewCommonServiceImpl(db, user_service.NewUserServiceImpl(db), postService, comment_service.NewCommentServiceImpl(db))
	couponService := coupon_service.NewCouponServiceImpl(db)
	couponUsageService := couponusage_service.NewCouponUsageServiceImpl(db)
	storageStrategyService := storagestrategy_service.NewStorageStrategyServiceImpl(db)
	smsTemplateService := smstemplate_service.NewSmsTemplateServiceImpl(db)
	themeService := theme_service.NewThemeServiceImpl(db, settingService)
	tagService := tag_service.NewTagServiceImpl(db)
	userService := user_service.NewUserServiceImpl(db)
	visitService := visit_service.NewVisitServiceImpl(db)
	walletService := wallet_service.NewWalletServiceImpl(db)
	migrationService := migration_service.NewMigrationServiceImpl(db)
	notificationService := notification_service.NewNotificationServiceImpl(db)
	scheduleJobService := schedulejob_service.NewScheduleJobServiceImpl(db, scheduleManager)

	permissionService.LoadPermissionsFromDef(assetsRes)

	serviceMap := ServiceMap{
		AIService:               aiService,
		AlbumService:            albumService,
		AlbumPhotoService:       albumPhotoService,
		AuthService:             authService,
		CategoryService:         categoryService,
		CommentService:          commentService,
		CommonService:           commonService,
		CouponService:           couponService,
		CouponUsageService:      couponUsageService,
		EssayService:            essayService,
		FileService:             fileService,
		FlinkService:            flinkService,
		FlinkApplicationService: flinkApplicationService,
		FriendCircleService:     friendCircleService,
		LicenseService:          licenseService,
		MailService:             mailService,
		MailTemplateService:     mailTemplateService,
		MenuService:             menuService,
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
		SmsTemplateService:      smsTemplateService,
		StorageStrategyService:  storageStrategyService,
		TagService:              tagService,
		ThemeService:            themeService,
		UserService:             userService,
		VisitService:            visitService,
		WalletService:           walletService,
	}

	return serviceMap
}
