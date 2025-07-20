package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/zhoudm1743/go-web/core/facades"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
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

// Menu 菜单模型
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`                                         // 主键ID
	CreatedAt time.Time      `json:"createdAt"`                                                    // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`                                                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                               // 删除时间
	ParentID  uint           `gorm:"default:0;comment:父菜单ID" json:"parentId"`                      // 父菜单ID
	Name      string         `gorm:"type:varchar(50);not null;comment:菜单名称" json:"name"`           // 菜单名称
	Path      string         `gorm:"type:varchar(100);comment:路由路径" json:"path"`                   // 路由路径
	Component string         `gorm:"type:varchar(100);comment:组件路径" json:"component"`              // 组件路径
	Redirect  string         `gorm:"type:varchar(100);comment:重定向路径" json:"redirect"`              // 重定向路径
	Icon      string         `gorm:"type:varchar(50);comment:图标" json:"icon"`                      // 图标
	Sort      uint           `gorm:"default:0;comment:排序" json:"sort"`                             // 排序
	Hidden    bool           `gorm:"default:false;comment:是否隐藏" json:"hidden"`                     // 是否隐藏
	KeepAlive bool           `gorm:"default:false;comment:是否缓存" json:"keepAlive"`                  // 是否缓存
	Type      uint           `gorm:"type:tinyint(1);default:1;comment:类型 1:菜单 2:按钮" json:"type"`   // 类型
	Status    uint           `gorm:"type:tinyint(1);default:1;comment:状态 1:启用 2:禁用" json:"status"` // 状态
	Roles     []*Role        `gorm:"many2many:role_menus;" json:"roles"`                           // 菜单角色关联
	Children  []*Menu        `gorm:"-" json:"children"`                                            // 子菜单
}

// 用户角色中间表
type UserRole struct {
	UserID uint `gorm:"primarykey;comment:用户ID" json:"userId"`
	RoleID uint `gorm:"primarykey;comment:角色ID" json:"roleId"`
}

// 角色菜单中间表
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
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 创建默认管理员角色
	adminRole := &Role{
		Name:   "超级管理员",
		Code:   "admin",
		Sort:   1,
		Status: 1,
		Remark: "系统默认创建的超级管理员角色",
	}
	if err := db.Create(adminRole).Error; err != nil {
		return err
	}

	// 创建默认菜单
	dashboardMenu := &Menu{
		ParentID:  0,
		Name:      "仪表盘",
		Path:      "/dashboard",
		Component: "LAYOUT",
		Icon:      "DashboardOutlined",
		Sort:      1,
		Type:      1,
		Status:    1,
	}
	if err := db.Create(dashboardMenu).Error; err != nil {
		return err
	}

	// 关联角色和菜单
	if err := db.Create(&RoleMenu{
		RoleID: adminRole.ID,
		MenuID: dashboardMenu.ID,
	}).Error; err != nil {
		return err
	}

	// 创建默认管理员账号
	hashedPassword, err := HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := &User{
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

// GetRoles 获取用户角色列表
func (u *User) GetRoles() []string {
	// 查询用户角色
	var role Role
	if err := facades.DB().First(&role, u.RoleID).Error; err != nil {
		return []string{}
	}
	return []string{role.Code}
}
