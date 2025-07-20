package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/zhoudm1743/go-web/core/facades"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID          uint           `gorm:"primarykey" json:"id"`                                         // 主键ID
	CreatedAt   time.Time      `json:"createdAt"`                                                    // 创建时间
	UpdatedAt   time.Time      `json:"updatedAt"`                                                    // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                               // 删除时间
	UUID        uuid.UUID      `gorm:"type:char(36);index;comment:用户UUID" json:"uuid"`               // 用户UUID
	Username    string         `gorm:"type:varchar(50);not null;unique;comment:用户名" json:"username"` // 用户名
	Password    string         `gorm:"type:varchar(100);not null;comment:密码" json:"-"`               // 密码
	Nickname    string         `gorm:"type:varchar(50);comment:昵称" json:"nickname"`                  // 昵称
	RealName    string         `gorm:"type:varchar(50);comment:真实姓名" json:"realName"`                // 真实姓名
	Avatar      string         `gorm:"type:varchar(255);comment:头像" json:"avatar"`                   // 头像
	Email       string         `gorm:"type:varchar(100);comment:邮箱" json:"email"`                    // 邮箱
	Mobile      string         `gorm:"type:varchar(20);comment:手机号" json:"mobile"`                   // 手机号
	Status      uint           `gorm:"type:tinyint(1);default:1;comment:状态 1:启用 2:禁用" json:"status"` // 状态
	RoleID      uint           `gorm:"comment:角色ID" json:"roleId"`                                   // 角色ID
	LastLoginAt time.Time      `gorm:"comment:最后登录时间" json:"lastLoginAt"`                            // 最后登录时间
	LastLoginIP string         `gorm:"type:varchar(50);comment:最后登录IP" json:"lastLoginIp"`           // 最后登录IP
}

// Role 角色模型
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`                                         // 主键ID
	CreatedAt time.Time      `json:"createdAt"`                                                    // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`                                                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                               // 删除时间
	Name      string         `gorm:"type:varchar(50);not null;comment:角色名称" json:"name"`           // 角色名称
	Code      string         `gorm:"type:varchar(50);not null;unique;comment:角色编码" json:"code"`    // 角色编码
	Sort      uint           `gorm:"default:0;comment:排序" json:"sort"`                             // 排序
	Status    uint           `gorm:"type:tinyint(1);default:1;comment:状态 1:启用 2:禁用" json:"status"` // 状态
	Remark    string         `gorm:"type:varchar(255);comment:备注" json:"remark"`                   // 备注
	Menus     []*Menu        `gorm:"many2many:role_menus;" json:"menus"`                           // 角色菜单关联
}

// Menu 菜单模型 - 重构符合前端要求的菜单结构
type Menu struct {
	ID           uint           `gorm:"primarykey" json:"id"`                                                  // 主键ID
	CreatedAt    time.Time      `json:"-"`                                                                     // 创建时间
	UpdatedAt    time.Time      `json:"-"`                                                                     // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                                                        // 删除时间
	PID          *uint          `gorm:"column:parent_id;default:null;comment:父菜单ID" json:"pid"`                // 父菜单ID，顶级菜单为null
	Name         string         `gorm:"type:varchar(50);not null;comment:路由名称" json:"name"`                    // 路由名称(唯一标识)
	Path         string         `gorm:"type:varchar(100);comment:路由路径" json:"path"`                            // 路由路径
	Component    string         `gorm:"type:varchar(100);comment:组件路径" json:"componentPath"`                   // 组件路径
	Redirect     string         `gorm:"type:varchar(100);comment:重定向路径" json:"redirect"`                       // 重定向路径
	Icon         string         `gorm:"type:varchar(50);comment:图标" json:"icon"`                               // 图标
	Title        string         `gorm:"type:varchar(50);comment:标题" json:"title"`                              // 菜单标题
	Order        int            `gorm:"default:0;comment:排序" json:"order"`                                     // 排序值
	Hidden       bool           `gorm:"default:false;comment:是否隐藏" json:"hide"`                                // 是否隐藏
	KeepAlive    bool           `gorm:"default:false;comment:是否缓存" json:"keepAlive"`                           // 是否缓存
	RequiresAuth bool           `gorm:"default:true;comment:是否需要认证" json:"requiresAuth"`                       // 是否需要认证
	WithoutTab   bool           `gorm:"default:false;comment:是否不添加到标签页" json:"withoutTab"`                     // 是否不添加到标签页
	PinTab       bool           `gorm:"default:false;comment:是否固定在标签页" json:"pinTab"`                          // 是否固定在标签页
	MenuType     string         `gorm:"type:varchar(10);default:'page';comment:菜单类型 page|dir" json:"menuType"` // 菜单类型
	Status       uint           `gorm:"type:tinyint(1);default:1;comment:状态 1:启用 2:禁用" json:"-"`               // 状态
	Roles        []*Role        `gorm:"many2many:role_menus;" json:"-"`                                        // 菜单角色关联
}

// AdminRole 管理员角色中间表
type AdminRole struct {
	AdminID uint `gorm:"primarykey;comment:管理员ID" json:"adminId"`
	RoleID  uint `gorm:"primarykey;comment:角色ID" json:"roleId"`
}

