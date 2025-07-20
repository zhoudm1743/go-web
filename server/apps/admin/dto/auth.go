package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	ID           uint     `json:"id"`
	Username     string   `json:"username"`
	RealName     string   `json:"realName"`
	Roles        []string `json:"role"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}

// AdminInfoResponse 用户信息响应
type AdminInfoResponse struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Roles    []string `json:"role"`
	Avatar   string   `json:"avatar"`
}

// MenuListRequest 菜单列表请求
type MenuListRequest struct {
	Title string `form:"title"`
	Path  string `form:"path"`
	Name  string `form:"name"`
}

// MenuCreateRequest 创建菜单请求
type MenuCreateRequest struct {
	PID          *uint  `json:"pid"`
	Name         string `json:"name" binding:"required"`
	Path         string `json:"path" binding:"required"`
	Component    string `json:"componentPath"`
	Redirect     string `json:"redirect"`
	Icon         string `json:"icon"`
	Title        string `json:"title" binding:"required"`
	Order        int    `json:"order"`
	Hidden       bool   `json:"hide"`
	KeepAlive    bool   `json:"keepAlive"`
	RequiresAuth bool   `json:"requiresAuth"`
	WithoutTab   bool   `json:"withoutTab"`
	PinTab       bool   `json:"pinTab"`
	MenuType     string `json:"menuType" binding:"required"`
}

// MenuUpdateRequest 更新菜单请求
type MenuUpdateRequest struct {
	ID           uint   `json:"id" binding:"required"`
	PID          *uint  `json:"pid"`
	Name         string `json:"name" binding:"required"`
	Path         string `json:"path" binding:"required"`
	Component    string `json:"componentPath"`
	Redirect     string `json:"redirect"`
	Icon         string `json:"icon"`
	Title        string `json:"title" binding:"required"`
	Order        int    `json:"order"`
	Hidden       bool   `json:"hide"`
	KeepAlive    bool   `json:"keepAlive"`
	RequiresAuth bool   `json:"requiresAuth"`
	WithoutTab   bool   `json:"withoutTab"`
	PinTab       bool   `json:"pinTab"`
	MenuType     string `json:"menuType" binding:"required"`
}

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   uint   `json:"sort"`
	Status uint   `json:"status"`
	Remark string `json:"remark"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Sort   uint   `json:"sort"`
	Status uint   `json:"status"`
	Remark string `json:"remark"`
}

// RoleMenuRequest 角色菜单关联请求
type RoleMenuRequest struct {
	RoleID  uint   `json:"roleId" binding:"required"`
	MenuIDs []uint `json:"menuIds"`
}

// AdminCreateRequest 创建用户请求
type AdminCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	RealName string `json:"realName"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleID   uint   `json:"roleId"`
	Status   uint   `json:"status"`
}

// AdminUpdateRequest 更新用户请求
type AdminUpdateRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RealName string `json:"realName"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleID   uint   `json:"roleId"`
	Status   uint   `json:"status"`
}
