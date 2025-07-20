package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/utils"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		DB: facades.DB(),
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	var user models.User

	// 查询用户
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户不存在")
		}
		return nil, "", err
	}

	// 验证密码
	if !models.CheckPassword(password, user.Password) {
		return nil, "", errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, "", errors.New("用户已被禁用")
	}

	// 查询用户角色
	var role models.Role
	if err := s.DB.First(&role, user.RoleID).Error; err != nil {
		return nil, "", errors.New("获取用户角色失败")
	}

	// 生成访问令牌
	token, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, "", err
	}

	// 更新最后登录时间
	s.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": time.Now(),
	})

	return &user, token, nil
}

// GenerateAccessToken 生成访问令牌
func (s *AuthService) GenerateAccessToken(user *models.User) (string, error) {
	// 从配置获取令牌过期时间，默认7天
	expiresDays := 7
	if value := facades.Config().Get("jwt.token_expire_days"); value != nil {
		if v, ok := value.(int); ok {
			expiresDays = v
		}
	}

	// 生成声明
	claims := models.CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RealName: user.RealName,
		UUID:     user.UUID.String(),
		RoleID:   user.RoleID,
		Roles:    user.GetRoles(),
		HomePath: "/dashboard",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(expiresDays))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成令牌
	return utils.GenerateToken(claims)
}

// GenerateRefreshToken 生成刷新令牌
func (s *AuthService) GenerateRefreshToken(user *models.User) (string, error) {
	// 从配置获取刷新令牌过期时间，默认30天
	expiresDays := 30
	if value := facades.Config().Get("jwt.refresh_token_expire_days"); value != nil {
		if v, ok := value.(int); ok {
			expiresDays = v
		}
	}

	// 生成声明
	claims := models.CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		UUID:     user.UUID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(expiresDays))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成令牌
	return utils.GenerateToken(claims)
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// 解析刷新令牌
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 查询用户
	var user models.User
	if err := s.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return "", err
	}

	// 生成新的访问令牌
	return s.GenerateAccessToken(&user)
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserAccessCodes 获取用户权限码
func (s *AuthService) GetUserAccessCodes(userID uint) ([]string, error) {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 这里简单返回用户角色作为权限码，实际项目中可能需要更复杂的权限计算
	var role models.Role
	if err := s.DB.First(&role, user.RoleID).Error; err != nil {
		return nil, err
	}

	// 返回一些示例权限码
	accessCodes := []string{"AC_100010"}

	// 管理员添加额外权限码
	if role.Code == "admin" {
		accessCodes = append(accessCodes, "AC_100020", "AC_100030")
	}

	// 超级管理员拥有全部权限
	if role.Code == "super" {
		accessCodes = append(accessCodes, "AC_100100", "AC_100110", "AC_100120")
	}

	return accessCodes, nil
}

// GetUserMenus 获取用户菜单
func (s *AuthService) GetUserMenus(userID uint) ([]dto.MenuResponse, error) {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 根据角色获取菜单
	var role models.Role
	if err := s.DB.Preload("Menus").First(&role, user.RoleID).Error; err != nil {
		return nil, err
	}

	// 构建菜单树
	menuTree := buildMenuTree(role.Menus)

	// 转换为前端所需格式
	return convertToMenuResponse(menuTree), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []*models.Menu) []*models.Menu {
	// 构建菜单映射
	menuMap := make(map[uint]*models.Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	// 构建菜单树
	var rootMenus []*models.Menu
	for _, menu := range menus {
		if parentMenu, exists := menuMap[menu.ParentID]; exists && menu.ParentID != 0 {
			if parentMenu.Children == nil {
				parentMenu.Children = []*models.Menu{}
			}
			parentMenu.Children = append(parentMenu.Children, menu)
		} else if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
		}
	}

	return rootMenus
}

// convertToMenuResponse 将菜单树转换为前端所需格式
func convertToMenuResponse(menus []*models.Menu) []dto.MenuResponse {
	var result []dto.MenuResponse
	for _, menu := range menus {
		menuResponse := dto.MenuResponse{
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: dto.MenuMeta{
				Title:     menu.Name,
				Icon:      menu.Icon,
				KeepAlive: menu.KeepAlive,
			},
		}

		if len(menu.Children) > 0 {
			menuResponse.Children = convertToMenuResponse(menu.Children)
		}

		result = append(result, menuResponse)
	}

	return result
}