// RoleMenu 角色菜单中间表
type RoleMenu struct {
	RoleID uint `gorm:"primarykey;comment:角色ID" json:"roleId"`
	MenuID uint `gorm:"primarykey;comment:菜单ID" json:"menuId"`
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 校验密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateDefaultAdminIfNotExists 创建默认管理员账号
func CreateDefaultAdminIfNotExists(db *gorm.DB) error {
	var count int64
	db.Model(&Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 创建默认管理员角色
	adminRole := &Role{
		Name:   "超级管理员",
		Code:   "super",
		Sort:   1,
		Status: 1,
		Remark: "系统默认创建的超级管理员角色",
	}
	if err := db.Create(adminRole).Error; err != nil {
		return err
	}

	// 创建默认仪表盘目录
	dashboardMenu := &Menu{
		PID:          nil,
		Name:         "dashboard",
		Path:         "/dashboard",
		Component:    "LAYOUT",
		Icon:         "icon-park-outline:analysis",
		Title:        "仪表盘",
		Order:        1,
		RequiresAuth: true,
		MenuType:     "dir",
		Status:       1,
	}
	if err := db.Create(dashboardMenu).Error; err != nil {
		return err
	}

	// 创建仪表盘子菜单 - 只保留工作台
	workbenchMenu := &Menu{
		PID:          &dashboardMenu.ID,
		Name:         "workbench",
		Path:         "/dashboard/workbench",
		Component:    "/dashboard/workbench/index.vue",
		Icon:         "icon-park-outline:alarm",
		Title:        "工作台",
		Order:        1,
		RequiresAuth: true,
		PinTab:       true,
		MenuType:     "page",
		Status:       1,
	}
	if err := db.Create(workbenchMenu).Error; err != nil {
		return err
	}

	// 创建系统设置目录
	settingMenu := &Menu{
		PID:          nil,
		Name:         "setting",
		Path:         "/setting",
		Component:    "",
		Icon:         "icon-park-outline:setting",
		Title:        "系统设置",
		Order:        2,
		RequiresAuth: true,
		MenuType:     "dir",
		Status:       1,
	}
	if err := db.Create(settingMenu).Error; err != nil {
		return err
	}

	// 创建菜单管理
	menuSettingMenu := &Menu{
		PID:          &settingMenu.ID,
		Name:         "menuSetting",
		Path:         "/setting/menu",
		Component:    "/setting/menu/index.vue",
		Icon:         "icon-park-outline:application-menu",
		Title:        "菜单管理",
		Order:        1,
		RequiresAuth: true,
		MenuType:     "page",
		Status:       1,
	}
	if err := db.Create(menuSettingMenu).Error; err != nil {
		return err
	}

	// 创建角色管理
	roleSettingMenu := &Menu{
		PID:          &settingMenu.ID,
		Name:         "roleSetting",
		Path:         "/setting/role",
		Component:    "/setting/role/index.vue",
		Icon:         "icon-park-outline:people-safe",
		Title:        "角色管理",
		Order:        2,
		RequiresAuth: true,
		MenuType:     "page",
		Status:       1,
	}
	if err := db.Create(roleSettingMenu).Error; err != nil {
		return err
	}

	// 创建管理员管理
	adminSettingMenu := &Menu{
		PID:          &settingMenu.ID,
		Name:         "adminSetting",
		Path:         "/setting/admin",
		Component:    "/setting/account/index.vue",
		Icon:         "icon-park-outline:every-user",
		Title:        "管理员管理",
		Order:        3,
		RequiresAuth: true,
		MenuType:     "page",
		Status:       1,
	}
	if err := db.Create(adminSettingMenu).Error; err != nil {
		return err
	}

	// 关联角色和菜单
	roleMenus := []RoleMenu{
		{RoleID: adminRole.ID, MenuID: dashboardMenu.ID},
		{RoleID: adminRole.ID, MenuID: workbenchMenu.ID},
		{RoleID: adminRole.ID, MenuID: settingMenu.ID},
		{RoleID: adminRole.ID, MenuID: menuSettingMenu.ID},
		{RoleID: adminRole.ID, MenuID: roleSettingMenu.ID},
		{RoleID: adminRole.ID, MenuID: adminSettingMenu.ID},
	}

	if err := db.Create(&roleMenus).Error; err != nil {
		return err
	}

	// 创建默认管理员账号
	hashedPassword, err := HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := &Admin{
		UUID:     uuid.New(),
		Username: "admin",
		Password: hashedPassword,
		Nickname: "管理员",
		RealName: "系统管理员",
		Avatar:   "https://avatar.vercel.sh/admin.svg?text=A",
		Email:    "admin@example.com",
		Status:   1,
		RoleID:   adminRole.ID,
	}

	return db.Create(admin).Error
}

// GetRoles 获取管理员角色列表
func (u *Admin) GetRoles() []string {
	// 查询用户角色
	var role Role
	if err := facades.DB().First(&role, u.RoleID).Error; err != nil {
		return []string{}
	}
	return []string{role.Code}
}
