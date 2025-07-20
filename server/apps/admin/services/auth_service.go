package services

import (
	"errors"
	"time"

	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/utils"
)

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*models.Admin, string, error) {
	var admin models.Admin
	db := facades.DB()

	// 查询管理员信息
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, "", errors.New("用户名不存在")
	}

	// 校验密码
	if !models.CheckPassword(password, admin.Password) {
		return nil, "", errors.New("密码错误")
	}

	// 检查管理员状态
	if admin.Status != 1 {
		return nil, "", errors.New("账号已被禁用")
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(int(admin.ID), admin.Username, int(admin.RoleID))
	if err != nil {
		return nil, "", err
	}

	// 更新登录信息
	admin.LastLoginAt = time.Now()
	db.Save(&admin)

	return &admin, token, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID int) (*models.Admin, error) {
	var admin models.Admin
	db := facades.DB()

	if err := db.First(&admin, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &admin, nil
}

// GetUserAccessCodes 获取用户权限码
func (s *AuthService) GetUserAccessCodes(userID int) ([]string, error) {
	// 这里可以根据实际需求实现更复杂的权限码逻辑
	// 暂时简单返回一些固定权限码
	return []string{"auth.user", "auth.role"}, nil
}

// GetUserMenus 获取用户菜单
func (s *AuthService) GetUserMenus(userID int) ([]models.Menu, error) {
	var admin models.Admin
	db := facades.DB()

	// 查询用户信息
	if err := db.First(&admin, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 查询用户角色
	var role models.Role
	if err := db.First(&role, admin.RoleID).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 超级管理员直接返回所有菜单
	if role.Code == "super" || role.ID == 1 {
		var allMenus []models.Menu
		if err := db.Where("status = 1").Order("`order` ASC").Find(&allMenus).Error; err != nil {
			return nil, err
		}
		return allMenus, nil
	}

	// 查询角色菜单
	var menus []models.Menu
	if err := db.Raw(`
		SELECT m.* FROM menus m 
		INNER JOIN role_menus rm ON m.id = rm.menu_id 
		WHERE rm.role_id = ? AND m.status = 1
		ORDER BY m.`+"`order`"+` ASC
	`, role.ID).Scan(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

// GetAllMenus 获取所有菜单
func (s *AuthService) GetAllMenus() ([]models.Menu, error) {
	var menus []models.Menu
	db := facades.DB()

	// 查询所有菜单
	if err := db.Order("`order` asc").Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}
