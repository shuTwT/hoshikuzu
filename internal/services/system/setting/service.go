package setting

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/setting"
)

var SettingKeys = struct {
	SiteName                 string //站点名称
	SiteDescription          string //站点描述
	SiteLogo                 string //站点logo
	SiteFavicon              string //站点favicon
	SiteKeywords             string //站点关键词
	SiteAuthor               string //站点作者
	SiteLanguage             string //站点语言
	SiteTimeZone             string //站点时区
	SiteDateFormat           string //站点日期格式
	SiteTimeFormat           string //站点时间格式
	SiteAllowRegister        string //是否允许注册
	SiteVerifyMail           string //是否验证邮箱
	SiteCommentReview        string //是否评论审核
	SiteUploadMaxSize        string //站点上传最大文件大小（MB）
	SiteEnableCDN            string //是否启用CDN
	SiteCDNURL               string //CDN URL
	SiteICPBeian             string //icp备案号
	SiteGonganBeian          string //公网安备号
	SecurityEnable2FA        string //是否启用2FA
	OpenaiAPIKey             string //openai api key
	OpenaiApiUrl             string //openai api url
	AiModel                  string //openai模型
	AiTemperature            string //openai温度
	AiMaxTokens              string //openai最大token数
	AiTopP                   string //openai top_p
	AiFrequencyPenalty       string //openai 频率惩罚
	AiPresencePenalty        string //openai 存在惩罚
	MailSmtpHost             string //smtp主机
	MailSmtpPort             string //smtp端口
	MailSmtpUsername         string //smtp用户名
	MailSmtpPassword         string //smtp密码
	MailSmtpEncryption       string //smtp加密方式
	MailSmtpSender           string //smtp发件人
	PayEnableEPay            string //启用已支付
	PayEnableAliPay          string //启用支付宝支付
	PayAliPayAppID           string //支付宝应用ID
	PayAliPayAppSecret       string //支付宝应用密钥
	PayEnableWxPay           string //启用微信支付
	PayWxPayAppID            string //微信支付应用ID
	PayWxPayAppSecret        string //微信支付应用密钥
	PayNotifyUrl             string //支付回调URL
	PayOrderExpire           string //订单过期时间（分钟）
	NotifyEnableEmail        string //是否开启邮件通知
	NotifyEnableRegister     string //是否开启注册通知
	NotifyEnableComment      string //是否开启评论通知
	NotifyEnableNewOrder     string //是否开启新订单通知
	NotifyEnablePaySuccess   string //是否开启支付成功通知
	NotifyEnableOrderDeliver string //是否开启订单发货通知
	NotifyEnableOrderSuccess string //是否开启订单成功通知
	NotifyEnableSysErr       string //是否开启系统错误通知
}{
	SiteName:                 "site_name",
	SiteDescription:          "site_description",
	SiteLogo:                 "site_logo",
	SiteFavicon:              "site_favicon",
	SiteKeywords:             "site_keywords",
	SiteAuthor:               "site_author",
	SiteLanguage:             "site_language",
	SiteTimeZone:             "site_time_zone",
	SiteDateFormat:           "site_date_format",
	SiteTimeFormat:           "site_time_format",
	SiteAllowRegister:        "site_allow_register",
	SiteVerifyMail:           "site_verify_mail",
	SiteCommentReview:        "site_comment_review",
	SiteUploadMaxSize:        "site_upload_max_size",
	SiteEnableCDN:            "site_enable_cdn",
	SiteCDNURL:               "site_cdn_url",
	SiteICPBeian:             "site_icp_beian",
	SiteGonganBeian:          "site_gongan_beian",
	SecurityEnable2FA:        "security_enable_2fa",
	OpenaiAPIKey:             "openai_api_key",
	OpenaiApiUrl:             "openai_api_url",
	AiModel:                  "openai_model",
	AiTemperature:            "openai_temperature",
	AiMaxTokens:              "openai_max_tokens",
	AiTopP:                   "openai_top_p",
	AiFrequencyPenalty:       "openai_frequency_penalty",
	AiPresencePenalty:        "openai_presence_penalty",
	MailSmtpHost:             "mail_smtp_host",
	MailSmtpPort:             "mail_smtp_port",
	MailSmtpUsername:         "mail_smtp_username",
	MailSmtpPassword:         "mail_smtp_password",
	MailSmtpEncryption:       "mail_smtp_encryption",
	MailSmtpSender:           "mail_smtp_sender",
	PayEnableEPay:            "pay_enable_epay",
	PayEnableAliPay:          "pay_enable_ali_pay",
	PayAliPayAppID:           "pay_ali_pay_app_id",
	PayAliPayAppSecret:       "pay_ali_pay_app_secret",
	PayEnableWxPay:           "pay_enable_wx_pay",
	PayWxPayAppID:            "pay_wx_pay_app_id",
	PayWxPayAppSecret:        "pay_wx_pay_app_secret",
	PayNotifyUrl:             "pay_notify_url",
	PayOrderExpire:           "pay_order_expire",
	NotifyEnableEmail:        "notify_enable_email",
	NotifyEnableRegister:     "notify_enable_register",
	NotifyEnableComment:      "notify_enable_comment",
	NotifyEnableNewOrder:     "notify_enable_new_order",
	NotifyEnablePaySuccess:   "notify_enable_pay_success",
	NotifyEnableOrderDeliver: "notify_enable_order_deliver",
	NotifyEnableOrderSuccess: "notify_enable_order_success",
	NotifyEnableSysErr:       "notify_enable_sys_err",
}

