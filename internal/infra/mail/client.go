package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"syscall"
	"time"

	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

// EncryptionType SMTP 加密方式
type EncryptionType string

const (
	// EncryptionTLS STARTTLS 加密（默认，常用端口 587）
	EncryptionTLS EncryptionType = "tls"
	// EncryptionSSL 隐式 SSL/TLS 加密（常用端口 465）
	EncryptionSSL EncryptionType = "ssl"
	// EncryptionNone 无加密（常用端口 25）
	EncryptionNone EncryptionType = "none"
)

// Config SMTP 客户端配置
type Config struct {
	Host       string // SMTP 服务器地址
	Port       int    // SMTP 端口
	Username   string // 登录用户名
	Password   string // 登录密码
	Encryption string // 加密方式：tls / ssl / none
	Sender     string // 发件人邮箱地址
	SenderName string // 发件人名称（可选）

	// tlsConfig 可选的 TLS 配置；为空时使用系统根证书并按 Host 校验。
	// 主要用于测试注入自定义证书信任。
	tlsConfig *tls.Config
}

// tlsConf 返回本次发送使用的 TLS 配置
func (c Config) tlsConf() *tls.Config {
	if c.tlsConfig != nil {
		return c.tlsConfig
	}
	return &tls.Config{ServerName: c.Host}
}

// Address 返回 smtp 服务器地址 host:port
func (c Config) Address() string {
	return net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))
}

// From 返回带名称的发件人地址，如 "Name <name@example.com>"
func (c Config) From() string {
	if c.SenderName != "" {
		return fmt.Sprintf("%s <%s>", c.SenderName, c.Sender)
	}
	return c.Sender
}

// Validate 校验配置是否完整
func (c Config) Validate() error {
	if c.Host == "" {
		return model.ErrMailConfigIncomplete
	}
	if c.Port <= 0 || c.Port > 65535 {
		return model.ErrMailConfigIncomplete
	}
	if c.Sender == "" {
		return model.ErrMailConfigIncomplete
	}
	return nil
}

// parseEncryption 解析加密方式，默认按 STARTTLS 处理
func (c Config) parseEncryption() EncryptionType {
	switch strings.ToLower(c.Encryption) {
	case "ssl", "tls", "none":
		return EncryptionType(strings.ToLower(c.Encryption))
	default:
		return EncryptionTLS
	}
}

// wrapErr 将底层错误转换为带步骤上下文与 SMTP 错误码的错误，
// 便于前端与日志定位具体失败环节。
func wrapErr(step string, err error) error {
	if err == nil {
		return nil
	}
	var te *textproto.Error
	if errors.As(err, &te) {
		return fmt.Errorf("%w: %s: SMTP %d %s", model.ErrMailSendFailed, step, te.Code, te.Msg)
	}
	return fmt.Errorf("%w: %s: %s", model.ErrMailSendFailed, step, err.Error())
}

// isConnClosed 判断错误是否为连接被对端关闭/重置
func isConnClosed(err error) bool {
	return errors.Is(err, io.EOF) ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE)
}

// SendMail 发送一封邮件。
// 支持 tls（STARTTLS）、ssl（隐式 TLS）、none（无加密）三种方式，
// 无加密方式下如果服务器支持 STARTTLS 也会自动升级为加密连接。
func (c Config) SendMail(to []string, subject, body string) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if len(to) == 0 {
		return model.ErrMailEmptyRecipient
	}

	enc := c.parseEncryption()
	addr := c.Address()

	var client *smtp.Client
	var err error

	switch enc {
	case EncryptionSSL:
		client, err = c.connectSSL(addr)
	case EncryptionNone:
		client, err = c.connectPlain(addr)
	default:
		client, err = c.connectSTARTTLS(addr)
	}
	if err != nil {
		return err
	}
	defer client.Close()

	// smtp.NewClient / StartTLS 均已完成 EHLO 握手，此处无需再发送 EHLO。
	if ok, _ := client.Extension("AUTH"); ok {
		auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
		if err = client.Auth(auth); err != nil {
			return wrapErr("认证", err)
		}
	}

	// MAIL FROM 只接受裸邮箱地址；发件人名称仅用于邮件头 From 字段
	if err = client.Mail(c.Sender); err != nil {
		return wrapErr("MAIL FROM", err)
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return wrapErr("RCPT TO "+recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return wrapErr("DATA", err)
	}

	msg := buildMessage(c.From(), to, subject, body)
	if _, err = w.Write([]byte(msg)); err != nil {
		return wrapErr("写入正文", err)
	}
	if err = w.Close(); err != nil {
		return wrapErr("结束正文", err)
	}

	if err = client.Quit(); err != nil {
		return wrapErr("QUIT", err)
	}
	return nil
}

// Send 发送一封文本邮件，收件人、主题、正文均不可为空
func (c Config) Send(to []string, subject, body string) error {
	return c.SendMail(to, subject, body)
}

// connectSSL 使用隐式 TLS 直连 465 端口
func (c Config) connectSSL(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, c.tlsConf())
	if err != nil {
		if isConnClosed(err) {
			return nil, fmt.Errorf("%w: 建立 SSL 连接: %s（服务器未响应 SSL 握手，若使用 587 端口可能要求 STARTTLS 加密，请将加密方式改为 tls）", model.ErrMailSendFailed, err.Error())
		}
		return nil, wrapErr("建立 SSL 连接", err)
	}
	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		conn.Close()
		return nil, wrapErr("SMTP 握手", err)
	}
	return client, nil
}

// connectPlain 先以明文连接，若服务器支持 STARTTLS 则自动升级
func (c Config) connectPlain(addr string) (*smtp.Client, error) {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, wrapErr("建立连接", err)
	}
	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		conn.Close()
		return nil, wrapErr("SMTP 握手", err)
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(c.tlsConf()); err != nil {
			client.Close()
			return nil, wrapErr("STARTTLS 升级", err)
		}
	}
	return client, nil
}

// connectSTARTTLS 连接后立即发送 STARTTLS 命令升级为加密连接
func (c Config) connectSTARTTLS(addr string) (*smtp.Client, error) {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, wrapErr("建立连接", err)
	}
	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		conn.Close()
		if isConnClosed(err) {
			return nil, fmt.Errorf("%w: SMTP 握手: %s（服务器在握手阶段关闭连接，若使用 465 端口可能要求隐式 SSL 加密，请将加密方式改为 ssl）", model.ErrMailSendFailed, err.Error())
		}
		return nil, wrapErr("SMTP 握手", err)
	}
	ok, _ := client.Extension("STARTTLS")
	if !ok {
		client.Close()
		return nil, fmt.Errorf("%w: 服务器不支持 STARTTLS 加密", model.ErrMailSendFailed)
	}
	if err := client.StartTLS(c.tlsConf()); err != nil {
		client.Close()
		return nil, wrapErr("STARTTLS 升级", err)
	}
	return client, nil
}

// buildMessage 构造 RFC 5322 邮件内容
func buildMessage(from string, to []string, subject, body string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", strings.Join(to, ", "))
	fmt.Fprintf(&b, "Subject: %s\r\n", mimeEncode(subject))
	fmt.Fprintf(&b, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&b, "Content-Type: text/plain; charset=UTF-8\r\n")
	fmt.Fprintf(&b, "Content-Transfer-Encoding: base64\r\n")
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(&b, "\r\n")
	b.WriteString(base64Encode(body))
	return b.String()
}
