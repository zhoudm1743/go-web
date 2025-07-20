package log

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLoggerMiddleware 创建适配我们日志系统的Gin中间件
func GinLoggerMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		// 执行时间
		latency := end.Sub(start)

		// 获取客户端IP
		clientIP := c.ClientIP()
		// 获取HTTP方法
		method := c.Request.Method
		// 获取请求路径
		path := c.Request.URL.Path
		// 获取状态码
		statusCode := c.Writer.Status()
		// 获取错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			logger.WithFields(logrus.Fields{
				"status_code": statusCode,
				"latency":     latency,
				"client_ip":   clientIP,
				"method":      method,
				"path":        path,
				"error":       errorMessage,
			}).Error("[GIN]")
		} else if statusCode >= 400 {
			logger.WithFields(logrus.Fields{
				"status_code": statusCode,
				"latency":     latency,
				"client_ip":   clientIP,
				"method":      method,
				"path":        path,
				"error":       errorMessage,
			}).Warn("[GIN]")
		} else {
			if path != "/favicon.ico" { // 忽略favicon请求的日志
				logger.WithFields(logrus.Fields{
					"status_code": statusCode,
					"latency":     latency,
					"client_ip":   clientIP,
					"method":      method,
					"path":        path,
				}).Debug("[GIN]") // 将正常请求的日志级别降低到DEBUG
			}
		}
	}
}

// GinWriter 实现Gin的Writer接口，转发日志到我们的Logger
type GinWriter struct {
	logger Logger
}

// NewGinWriter 创建GinWriter
func NewGinWriter(logger Logger) *GinWriter {
	return &GinWriter{logger: logger}
}

// Write 实现io.Writer接口
func (w *GinWriter) Write(p []byte) (n int, err error) {
	w.logger.Debug(string(p)) // 将Gin的内部日志转到Debug级别
	return len(p), nil
}

// DisableGinDefaultLogger 禁用Gin的默认日志输出
func DisableGinDefaultLogger() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = &NullWriter{}
}

// NullWriter 一个什么都不做的Writer
type NullWriter struct{}

// Write 实现io.Writer接口但不做任何事
func (n *NullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// GinDebugPrintRouteRegistration Gin路由注册日志
func GinDebugPrintRouteRegistration(logger Logger) {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		// 		logger.Debugf("Route registered: %s %s (%d handlers)", httpMethod, absolutePath, nuHandlers)
	}
}

// CustomGinEngine 创建使用自定义日志的Gin引擎
func CustomGinEngine(logger Logger, disableConsoleLog bool) *gin.Engine {
	// 如果需要禁用控制台日志
	if disableConsoleLog {
		DisableGinDefaultLogger()
	}

	// 禁用路由注册日志
	GinDebugPrintRouteRegistration(logger)

	// 创建一个不带默认中间件的引擎
	r := gin.New()

	// 添加Recovery中间件防止程序崩溃
	r.Use(gin.Recovery())

	// 添加我们的日志中间件
	r.Use(GinLoggerMiddleware(logger))

	return r
}
