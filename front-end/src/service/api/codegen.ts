import { request } from '../http';

// 表结构信息
export interface TableInfo {
  tableName: string;
  tableComment: string;
}

// 字段结构信息
export interface ColumnInfo {
  columnName: string;
  dataType: string;
  columnComment: string;
  isNullable: string;
  columnKey: string;
}

// 关系类型
export enum RelationType {
  BelongsTo = 'belongs_to',   // 从属于
  HasOne = 'has_one',         // 拥有一个
  HasMany = 'has_many',       // 拥有多个
  ManyToMany = 'many_to_many' // 多对多
}

// 字段结构
export interface FieldInfo {
  fieldName: string;    // 结构体字段名称
  fieldType: string;    // Go类型
  columnName: string;   // 数据库字段名
  fieldDesc: string;    // 字段描述
  required: boolean;    // 是否必填
  isPrimaryKey: boolean;// 是否主键
  isSearchable: boolean;// 是否可搜索
  isFilterable: boolean;// 是否可过滤
  isSortable: boolean;  // 是否可排序
  
  // 关系字段
  isRelation?: boolean;         // 是否关系字段
  relationType?: RelationType;  // 关系类型
  relatedModel?: string;        // 关联模型
  foreignKey?: string;          // 外键 (本模型字段)
  references?: string;          // 引用字段 (关联模型字段)
  preload?: boolean;            // 是否预加载
  joinTable?: string;           // 多对多关联表名
  
  // JOIN查询相关
  joinable?: boolean;           // 是否支持JOIN查询
  joinCondition?: string;       // JOIN条件字段
  filterCondition?: string;     // 过滤条件字段
}

// 代码生成配置
export interface CodegenConfig {
  structName: string;     // 结构体名称
  tableName: string;      // 表名
  packageName: string;    // 包名
  description: string;    // 描述
  apiPrefix: string;      // API前缀
  appName: string;        // 应用名称
  hasList: boolean;       // 是否有列表
  hasCreate: boolean;     // 是否有创建
  hasUpdate: boolean;     // 是否有更新
  hasDelete: boolean;     // 是否有删除
  hasDetail: boolean;     // 是否有详情
  hasPagination: boolean; // 是否分页
  fields: FieldInfo[];    // 字段列表
}

// 历史记录
export interface HistoryRecord {
  id: number;
  createdAt: string;
  updatedAt: string;
  table: string;
  structName: string;
  packageName: string;
  moduleName: string;
  description: string;
  flag: number;
}

// 获取应用列表
export function getAppList() {
  return request.Get<Service.ResponseResult<string[]>>('/admin/codegen/apps');
}

// 获取数据库表列表
export function getTables() {
  return request.Get<Service.ResponseResult<TableInfo[]>>('/admin/codegen/tables');
}

// 获取表字段列表
export function getColumns(tableName: string) {
  return request.Get<Service.ResponseResult<ColumnInfo[]>>('/admin/codegen/columns', { params: { tableName } });
}

// 生成代码
export function generateCode(data: CodegenConfig) {
  return request.Post<Service.ResponseResult<void>>('/admin/codegen/generate', data);
}

// 获取历史记录列表
export function getHistoryList(page: number, pageSize: number) {
  return request.Get<Service.ResponseResult<{
    list: HistoryRecord[];
    total: number;
    page: number;
    pageSize: number;
  }>>('/admin/codegen/history', { params: { page, pageSize } });
}

// 回滚代码生成
export function rollbackCode(data: {
  id: number;
  deleteFiles: boolean;
  deleteApi: boolean;
  deleteMenu: boolean;
  deleteTable: boolean;
}) {
  return request.Post<Service.ResponseResult<void>>('/admin/codegen/rollback', data);
}

// 删除历史记录
export function deleteHistory(id: number) {
  return request.Delete<Service.ResponseResult<void>>(`/admin/codegen/history/${id}`);
} 