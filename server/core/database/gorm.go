package database

import (
	"fmt"
	"os"
	"time"

	"github.com/glebarez/sqlite" // 纯Go的SQLite实现，不需要CGO
	"github.com/zhoudm1743/go-web/core/conf"
	"github.com/zhoudm1743/go-web/core/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormDB 包装gorm.DB，添加额外方法
type GormDB struct {
	DB *gorm.DB
}

// DBParams 数据库参数
type DBParams struct {
	Config *conf.Config
	Logger log.Logger
}

// NewDB 创建数据库连接
func NewDB(p DBParams) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// 根据驱动类型创建对应的方言
	switch p.Config.Database.Driver {
	case "mysql":
		dialector = mysql.Open(p.Config.Database.DSN)
	case "postgres":
		dialector = postgres.Open(p.Config.Database.DSN)
	case "sqlite":
		// 判断文件是否存在
		if _, err := os.Stat(p.Config.Database.DSN); os.IsNotExist(err) {
			// 创建文件
			os.Create(p.Config.Database.DSN)
		}
		dialector = sqlite.Open(p.Config.Database.DSN)
	case "memory":
		// 使用内存SQLite，不需要CGO
		dialector = sqlite.Open(":memory:")
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", p.Config.Database.Driver)
	}

	// 创建日志记录器
	logLevel := logger.Error
	switch p.Config.Database.LogLevel {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	}

	// 自定义GORM日志适配器
	gormLogger := logger.New(
		&logWriter{p.Logger},
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,  // 忽略记录未找到错误
			Colorful:                  false, // 禁用彩色打印
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// 获取底层的SQL DB以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(p.Config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(p.Config.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(p.Config.Database.ConnMaxLifetime)

	return db, nil
}

// logWriter 日志写入器，将GORM日志适配到我们的Logger接口
type logWriter struct {
	Logger log.Logger
}

func (w *logWriter) Printf(format string, args ...interface{}) {
	w.Logger.Infof(format, args...)
}

// OnStop 数据库关闭钩子
func OnStop(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction 执行事务
func (g *GormDB) Transaction(fn func(tx *gorm.DB) error) error {
	return g.DB.Transaction(fn)
}

// Close 关闭数据库连接
func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
