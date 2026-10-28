package mail

import (
	"encoding/base64"
	"mime"
	"strings"
)

// mimeEncode 对邮件头字段进行 MIME 编码，保证中文主题/发件人名称可正常显示
func mimeEncode(s string) string {
	if s == "" || isASCII(s) {
		return s
	}
	return mime.QEncoding.Encode("UTF-8", s)
}

// base64Encode 将正文按 UTF-8 进行 Base64 编码（配合 Content-Transfer-Encoding: base64）
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return strings.TrimSpace(s) != ""
}
