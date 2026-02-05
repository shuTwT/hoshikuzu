package config

import (
	"os"
	"path/filepath"

	"github.com/shuTwT/hoshikuzu/internal/infra/logger"

	"github.com/spf13/viper"
)

var configKeys = []string{
	DATABASE_URL,
	SERVER_PORT,
	SERVER_STAGE,
}

const (
	DATABASE_TYPE          = "database.type"
	DATABASE_URL           = "database.url"
	SERVER_PORT            = "server.port"
	SERVER_STAGE           = "server.stage"
	SERVER_DEBUG           = "server.debug"
	SERVER_TRUSTED_PROXIES = "server.trusted_proxies"
	SWAGGER_ENABLE         = "swagger.enable"
	AUTH_TOKEN_SECRET      = "auth.token_secret"
	AUTH_PAT_SECRET        = "auth.pat_secret"
	PAY_EPAY_ENABLE        = "pay.epay_enable"
	PAY_EPAY_API_URL       = "pay.epay_api_url"
	PAY_EPAY_MERCHANT_ID   = "pay.epay_merchant_id"
	PAY_EPAY_MERCHANT_KEY  = "pay.epay_merchant_key"
	PAY_EPAY_NOTIFY_URL    = "pay.epay_notify_url"
	PAY_EPAY_RETURN_URL    = "pay.epay_return_url"
	// Redis 相关
	Redis_Enable   = "redis.enable"
	REDIS_ADDR     = "redis.addr"
	REDIS_PASSWORD = "redis.password"
	REDIS_DB       = "redis.db"
)

func Init() {
	viper.SetDefault(DATABASE_TYPE, "sqlite")
	viper.SetDefault(DATABASE_URL, "file:./data/sql.db?cache=shared&_fk=1")
	viper.SetDefault(SERVER_PORT, "8000")
	viper.SetDefault(SERVER_STAGE, "dev")
	viper.SetDefault(SERVER_DEBUG, false)
	viper.SetDefault(SERVER_TRUSTED_PROXIES, []string{"127.0.0.1"})
	viper.SetDefault(SWAGGER_ENABLE, true)
	viper.SetDefault(AUTH_TOKEN_SECRET, "your-secret-key")
	viper.SetDefault(AUTH_PAT_SECRET, "your-pat-secret")
	// 易支付相关
	viper.SetDefault(PAY_EPAY_ENABLE, true)
	viper.SetDefault(PAY_EPAY_API_URL, "https://api.pay.com")
	viper.SetDefault(PAY_EPAY_MERCHANT_ID, "your-merchant-id")
	viper.SetDefault(PAY_EPAY_MERCHANT_KEY, "your-merchant-key")
	viper.SetDefault(PAY_EPAY_NOTIFY_URL, "https://api.shhsu.com/api/v1/payorder/notify")
	viper.SetDefault(PAY_EPAY_RETURN_URL, "https://api.shhsu.com/api/v1/payorder/return")
	// Redis 相关
	viper.SetDefault(Redis_Enable, false)
	viper.SetDefault(REDIS_ADDR, "localhost:6379")
	viper.SetDefault(REDIS_PASSWORD, "")
	viper.SetDefault(REDIS_DB, 0)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./data")
	viper.AddConfigPath("$HOME/.hoshikuzu")
	viper.AutomaticEnv()

	for _, key := range configKeys {
		if viper.Get(key) == nil {
			panic("config key not found: " + key)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，在工作空间/data创建默认配置文件
			createDefaultConfig()
		} else {
			// 配置文件存在但读取失败
			panic("fatal error config file: " + err.Error())
		}
	} else {
		logger.Info("已加载配置文件", "config_file", viper.ConfigFileUsed())
	}
}

func createDefaultConfig() {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		logger.Warn("无法获取当前工作目录", "error", err.Error())
		return
	}

	configDir := filepath.Join(workDir, "data")
	configPath := filepath.Join(configDir, "config.toml")

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Panic("无法创建配置目录", "error", err.Error())
		return
	}

	// 在可执行文件同级目录创建配置文件
	if err := viper.WriteConfigAs(configPath); err != nil {
		// 如果写入失败，记录警告但不影响程序运行
		logger.Warn("无法创建默认配置文件", "error", err.Error())
	} else {
		logger.Info("已创建默认配置文件", "config_path", configPath)
	}
}

func GetDatabaseUrl() string {
	return os.Getenv("DATABASE_URL")
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetTrustedProxies() []string {
	return viper.GetStringSlice(SERVER_TRUSTED_PROXIES)
}
