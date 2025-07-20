package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
)

// UserController 用户管理控制器
type UserController struct{}

// List 获取用户列表
func (c *UserController) List(ctx *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取数据库连接
	db := facades.DB()
	if db == nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 查询用户列表
	var users []models.User
	var total int64

	// 计算总数
	db.Model(&models.User{}).Count(&total)

	// 查询带分页的数据
	offset := (page - 1) * pageSize
	result := db.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		response.FailWithMsg(ctx, response.SystemError, "获取用户列表失败: "+result.Error.Error())
		return
	}

	// 返回用户列表
	response.OkWithData(ctx, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      users,
	})
}

// GetByID 根据ID获取用户
func (c *UserController) GetByID(ctx *gin.Context) {
	// 获取用户ID
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || userID <= 0 {
		response.FailWithMsg(ctx, response.ParamsValidError, "无效的用户ID")
		return
	}

	// 获取数据库连接
	db := facades.DB()
	if db == nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 查询用户
	var user models.User
	result := db.First(&user, userID)
	if result.Error != nil {
		response.Fail(ctx, response.AssertArgumentError.Make("用户不存在"))
		return
	}

	// 返回用户信息
	response.OkWithData(ctx, user)
}

// Update 更新用户信息
func (c *UserController) Update(ctx *gin.Context) {
	// 获取用户ID
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || userID <= 0 {
		response.FailWithMsg(ctx, response.ParamsValidError, "无效的用户ID")
		return
	}

	// 获取更新数据
	var userData struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
		RoleID   uint   `json:"roleId"`
		Status   int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, err.Error())
		return
	}

	// 获取数据库连接
	db := facades.DB()
	if db == nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 查找用户是否存在
	var user models.User
	result := db.First(&user, userID)
	if result.Error != nil {
		response.Fail(ctx, response.AssertArgumentError.Make("用户不存在"))
		return
	}

	// 更新用户信息
	updates := map[string]interface{}{}
	if userData.Nickname != "" {
		updates["nickname"] = userData.Nickname
	}
	if userData.Email != "" {
		updates["email"] = userData.Email
	}
	if userData.Avatar != "" {
		updates["avatar"] = userData.Avatar
	}
	if userData.RoleID > 0 {
		updates["role_id"] = userData.RoleID
	}
	if userData.Status != 0 {
		updates["status"] = userData.Status
	}

	if len(updates) > 0 {
		if err := db.Model(&user).Updates(updates).Error; err != nil {
			response.FailWithMsg(ctx, response.SystemError, "更新用户失败: "+err.Error())
			return
		}
	}

	// 返回更新后的用户信息
	response.OkWithData(ctx, user)
}

// Delete 删除用户
func (c *UserController) Delete(ctx *gin.Context) {
	// 获取用户ID
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || userID <= 0 {
		response.FailWithMsg(ctx, response.ParamsValidError, "无效的用户ID")
		return
	}

	// 获取数据库连接
	db := facades.DB()
	if db == nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 确保不会删除最后一个管理员
	var adminCount int64
	db.Model(&models.User{}).Where("role_id = ?", 1).Count(&adminCount) // 假设角色ID 1为超级管理员

	var user models.User
	if db.First(&user, userID).Error == nil && user.RoleID == 1 && adminCount <= 1 {
		response.FailWithMsg(ctx, response.Failed, "无法删除唯一的管理员账号")
		return
	}

	// 删除用户
	if result := db.Delete(&models.User{}, userID); result.Error != nil {
		response.FailWithMsg(ctx, response.SystemError, "删除用户失败: "+result.Error.Error())
		return
	}

	// 返回成功信息
	response.OkWithMsg(ctx, "删除成功")
}
