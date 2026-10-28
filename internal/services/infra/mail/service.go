package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	mail_infra "github.com/shuTwT/hoshikuzu/internal/infra/mail"
	setting_service "github.com/shuTwT/hoshikuzu/internal/services/system/setting"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

// EmailConfig 邮件配置（对应前端 email 组设置的 JSON 结构）
type EmailConfig struct {
	SmtpHost       string `json:"smtpHost"`
	SmtpPort       int    `json:"smtpPort"`
	SmtpUsername   string `json:"smtpUsername"`
	SmtpPassword   string `json:"smtpPassword"`
	SmtpEncryption string `json:"smtpEncryption"`
	SenderEmail    string `json:"senderEmail"`
	SenderName     string `json:"senderName"`
}

// MailService 邮件服务接口
type MailService interface {
	// SendTestEmail 发送一封测试邮件到指定邮箱，验证 SMTP 配置是否可用
	SendTestEmail(ctx context.Context, to, subject, content string) error
	// GetMailConfig 读取当前保存的邮件配置
	GetMailConfig(ctx context.Context) (*EmailConfig, error)
}

type MailServiceImpl struct {
	settingService setting_service.SettingService
}

func NewMailServiceImpl(settingService setting_service.SettingService) *MailServiceImpl {
	return &MailServiceImpl{settingService: settingService}
}

// GetMailConfig 从通用设置中读取 email 组配置
func (s *MailServiceImpl) GetMailConfig(ctx context.Context) (*EmailConfig, error) {
	setting, err := s.settingService.GetSettingByKey(ctx, "email")
	if err != nil {
		return nil, fmt.Errorf("%w: 未找到邮件配置: %s", model.ErrMailConfigIncomplete, err.Error())
	}

	var cfg EmailConfig
	if err := json.Unmarshal([]byte(setting.Value), &cfg); err != nil {
		return nil, fmt.Errorf("%w: 邮件配置格式错误: %s", model.ErrMailConfigIncomplete, err.Error())
	}
	return &cfg, nil
}

// SendTestEmail 使用当前保存的 SMTP 配置发送测试邮件
func (s *MailServiceImpl) SendTestEmail(ctx context.Context, to, subject, content string) error {
	cfg, err := s.GetMailConfig(ctx)
	if err != nil {
		return err
	}

	if subject == "" {
		subject = "Hoshikuzu 测试邮件"
	}
	if content == "" {
		content = "这是一封来自 Hoshikuzu CMS 的测试邮件，收到即表示 SMTP 配置正确。"
	}

	client := mail_infra.Config{
		Host:       cfg.SmtpHost,
		Port:       cfg.SmtpPort,
		Username:   cfg.SmtpUsername,
		Password:   cfg.SmtpPassword,
		Encryption: cfg.SmtpEncryption,
		Sender:     cfg.SenderEmail,
		SenderName: cfg.SenderName,
	}

	if err := client.SendMail([]string{to}, subject, content); err != nil {
		slog.Error("发送测试邮件失败",
			"to", to,
			"smtpHost", cfg.SmtpHost,
			"smtpPort", cfg.SmtpPort,
			"smtpEncryption", cfg.SmtpEncryption,
			"error", err.Error())
		return err
	}
	return nil
}

// ToMailInfraConfig 将配置转换为基础设施层配置（供其他模块复用发送能力）
func (c *EmailConfig) ToMailInfraConfig() mail_infra.Config {
	port := c.SmtpPort
	if port <= 0 || port > 65535 {
		port = 587
	}
	return mail_infra.Config{
		Host:       c.SmtpHost,
		Port:       port,
		Username:   c.SmtpUsername,
		Password:   c.SmtpPassword,
		Encryption: c.SmtpEncryption,
		Sender:     c.SenderEmail,
		SenderName: c.SenderName,
	}
}
