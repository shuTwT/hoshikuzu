package config

import (
	"crypto/rand"
	"encoding/base64"
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
	// Redis 相关
	Redis_Enable   = "redis.enable"
	REDIS_ADDR     = "redis.addr"
	REDIS_PASSWORD = "redis.password"
	REDIS_DB       = "redis.db"
	// AI 配置加密密钥（Base64 编码的 32 字节 AES-256 密钥）
	AI_CONFIG_ENCRYPTION_KEY = "ai.config_encryption_key"
)

// 认证密钥的占位默认值。默认配置文件也会写入这两个值，因此仅凭"已配置"无法判断是否安全。
const (
	defaultAuthTokenSecret = "your-secret-key"
	defaultAuthPatSecret   = "your-pat-secret"
)

func Init() {
	viper.SetDefault(DATABASE_TYPE, "sqlite")
	viper.SetDefault(DATABASE_URL, "file:./data/sql.db?cache=shared&_fk=1")
	viper.SetDefault(SERVER_PORT, "8000")
	viper.SetDefault(SERVER_STAGE, "dev")
	viper.SetDefault(SERVER_DEBUG, false)
	viper.SetDefault(SERVER_TRUSTED_PROXIES, []string{"127.0.0.1"})
	viper.SetDefault(SWAGGER_ENABLE, true)
	viper.SetDefault(AUTH_TOKEN_SECRET, defaultAuthTokenSecret)
	viper.SetDefault(AUTH_PAT_SECRET, defaultAuthPatSecret)
	// Redis 相关
	viper.SetDefault(Redis_Enable, false)
	viper.SetDefault(REDIS_ADDR, "localhost:6379")
	viper.SetDefault(REDIS_PASSWORD, "")
	viper.SetDefault(REDIS_DB, 0)
	// AI 配置加密密钥：空则首次启动自动生成并持久化，保证 AI 功能开箱即用
	viper.SetDefault(AI_CONFIG_ENCRYPTION_KEY, "")

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
			ensureAIEncryptionKey()
			createDefaultConfig()
		} else {
			// 配置文件存在但读取失败
			panic("fatal error config file: " + err.Error())
		}
	} else {
		logger.Info("已加载配置文件", "config_file", viper.ConfigFileUsed())
		ensureAIEncryptionKey()
	}

	warnInsecureAuthSecrets()
}

// warnInsecureAuthSecrets 认证密钥为空或仍是占位默认值时打印警告：任何人都可以用公开的默认值伪造登录令牌。
func warnInsecureAuthSecrets() {
	secrets := []struct {
		key          string
		defaultValue string
	}{
		{AUTH_TOKEN_SECRET, defaultAuthTokenSecret},
		{AUTH_PAT_SECRET, defaultAuthPatSecret},
	}
	for _, s := range secrets {
		if v := viper.GetString(s.key); v == "" || v == s.defaultValue {
			logger.Warn("认证密钥未配置或仍为默认值，令牌可被伪造，请在配置文件中设置随机密钥", "config_key", s.key)
		}
	}
}

// ensureAIEncryptionKey 确保 AI 配置加密密钥存在：缺失时生成随机密钥并持久化到配置文件。
// 生成的密钥稳定不变，保证数据库中已加密的 AI API Key 在重启后仍可解密。
func ensureAIEncryptionKey() {
	if viper.GetString(AI_CONFIG_ENCRYPTION_KEY) != "" {
		return
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		logger.Warn("生成 AI 配置加密密钥失败", "error", err.Error())
		return
	}
	key := base64.StdEncoding.EncodeToString(raw)
	viper.Set(AI_CONFIG_ENCRYPTION_KEY, key)
	if viper.ConfigFileUsed() == "" {
		// 配置文件尚未存在，由 createDefaultConfig 负责写盘
		return
	}
	if err := viper.WriteConfig(); err != nil {
		logger.Warn("无法持久化 AI 配置加密密钥", "error", err.Error())
	} else {
		logger.Info("已生成 AI 配置加密密钥并写入配置文件", "config_file", viper.ConfigFileUsed())
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
