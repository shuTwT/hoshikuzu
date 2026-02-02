// 插件接口 pkg/plugin/shared/interface.go
package plugin

import (
	shared "github.com/shuTwT/hoshikuzu/pkg/plugin/shared"

	"github.com/hashicorp/go-plugin"
)

type PluginDemo struct {
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GO_PLUGIN",
	MagicCookieValue: "gobee",
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"demo": &shared.GreeterPlugin{},
		},
	})
}