type SettingService interface {
	GetAllSettings(ctx context.Context) ([]*ent.Setting, error)
	ExistSettingByKey(ctx context.Context, key string) (bool, error)
	GetSettingByKey(ctx context.Context, key string) (*ent.Setting, error)
	UpdateSettingByKey(ctx context.Context, key string, value string) error
	CreateSettingIfNotExist(ctx context.Context, key string, value string) error
	IsSystemInitialized(ctx context.Context) (bool, error)
	SetSystemInitialized(ctx context.Context) error
}

type SettingServiceImpl struct {
	client *ent.Client
}

func NewSettingServiceImpl(client *ent.Client) *SettingServiceImpl {
	return &SettingServiceImpl{client: client}
}

// GetAllSettings 获取所有系统设置
func (s *SettingServiceImpl) GetAllSettings(ctx context.Context) ([]*ent.Setting, error) {
	return s.client.Setting.Query().All(ctx)
}

func (s *SettingServiceImpl) ExistSettingByKey(ctx context.Context, key string) (bool, error) {
	return s.client.Setting.Query().Where(setting.KeyEQ(key)).Exist(ctx)
}

// GetSettingByKey 根据键获取设置
func (s *SettingServiceImpl) GetSettingByKey(ctx context.Context, key string) (*ent.Setting, error) {
	return s.client.Setting.Query().Where(setting.KeyEQ(key)).Only(ctx)
}

func (s *SettingServiceImpl) UpdateSettingByKey(ctx context.Context, key string, value string) error {
	_, err := s.client.Setting.Update().Where(setting.KeyEQ(key)).SetValue(value).Save(ctx)
	return err
}

func (s *SettingServiceImpl) CreateSettingIfNotExist(ctx context.Context, key string, value string) error {

	_, err := s.client.Setting.Create().SetKey(key).SetValue(value).Save(ctx)

	return err
}

// IsSystemInitialized 检查系统是否已初始化
// 通过检查是否存在特定的初始化标记设置来判断
func (s *SettingServiceImpl) IsSystemInitialized(ctx context.Context) (bool, error) {
	exists, err := s.client.Setting.Query().
		Where(setting.KeyEQ("system_initialized")).
		Exist(ctx)

	return exists, err
}

// SetSystemInitialized 标记系统已初始化
func (s *SettingServiceImpl) SetSystemInitialized(ctx context.Context) error {
	// 检查是否存在 system_initialized 设置
	exists, err := s.client.Setting.Query().Where(setting.KeyEQ("system_initialized")).Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		// 如果存在，则更新其值为 true
		_, err = s.client.Setting.Update().Where(setting.KeyEQ("system_initialized")).SetValue("true").Save(ctx)
	} else {
		// 如果不存在，则创建新的设置
		_, err = s.client.Setting.Create().SetKey("system_initialized").SetValue("true").Save(ctx)
	}
	return err
}
