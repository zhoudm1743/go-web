package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/admin/dto"
	"github.com/zhoudm1743/go-web/apps/admin/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
)

// MenuController 菜单控制器
type MenuController struct{}

// NewMenuController 创建菜单控制器
func NewMenuController() *MenuController {
	return &MenuController{}
}

// GetMenus 获取菜单列表
func (c *MenuController) GetMenus(ctx *gin.Context) {
	var menus []models.Menu
	db := facades.DB()

	// 查询条件
	title := ctx.Query("title")
	name := ctx.Query("name")
	path := ctx.Query("path")

	// 构建查询
	query := db.Order("`order` asc")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}

	if err := query.Find(&menus).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithData(ctx, menus)
}

// CreateMenu 创建菜单
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var req dto.MenuCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 检查菜单名称是否已存在
	var count int64
	db := facades.DB()
	if err := db.Model(&models.Menu{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	if count > 0 {
		response.FailWithMsg(ctx, response.Failed, "菜单名称已存在")
		return
	}

	// 创建菜单
	menu := &models.Menu{}
	response.Copy(menu, req)
	menu.Status = 1 // 默认启用

	if err := db.Create(menu).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "创建成功")
}

// UpdateMenu 更新菜单
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	var req dto.MenuUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()
	var menu models.Menu
	if err := db.First(&menu, req.ID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "菜单不存在")
		return
	}

	// 检查菜单名称是否已被其他菜单使用
	if req.Name != "" && req.Name != menu.Name {
		var count int64
		if err := db.Model(&models.Menu{}).Where("name = ? AND id != ?", req.Name, req.ID).Count(&count).Error; err != nil {
			response.Fail(ctx, response.SystemError)
			return
		}

		if count > 0 {
			response.FailWithMsg(ctx, response.Failed, "菜单名称已存在")
			return
		}
	}

	// 检查是否将自己设为了自己的父级
	if req.PID != nil && *req.PID == req.ID {
		response.FailWithMsg(ctx, response.Failed, "不能将菜单设为自己的父级")
		return
	}

	// 更新菜单
	response.Copy(&menu, req)

	if err := db.Save(&menu).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}

// DeleteMenu 删除菜单
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	id := ctx.Param("id")
	menuID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "菜单ID无效")
		return
	}

	db := facades.DB()

	// 检查是否有子菜单
	var childCount int64
	if err := db.Model(&models.Menu{}).Where("parent_id = ?", menuID).Count(&childCount).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	if childCount > 0 {
		response.FailWithMsg(ctx, response.Failed, "该菜单下有子菜单，请先删除子菜单")
		return
	}

	// 删除菜单
	if err := db.Delete(&models.Menu{}, menuID).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 同时删除角色菜单关联
	if err := db.Where("menu_id = ?", menuID).Delete(&models.RoleMenu{}).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}
