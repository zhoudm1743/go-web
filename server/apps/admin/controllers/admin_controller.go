package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
)

// AdminController 管理员控制器
type AdminController struct{}

// NewAdminController 创建管理员控制器
func NewAdminController() *AdminController {
	return &AdminController{}
}

// GetAdmins 获取管理员列表
func (c *AdminController) GetAdmins(ctx *gin.Context) {
	var admins []models.Admin
	db := facades.DB()

	if err := db.Find(&admins).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, admins)
}

// CreateAdmin 创建管理员
func (c *AdminController) CreateAdmin(ctx *gin.Context) {
	var req dto.AdminCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 检查管理员名是否已存在
	var count int64
	db := facades.DB()
	if err := db.Model(&models.Admin{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	if count > 0 {
		response.FailWithMsg(ctx, response.Failed, "管理员名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 创建管理员
	admin := &models.Admin{}
	response.Copy(admin, req)
	admin.Password = hashedPassword

	// 确保状态值有效
	if admin.Status == 0 {
		admin.Status = 1 // 默认启用
	}

	if err := db.Create(admin).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "创建成功")
}

// UpdateAdmin 更新管理员
func (c *AdminController) UpdateAdmin(ctx *gin.Context) {
	var req dto.AdminUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()
	var admin models.Admin
	if err := db.First(&admin, req.ID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "管理员不存在")
		return
	}

	// 检查管理员名是否已被其他管理员使用
	if req.Username != "" && req.Username != admin.Username {
		var count int64
		if err := db.Model(&models.Admin{}).Where("username = ? AND id != ?", req.Username, req.ID).Count(&count).Error; err != nil {
			response.Fail(ctx, response.SystemError)
			return
		}

		if count > 0 {
			response.FailWithMsg(ctx, response.Failed, "管理员名已存在")
			return
		}
	}

	// 只更新提供的字段
	updates := map[string]interface{}{}

	// 创建一个临时对象，用于复制非空字段
	tempAdmin := &models.Admin{}
	response.Copy(tempAdmin, req)

	if req.Username != "" {
		updates["username"] = tempAdmin.Username
	}
	if req.Nickname != "" {
		updates["nickname"] = tempAdmin.Nickname
	}
	if req.RealName != "" {
		updates["real_name"] = tempAdmin.RealName
	}
	if req.Email != "" {
		updates["email"] = tempAdmin.Email
	}
	if req.Mobile != "" {
		updates["mobile"] = tempAdmin.Mobile
	}
	if req.RoleID > 0 {
		updates["role_id"] = tempAdmin.RoleID
	}
	if req.Status > 0 {
		updates["status"] = tempAdmin.Status
	}

	if err := db.Model(&admin).Updates(updates).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}

// DeleteAdmin 删除管理员
func (c *AdminController) DeleteAdmin(ctx *gin.Context) {
	id := ctx.Param("id")
	AdminID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "管理员ID无效")
		return
	}

	db := facades.DB()
	if err := db.Delete(&models.Admin{}, AdminID).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}
