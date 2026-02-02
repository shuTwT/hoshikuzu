package consts

// go-plugin 握手配置常量（宿主和插件必须完全一致）
const (
	MagicCookieKey   = "PLUGIN_STORE_MAGIC"
	MagicCookieValue = "custom-magic-key-123456"
	ProtocolVersion  = 1 // 与 Plugin.ProtocolVersion() 返回值一致
)
