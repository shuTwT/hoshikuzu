package database

import (
	"context"
	"log"

	"github.com/shuTwT/hoshikuzu/ent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	DBType string

	DBUrl string
}

func InitializeDB(cfg DBConfig, autoMigrate bool) (*ent.Client, error) {
	var client *ent.Client
	var err error

	switch cfg.DBType {
	case "sqlite":

		client, err = ent.Open("sqlite3", cfg.DBUrl)
	case "mysql":
		client, err = ent.Open("mysql", cfg.DBUrl)
	case "postgresql":
		client, err = ent.Open("postgres", cfg.DBUrl)
	default:
		log.Fatalf("unsupported database type: %s", cfg.DBType)
	}

	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", cfg.DBType, err)
		return nil, err
	}

	// 主进程执行迁移
	if autoMigrate && !fiber.IsChild() {
		log.Println("Auto migrating schema...")
		if err = client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}
	return client, nil
}
