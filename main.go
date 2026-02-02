package main

import (
	"embed"
	"log"

	_ "github.com/mattn/go-sqlite3"

	server "github.com/shuTwT/hoshikuzu/cmd/server"

	_ "github.com/shuTwT/hoshikuzu/docs"
	"github.com/shuTwT/hoshikuzu/internal/infra/logger"
)

//go:embed assets
var assetsRes embed.FS

//go:embed ui/dist
var frontendRes embed.FS

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:13000
// @BasePath /
func main() {
	logger.NewLogger()
	app, cleanup := server.InitializeApp(assetsRes, frontendRes)
	defer cleanup()

	log.Fatal(app.Listen(":13000"))
}
