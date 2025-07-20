/// <reference path="../global.d.ts"/>

/* 角色数据库表字段 */
declare namespace Entity {
  /** 角色类型 */
  type RoleType = 'super' | 'admin' | 'user' | string

  /** 角色实体 */
  interface Role {
    /** 角色ID */
    id: number
    /** 角色名称 */
    name: string
    /** 角色编码 */
    code: string
    /** 排序值 */
    sort: number
    /** 状态：1-启用，2-禁用 */
    status: number
    /** 备注 */
    remark?: string
    /** 创建时间 */
    createdAt?: string
    /** 更新时间 */
    updatedAt?: string
  }
}
