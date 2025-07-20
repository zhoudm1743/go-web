package utils

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/zhoudm1743/go-web/core/facades"
	"gorm.io/gorm"
)

// 注意: 使用本模块前需要先安装以下依赖:
// go get github.com/casbin/casbin/v2
// go get github.com/casbin/gorm-adapter/v3

var (
	casbinEnforcer *casbin.Enforcer
	once           sync.Once
)

// Casbin 获取casbin实例
func Casbin() *casbin.Enforcer {
	once.Do(func() {
		// 获取数据库连接
		db := facades.DB()
		// 创建Casbin Adapter
		adapter, err := gormadapter.NewAdapterByDB(db)
		if err != nil {
			panic("初始化casbin adapter失败: " + err.Error())
		}

		// 从字符串初始化模型
		m := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
		`

		// 创建enforcer
		modelObj, err := model.NewModelFromString(m)
		if err != nil {
			panic("创建模型失败: " + err.Error())
		}

		enforcer, err := casbin.NewEnforcer(modelObj, adapter)
		if err != nil {
			panic("初始化casbin enforcer失败: " + err.Error())
		}

		// 加载策略
		err = enforcer.LoadPolicy()
		if err != nil {
			panic("加载casbin策略失败: " + err.Error())
		}

		casbinEnforcer = enforcer
	})

	return casbinEnforcer
}

// InitCasbinTables 初始化Casbin数据表和基本策略
func InitCasbinTables(db *gorm.DB) error {
	// 创建casbin表
	// 此处不需要实际使用adapter，只是为了确保表结构正确创建
	_, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}

	// 获取enforcer
	enforcer := Casbin()

	// 创建基础角色和权限
	// 1. 超级管理员角色拥有所有权限
	_, err = enforcer.AddPolicy("1", "/*", "*") // roleID=1 为超级管理员
	if err != nil {
		return err
	}

	// 2. 普通用户只能访问部分API
	_, err = enforcer.AddPolicy("2", "/api/user/info", "GET") // roleID=2 为普通用户
	if err != nil {
		return err
	}

	_, err = enforcer.AddPolicy("2", "/api/user/changePassword", "PUT")
	if err != nil {
		return err
	}

	// 保存策略到数据库
	return enforcer.SavePolicy()
}
