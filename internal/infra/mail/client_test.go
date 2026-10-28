package mail

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

// generateCert 生成一个自签名证书，用于 mock 服务器的 TLS 握手
func generateCert(t *testing.T) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatal(err)
	}
	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return cert, priv
}

// handleSMTPConn 处理一个 SMTP 客户端连接。
// implicit 为 true 时先做 TLS 握手（对应 ssl 模式）；
// startTLS 为 true 时在收到 STARTTLS 后升级为 TLS（对应 tls 模式）；
// failAuth 为 true 时对 AUTH 返回 535，用于验证错误信息透传。
func handleSMTPConn(conn net.Conn, implicit, startTLS, failAuth bool, tlsConf *tls.Config) {
	defer conn.Close()

	if implicit {
		tconn := tls.Server(conn, tlsConf)
		if err := tconn.Handshake(); err != nil {
			return
		}
		conn = tconn
	}

	reader := bufio.NewReader(conn)
	write := func(line string) {
		conn.Write([]byte(line + "\r\n"))
	}

	write("220 localhost ESMTP mock")
	tlsActive := false
	inData := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")

		if inData {
			if line == "." {
				write("250 2.0.0 OK: queued")
				inData = false
			}
			continue
		}

		cmd := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(cmd, "EHLO") || strings.HasPrefix(cmd, "HELO"):
			if tlsActive && startTLS {
				write("250-localhost")
				write("250 AUTH PLAIN")
			} else if startTLS {
				write("250-localhost")
				write("250-AUTH PLAIN LOGIN")
				write("250 STARTTLS")
			} else {
				write("250-localhost")
				write("250 AUTH PLAIN")
			}
		case strings.HasPrefix(cmd, "STARTTLS"):
			write("220 2.0.0 Ready to start TLS")
			tconn := tls.Server(conn, tlsConf)
			if err := tconn.Handshake(); err != nil {
				return
			}
			conn = tconn
			reader = bufio.NewReader(conn)
			tlsActive = true
		case strings.HasPrefix(cmd, "AUTH"):
			if failAuth {
				write("535 5.7.8 Authentication credentials invalid")
			} else {
				write("235 2.7.0 Authentication successful")
			}
		case strings.HasPrefix(cmd, "MAIL FROM"):
			addr := extractAddr(line)
			// MAIL FROM 必须是裸邮箱地址，不允许出现显示名称、空格或嵌套尖括号
			if addr == "" || strings.ContainsAny(addr, "<> ") {
				write("500 5.5.1 Error: bad syntax")
			} else {
				write("250 2.1.0 OK")
			}
		case strings.HasPrefix(cmd, "RCPT TO"):
			write("250 2.1.5 OK")
		case strings.HasPrefix(cmd, "DATA"):
			write("354 End data with <CR><LF>.<CR><LF>")
			inData = true
		case strings.HasPrefix(cmd, "QUIT"):
			write("221 2.0.0 Bye")
			return
		default:
			write("250 2.0.0 OK")
		}
	}
}

// extractAddr 从 "MAIL FROM:<addr>" 中提取尖括号内的地址
func extractAddr(line string) string {
	start := strings.Index(line, "<")
	end := strings.LastIndex(line, ">")
	if start < 0 || end <= start {
		return ""
	}
	return line[start+1 : end]
}

// startMockServer 启动 mock SMTP 服务器，返回主机名、端口与服务器证书
func startMockServer(t *testing.T, implicit, startTLS, failAuth bool) (host string, port int, cert *x509.Certificate) {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { lis.Close() })

	cert, key := generateCert(t)
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{cert.Raw}, PrivateKey: key}},
	}

	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				return
			}
			go handleSMTPConn(conn, implicit, startTLS, failAuth, tlsConf)
		}
	}()

	_, portStr, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	p, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatal(err)
	}
	return "localhost", p, cert
}

// clientTLSConfig 构建信任自签名证书的客户端 TLS 配置
func clientTLSConfig(cert *x509.Certificate) *tls.Config {
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	return &tls.Config{
		ServerName: "localhost",
		RootCAs:    pool,
	}
}

func TestSendMail(t *testing.T) {
	tests := []struct {
		name       string
		encryption string
		implicit   bool
		startTLS   bool
	}{
		{name: "tls_STARTTLS", encryption: "tls", startTLS: true},
		{name: "ssl_implicit", encryption: "ssl", implicit: true},
		{name: "none_plain", encryption: "none", startTLS: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, cert := startMockServer(t, tt.implicit, tt.startTLS, false)
			cfg := Config{
				Host:       host,
				Port:       port,
				Username:   "user",
				Password:   "pass",
				Encryption: tt.encryption,
				Sender:     "no-reply@example.com",
				SenderName: "Hoshikuzu",
				tlsConfig:  clientTLSConfig(cert),
			}
			if err := cfg.SendMail([]string{"a@example.com", "b@example.com"}, "测试主题", "测试正文"); err != nil {
				t.Fatalf("SendMail failed: %v", err)
			}
		})
	}
}

// TestSendMailAuthError 验证认证失败时错误信息会携带具体步骤与 SMTP 错误码
func TestSendMailAuthError(t *testing.T) {
	host, port, cert := startMockServer(t, false, true, true)
	cfg := Config{
		Host:       host,
		Port:       port,
		Username:   "user",
		Password:   "wrong-pass",
		Encryption: "tls",
		Sender:     "no-reply@example.com",
		tlsConfig:  clientTLSConfig(cert),
	}
	err := cfg.SendMail([]string{"a@example.com"}, "主题", "正文")
	if err == nil {
		t.Fatal("expected authentication error")
	}
	for _, want := range []string{"认证", "535"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err.Error(), want)
		}
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for empty config")
	}
	cfg = Config{Host: "smtp.example.com", Port: 587, Sender: "a@example.com"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}
