package conf

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Cache    CacheConfig    `mapstructure:"cache"`
	viper    *viper.Viper   // 存储viper实例，用于获取配置
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string
	Version string
	Mode    string // dev, test, prod
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Host           string
	Port           int
	Engine         string // 引擎类型："gin" 或 "fiber"
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	MaxBodySize    int // 请求体大小限制
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string `mapstructure:"secret"`        // JWT密钥
	AccessExpire  int64  `mapstructure:"accessExpire"`  // 访问令牌过期时间(秒)
	RefreshExpire int64  `mapstructure:"refreshExpire"` // 刷新令牌过期时间(秒)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        string
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, text
	OutputPath string
}

// CacheConfig Cache缓存配置
type CacheConfig struct {
	Type     string // 缓存类型：memory、redis 或 file
	Host     string
	Port     int
	Password string
	DB       int
	Prefix   string // 键前缀
	FilePath string // 文件缓存路径，仅当 Type 为 file 时使用
}

// setDefaultConfig 设置配置的默认值
func setDefaultConfig(config *Config) {
	// 应用配置默认值
	config.App.Name = "go-web"
	config.App.Version = "0.1.0"
	config.App.Mode = "dev"

	// HTTP服务配置默认值
	config.HTTP.Host = "0.0.0.0"
	config.HTTP.Port = 8080
	config.HTTP.Engine = "gin"
	config.HTTP.ReadTimeout = 10 * time.Second
	config.HTTP.WriteTimeout = 10 * time.Second
	config.HTTP.MaxHeaderBytes = 5 << 20 // 5MB
	config.HTTP.MaxBodySize = 10 << 20   // 10MB

	// JWT配置默认值
	config.JWT.Secret = "go-web-secret-key"
	config.JWT.AccessExpire = 7200    // 2小时
	config.JWT.RefreshExpire = 604800 // 7天

	// 数据库配置默认值
	config.Database.Driver = "sqlite"
	config.Database.DSN = "file:go-web.db?cache=shared"
	config.Database.MaxOpenConns = 100
	config.Database.MaxIdleConns = 10
	config.Database.ConnMaxLifetime = time.Hour
	config.Database.LogLevel = "info"

	// 日志配置默认值
	config.Log.Level = "info"
	config.Log.Format = "text"
	config.Log.OutputPath = "logs/app.log"

	// 缓存配置默认值
	config.Cache.Type = "memory"
	config.Cache.Host = "127.0.0.1"
	config.Cache.Port = 6379
	config.Cache.Password = ""
	config.Cache.DB = 0
	config.Cache.Prefix = "go-web:"
	config.Cache.FilePath = "cache"
}

// NewConfig 创建配置
func NewConfig() (*Config, error) {
	// 创建默认配置
	config := &Config{}

	// 设置默认值（无论配置文件是否存在）
	setDefaultConfig(config)

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	// 确保配置目录存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 使用默认配置
		return config, nil
	}

	v := viper.New()
	v.AddConfigPath(configPath)
	configName := os.Getenv("CONFIG_NAME")
	v.SetConfigName(configName)
	if configName == "" {
		v.SetConfigName("config")
	}
	v.SetConfigType("yaml")

	// 读取环境变量
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 使用默认配置
		return config, nil
	}

	// 将文件配置合并到默认配置
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置错误: %w", err)
	}

	// 保存viper实例以便后续使用
	config.viper = v

	return config, nil
}

// Get 通过点分隔的路径获取配置项值
func (c *Config) Get(key string) interface{} {
	// 如果viper实例存在，优先使用它获取
	if c.viper != nil {
		return c.viper.Get(key)
	}

	// viper实例不存在，手动获取
	parts := strings.Split(key, ".")
	if len(parts) == 0 {
		return nil
	}

	// 根据配置路径获取值
	switch parts[0] {
	case "app":
		if len(parts) == 1 {
			return c.App
		}
		switch parts[1] {
		case "name":
			return c.App.Name
		case "version":
			return c.App.Version
		case "mode":
			return c.App.Mode
		}
	case "http":
		if len(parts) == 1 {
			return c.HTTP
		}
		switch parts[1] {
		case "host":
			return c.HTTP.Host
		case "port":
			return c.HTTP.Port
		case "engine":
			return c.HTTP.Engine
		case "readTimeout":
			return c.HTTP.ReadTimeout
		case "writeTimeout":
			return c.HTTP.WriteTimeout
		case "maxHeaderBytes":
			return c.HTTP.MaxHeaderBytes
		case "maxBodySize":
			return c.HTTP.MaxBodySize
		}
	case "jwt":
		if len(parts) == 1 {
			return c.JWT
		}
		switch parts[1] {
		case "secret":
			return c.JWT.Secret
		case "accessExpire":
			return c.JWT.AccessExpire
		case "refreshExpire":
			return c.JWT.RefreshExpire
		}
	case "database":
		if len(parts) == 1 {
			return c.Database
		}
		switch parts[1] {
		case "driver":
			return c.Database.Driver
		case "dsn":
			return c.Database.DSN
		case "maxOpenConns":
			return c.Database.MaxOpenConns
		case "maxIdleConns":
			return c.Database.MaxIdleConns
		case "connMaxLifetime":
			return c.Database.ConnMaxLifetime
		case "logLevel":
			return c.Database.LogLevel
		}
	case "log":
		if len(parts) == 1 {
			return c.Log
		}
		switch parts[1] {
		case "level":
			return c.Log.Level
		case "format":
			return c.Log.Format
		case "outputPath":
			return c.Log.OutputPath
		}
	case "cache":
		if len(parts) == 1 {
			return c.Cache
		}
		switch parts[1] {
		case "type":
			return c.Cache.Type
		case "host":
			return c.Cache.Host
		case "port":
			return c.Cache.Port
		case "password":
			return c.Cache.Password
		case "db":
			return c.Cache.DB
		case "prefix":
			return c.Cache.Prefix
		case "filePath":
			return c.Cache.FilePath
		}
	}

	return nil
}
