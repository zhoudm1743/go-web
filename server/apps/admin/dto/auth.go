package dto

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	RealName    string   `json:"realName"`
	Roles       []string `json:"roles"`
	AccessToken string   `json:"accessToken"`
	HomePath    string   `json:"homePath,omitempty"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Roles    []string `json:"roles"`
	HomePath string   `json:"homePath,omitempty"`
}

// MenuItem 菜单项
type MenuItem struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Component string     `json:"component,omitempty"`
	Redirect  string     `json:"redirect,omitempty"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuItem `json:"children,omitempty"`
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Icon                     string   `json:"icon,omitempty"`
	Title                    string   `json:"title"`
	HideMenu                 bool     `json:"hideMenu,omitempty"`
	HideChildrenInMenu       bool     `json:"hideChildrenInMenu,omitempty"`
	HideTab                  bool     `json:"hideTab,omitempty"`
	IgnoreAuth               bool     `json:"ignoreAuth,omitempty"`
	IgnoreKeepAlive          bool     `json:"ignoreKeepAlive,omitempty"`
	AffixTab                 bool     `json:"affixTab,omitempty"`
	Order                    int      `json:"order,omitempty"`
	FrameSrc                 string   `json:"frameSrc,omitempty"`
	CarryParam               bool     `json:"carryParam,omitempty"`
	SingleLayout             string   `json:"singleLayout,omitempty"`
	KeepAlive                bool     `json:"keepAlive,omitempty"`
	Authority                []string `json:"authority,omitempty"`
	MenuVisibleWithForbidden bool     `json:"menuVisibleWithForbidden,omitempty"`
}

// MenuResponse 菜单响应别名
type MenuResponse = MenuItem
