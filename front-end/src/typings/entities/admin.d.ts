/// <reference path="../global.d.ts"/>

/** 管理员数据库表字段 */
namespace Entity {
  interface Admin {
    /** ID */
    id?: number
    /** 用户名 */
    username?: string
    /** 昵称 */
    nickname?: string
    /** 真实姓名 */
    realName?: string
    /** 头像 */
    avatar?: string
    /** 邮箱 */
    email?: string
    /** 手机号 */
    mobile?: string
    /** 状态：1启用 2禁用 */
    status?: number
    /** 角色ID */
    roleId?: number
    /** 最后登录时间 */
    lastLoginAt?: string
    /** 最后登录IP */
    lastLoginIp?: string
    /** 创建时间 */
    createdAt?: string
    /** 更新时间 */
    updatedAt?: string
  }
} 