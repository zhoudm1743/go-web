package facades

import (
	"github.com/zhoudm1743/go-web/core/conf"
)

var configInstance *conf.Config

// SetConfig 设置全局配置实例
func SetConfig(config *conf.Config) {
	configInstance = config
}

// Config 获取全局配置实例
func Config() *conf.Config {
	return configInstance
}
