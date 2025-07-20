package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
	"github.com/zhoudm1743/go-web/core/utils"
)

// RoleController 角色控制器
type RoleController struct{}

// NewRoleController 创建角色控制器
func NewRoleController() *RoleController {
	return &RoleController{}
}

// GetRoles 获取角色列表
func (c *RoleController) GetRoles(ctx *gin.Context) {
	var roles []models.Role
	db := facades.DB()

	if err := db.Find(&roles).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, roles)
}

// GetRoleList 获取角色简易列表（用于下拉选择）
func (c *RoleController) GetRoleList(ctx *gin.Context) {
	var roles []models.Role
	db := facades.DB()

	if err := db.Select("id, name, code").Find(&roles).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, roles)
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req dto.RoleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 检查角色编码是否已存在
	var count int64
	db := facades.DB()
	if err := db.Model(&models.Role{}).Where("code = ?", req.Code).Count(&count).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	if count > 0 {
		response.FailWithMsg(ctx, response.Failed, "角色编码已存在")
		return
	}

	// 创建角色
	role := &models.Role{}
	response.Copy(role, req)

	if err := db.Create(role).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "创建成功")
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	var req dto.RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()
	var role models.Role
	if err := db.First(&role, req.ID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "角色不存在")
		return
	}

	// 检查角色编码是否已被其他角色使用
	if req.Code != "" && req.Code != role.Code {
		var count int64
		if err := db.Model(&models.Role{}).Where("code = ? AND id != ?", req.Code, req.ID).Count(&count).Error; err != nil {
			response.Fail(ctx, response.SystemError)
			return
		}

		if count > 0 {
			response.FailWithMsg(ctx, response.Failed, "角色编码已存在")
			return
		}
	}

	// 更新角色，只更新非空字段
	updates := map[string]interface{}{}

	// 创建一个临时结构体，用于复制非空字段
	tempRole := &models.Role{}
	response.Copy(tempRole, req)

	// 构建更新字段
	if req.Name != "" {
		updates["name"] = tempRole.Name
	}
	if req.Code != "" {
		updates["code"] = tempRole.Code
	}
	if req.Sort > 0 {
		updates["sort"] = tempRole.Sort
	}
	if req.Status > 0 {
		updates["status"] = tempRole.Status
	}
	if req.Remark != "" {
		updates["remark"] = tempRole.Remark
	}

	if err := db.Model(&role).Updates(updates).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "角色ID无效")
		return
	}

	// 检查是否有管理员在使用该角色
	db := facades.DB()
	var count int64
	if err := db.Model(&models.Admin{}).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	if count > 0 {
		response.FailWithMsg(ctx, response.Failed, "该角色正在被使用，无法删除")
		return
	}

	// 删除角色
	if err := db.Delete(&models.Role{}, roleID).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 同时删除角色菜单关联
	if err := db.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}

// GetRoleMenus 获取角色菜单
func (c *RoleController) GetRoleMenus(ctx *gin.Context) {
	id := ctx.Query("roleId")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "角色ID无效")
		return
	}

	// 查询该角色关联的菜单ID
	var menuIDs []uint
	db := facades.DB()
	if err := db.Model(&models.RoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, menuIDs)
}

// UpdateRoleMenus 更新角色菜单
func (c *RoleController) UpdateRoleMenus(ctx *gin.Context) {
	var req struct {
		RoleID  uint   `json:"roleId" binding:"required"`
		MenuIDs []uint `json:"menuIds"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()

	// 开启事务
	tx := db.Begin()

	// 先删除原有的角色菜单关联
	if err := tx.Where("role_id = ?", req.RoleID).Delete(&models.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, response.SystemError)
		return
	}

	// 添加新的角色菜单关联
	if len(req.MenuIDs) > 0 {
		var roleMenus []models.RoleMenu
		for _, menuID := range req.MenuIDs {
			roleMenus = append(roleMenus, models.RoleMenu{
				RoleID: req.RoleID,
				MenuID: menuID,
			})
		}

		if err := tx.Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, response.SystemError)
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, response.SystemError)
		return
	}

	// 更新casbin权限
	enforcer := utils.Casbin()

	// 清除该角色原有权限
	enforcer.RemoveFilteredPolicy(0, strconv.Itoa(int(req.RoleID)))

	// 添加新权限
	// 超级管理员可以访问所有API
	if req.RoleID == 1 {
		enforcer.AddPolicy(strconv.Itoa(int(req.RoleID)), "/*", "*")
	} else {
		// 其他角色只能访问有权限的API
		// 这里根据菜单ID设置权限
		// 实际项目中可能需要更复杂的权限控制
		// 简单示例：为每个菜单ID添加一个API权限
		for _, menuID := range req.MenuIDs {
			enforcer.AddPolicy(strconv.Itoa(int(req.RoleID)), "/api/menu/"+strconv.Itoa(int(menuID)), "GET")
		}

		// 添加一些基础权限
		enforcer.AddPolicy(strconv.Itoa(int(req.RoleID)), "/api/user/info", "GET")
	}

	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}
