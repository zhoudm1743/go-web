# GoWeb 代码生成器

这是一个为GoWeb项目设计的代码生成器，可以快速生成CRUD相关的后端和前端代码。

## 功能特性

- 生成模型（Model）
- 生成数据传输对象（DTO）
- 生成控制器（Controller）
- 生成路由（Route）
- 生成前端页面（Vue 3 + TypeScript）
- 生成API接口文件

## 安装

1. 确保您的GoWeb项目已经设置好
2. 构建代码生成器：

```bash
cd server/pkg/generator/cmd
go build -o codegen
```

## 使用方法

### 基本命令

```bash
./codegen -struct User -desc "用户" -fields "Username:string:username:用户名:true:false:true:true:false,Password:string:password:密码:true:false:false:false:false,Email:string:email:邮箱:false:false:true:true:false"
```

### 参数说明

| 参数 | 描述 | 默认值 | 是否必填 |
| ---- | ---- | ------ | -------- |
| -struct | 结构体名称 | - | 是 |
| -table | 表名 | 结构体名称的蛇形命名 | 否 |
| -package | 包名 | admin | 否 |
| -desc | 描述 | - | 是 |
| -module | 模块名 | github.com/zhoudm1743/go-web | 否 |
| -router | 路由分组 | privateRoutes | 否 |
| -api | API前缀 | 表名 | 否 |
| -root | 项目根目录 | ./ | 否 |
| -fields | 字段定义 | - | 是 |
| -all | 是否生成所有CRUD操作 | true | 否 |
| -list | 是否有列表 | false | 否 |
| -create | 是否有创建 | false | 否 |
| -update | 是否有更新 | false | 否 |
| -delete | 是否有删除 | false | 否 |
| -detail | 是否有详情 | false | 否 |
| -pagination | 是否分页 | true | 否 |

### 字段定义格式

字段定义使用逗号分隔不同字段，每个字段的格式为：

```
字段名:类型:数据库字段名:描述:是否必填:是否主键:是否可搜索:是否可过滤:是否可排序
```

例如：

```
"Username:string:username:用户名:true:false:true:true:false,Email:string:email:邮箱:false:false:true:true:false"
```

说明：

1. `Username` - 字段名
2. `string` - 字段类型
3. `username` - 数据库字段名
4. `用户名` - 字段描述
5. `true` - 是否必填（创建时）
6. `false` - 是否为主键
7. `true` - 是否可搜索
8. `true` - 是否可过滤
9. `false` - 是否可排序

## 生成的文件

代码生成器会生成以下文件：

1. 模型文件：`server/apps/admin/models/{name}.go`
2. DTO文件：`server/apps/admin/dto/{name}.go`
3. 控制器文件：`server/apps/admin/controllers/{name}_controller.go`
4. 更新路由文件：`server/apps/admin/routes/routes.go`
5. 前端页面：`front-end/src/views/setting/{name}/index.vue`
6. 前端组件：`front-end/src/views/setting/{name}/components/TableModal.vue`
7. 前端API文件：`front-end/src/service/api/{name}.ts`

## 示例

### 生成文章模块

```bash
./codegen -struct Article -desc "文章" -fields "Title:string:title:标题:true:false:true:true:true,Content:string:content:内容:true:false:false:false:false,Author:string:author:作者:true:false:true:true:false,Status:uint:status:状态:true:false:false:true:false"
```

### 生成产品模块

```bash
./codegen -struct Product -desc "产品" -fields "Name:string:name:名称:true:false:true:true:true,Price:float64:price:价格:true:false:false:true:true,Description:string:description:描述:false:false:false:false:false,Status:uint:status:状态:true:false:false:true:false"
``` 