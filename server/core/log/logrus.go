package log

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/zhoudm1743/go-web/core/conf"
)

// LoggerParams 日志参数
type LoggerParams struct {
	Config *conf.Config
}

// NewLogger 创建日志实例
func NewLogger(p LoggerParams) (Logger, error) {
	log := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(p.Config.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// 设置日志格式
	if p.Config.Log.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		// 创建文本格式化器并启用颜色
		formatter := &logrus.TextFormatter{
			TimestampFormat:           "2006-01-02 15:04:05",
			FullTimestamp:             true,
			ForceColors:               true,
			DisableColors:             false,
			EnvironmentOverrideColors: true,
		}

		// 在Windows环境下，设置ForceColors为true可能不够
		// 如果是在Windows系统下，我们启用额外的设置
		if os.PathSeparator == '\\' { // Windows路径分隔符是反斜杠
			formatter.ForceColors = true
			// 在某些Windows环境中，即使ForceColors为true，某些终端仍可能不显示颜色
		}

		log.SetFormatter(formatter)
	}

	// 设置输出
	var output io.Writer
	if p.Config.Log.OutputPath == "stdout" {
		if runtime.GOOS == "windows" {
			// 在Windows下使用go-colorable
			output = colorable.NewColorableStdout()
		} else {
			output = os.Stdout
		}
	} else if p.Config.Log.OutputPath == "stderr" {
		if runtime.GOOS == "windows" {
			// 在Windows下使用go-colorable
			output = colorable.NewColorableStderr()
		} else {
			output = os.Stderr
		}
	} else {
		// 确保日志目录存在
		dir := filepath.Dir(p.Config.Log.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(p.Config.Log.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = file
	}
	log.SetOutput(output)

	return log, nil
}
