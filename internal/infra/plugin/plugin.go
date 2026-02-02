package plugin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	plugin_shared "github.com/shuTwT/hoshikuzu/pkg/plugin/shared"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GO_PLUGIN",
	MagicCookieValue: "gobee",
}

var pluginMap = map[string]plugin.Plugin{
	"greeter": &plugin_shared.GreeterPlugin{},
}

func RunPlugin() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	pluginPath := filepath.Join(".", "plugins")
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		log.Fatalf("插件不存在:%s", pluginPath)
	}

	// 启动插件客户端
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugins/greeter"),
		Logger:          logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Fatalf("插件连接失败:%s", err)
	}

	raw, err := rpcClient.Dispense("greeter")
	if err != nil {
		log.Fatal(err)
	}
	greeter := raw.(plugin_shared.Greeter)
	fmt.Println(greeter.Init())
}
