package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/services"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/core/utils"
)

// AuthController 认证控制器
type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController() *AuthController {
	return &AuthController{
		AuthService: services.NewAuthService(),
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	admin, token, err := c.AuthService.Login(req.Username, req.Password)
	if err != nil {
		response.FailWithMsg(ctx, response.LoginAccountError, err.Error())
		return
	}

	// 使用Copy函数构造响应
	loginResp := &dto.LoginResponse{}
	response.Copy(loginResp, admin)
	loginResp.AccessToken = token
	loginResp.Roles = admin.GetRoles()

	response.OkWithData(ctx, loginResp)
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// JWT无状态，服务端只需清除cookie
	response.Ok(ctx)
}

// GetUserInfo 获取用户信息
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	// 从JWT中获取用户信息
	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, response.TokenInvalid)
		return
	}

	admin, err := c.AuthService.GetUserInfo(int(claims.UserID))
	if err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 使用Copy函数返回用户信息
	userInfo := &dto.AdminInfoResponse{}
	response.Copy(userInfo, admin)
	userInfo.Roles = admin.GetRoles()

	response.OkWithData(ctx, userInfo)
}

// GetAccessCodes 获取用户权限码
func (c *AuthController) GetAccessCodes(ctx *gin.Context) {
	// 从JWT中获取用户信息
	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, response.TokenInvalid)
		return
	}

	codes, err := c.AuthService.GetUserAccessCodes(int(claims.UserID))
	if err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, codes)
}

// GetUserRoutes 获取用户路由菜单
func (c *AuthController) GetUserRoutes(ctx *gin.Context) {
	// 从JWT中获取用户信息
	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, response.TokenInvalid)
		return
	}

	menus, err := c.AuthService.GetUserMenus(int(claims.UserID))
	if err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, menus)
}

// GetAllRoutes 获取所有路由菜单（管理菜单页面使用）
func (c *AuthController) GetAllRoutes(ctx *gin.Context) {
	menus, err := c.AuthService.GetAllMenus()
	if err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, menus)
}
