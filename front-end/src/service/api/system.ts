import { request } from '../http'

// 获取所有路由信息
export function fetchAllRoutes() {
  return request.Get<Service.ResponseResult<AppRoute.RowRoute[]>>('/getAllRoutes')
}

// 获取菜单列表
export function fetchMenus() {
  return request.Get<Service.ResponseResult<AppRoute.RowRoute[]>>('/admin/menu/list')
}

// 创建菜单
export function createMenu(data: Partial<AppRoute.RowRoute>) {
  return request.Post<Service.ResponseResult<any>>('/admin/menu/create', data)
}

// 更新菜单
export function updateMenu(data: Partial<AppRoute.RowRoute>) {
  return request.Put<Service.ResponseResult<any>>('/admin/menu/update', data)
}

// 删除菜单
export function deleteMenu(id: number) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/menu/delete/${id}`)
}

// 获取所有管理员信息
export function fetchAdminPage() {
  return request.Get<Service.ResponseResult<Entity.Admin[]>>('/user/list')
}

// 创建管理员
export function createAdmin(data: Partial<Entity.Admin>) {
  return request.Post<Service.ResponseResult<any>>('/admin/create', data)
}

// 更新管理员
export function updateAdmin(data: Partial<Entity.Admin>) {
  return request.Put<Service.ResponseResult<any>>('/admin/update', data)
}

// 删除管理员
export function deleteAdmin(id: number) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/delete/${id}`)
}

// 获取所有角色列表
export function fetchRoleList() {
  return request.Get<Service.ResponseResult<Entity.Role[]>>('/role/list')
}

// 获取角色列表 - 完整信息
export function fetchRoles() {
  return request.Get<Service.ResponseResult<Entity.Role[]>>('/role/list')
}

// 创建角色
export function createRole(data: Partial<Entity.Role>) {
  return request.Post<Service.ResponseResult<any>>('/role/create', data)
}

// 更新角色
export function updateRole(data: Partial<Entity.Role>) {
  return request.Put<Service.ResponseResult<any>>('/role/update', data)
}

// 删除角色
export function deleteRole(id: number) {
  return request.Delete<Service.ResponseResult<any>>(`/role/delete/${id}`)
}

// 获取角色菜单
export function getRoleMenus(roleId: number) {
  return request.Get<Service.ResponseResult<number[]>>('/role/menus', { params: { roleId } })
}

// 更新角色菜单
export function updateRoleMenus(data: { roleId: number; menuIds: number[] }) {
  return request.Post<Service.ResponseResult<any>>('/role/menus', data)
}

/**
 * 请求获取字典列表
 *
 * @param code - 字典编码，用于筛选特定的字典列表
 * @returns 返回的字典列表数据
 */
export function fetchDictList(code?: string) {
  const params = { code }
  return request.Get<Service.ResponseResult<Entity.Dict[]>>('/dict/list', { params })
}

// 创建字典
export function createDict(data: Partial<Entity.Dict>) {
  return request.Post<Service.ResponseResult<any>>('/dict/create', data)
}

// 更新字典
export function updateDict(data: Partial<Entity.Dict>) {
  return request.Put<Service.ResponseResult<any>>('/dict/update', data)
}

// 删除字典
export function deleteDict(id: number) {
  return request.Delete<Service.ResponseResult<any>>(`/dict/delete/${id}`)
}
