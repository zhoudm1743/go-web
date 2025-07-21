# GoWeb 代码生成器

这是一个为GoWeb项目设计的代码生成器，可以快速生成CRUD相关的后端和前端代码。

## 功能特性

- 生成模型（Model）
- 生成数据传输对象（DTO）
- 生成控制器（Controller）
- 生成路由（Route）
- 生成前端页面（Vue 3 + TypeScript）
- 生成API接口文件
- 支持关系字段（BelongsTo, HasOne, HasMany, ManyToMany）
- 支持生成历史记录和回滚功能
- 支持自动生成前端表单、表格和API调用

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

#### 基本字段

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

#### 关系字段

关系字段的格式：

```
字段名:关系类型:关联模型:外键:引用字段:描述:是否预加载:多对多关联表
```

例如：

```
"User:belongs_to:User:UserID:ID:用户:true:,Comments:has_many:Comment:ArticleID:ID:评论:false:,Tags:many_to_many:Tag:article_id:tag_id:标签:true:article_tags"
```

说明：

1. `User` - 字段名
2. `belongs_to` - 关系类型（可选值：belongs_to, has_one, has_many, many_to_many）
3. `User` - 关联模型名称
4. `UserID` - 外键字段名
5. `ID` - 引用字段名
6. `用户` - 字段描述
7. `true` - 是否预加载
8. `article_tags` - 多对多关联表名（仅多对多关系需要）

支持的关系类型：
- `belongs_to`: 从属于（例如：文章从属于用户）
- `has_one`: 拥有一个（例如：用户拥有一个个人资料）
- `has_many`: 拥有多个（例如：用户拥有多篇文章）
- `many_to_many`: 多对多（例如：文章和标签的多对多关系）

## 生成的文件

代码生成器会生成以下文件：

1. 模型文件：`server/apps/admin/models/{name}.go`
2. DTO文件：`server/apps/admin/dto/{name}.go`
3. 控制器文件：`server/apps/admin/controllers/{name}_controller.go`
4. 更新路由文件：`server/apps/admin/routes/routes.go`
5. 前端页面：`front-end/src/views/setting/{name}/index.vue`
6. 前端组件：`front-end/src/views/setting/{name}/components/TableModal.vue`
7. 前端API文件：`front-end/src/service/api/{name}.ts`

## 历史记录和回滚

代码生成器会记录每次生成操作，支持回滚功能。

### 历史记录

每次代码生成操作都会在数据库中记录以下信息：
- 结构体名称和表名
- 包名和描述
- 生成的文件路径
- 字段配置信息
- 生成时间

### 回滚功能

回滚功能可以撤销之前的代码生成操作，包括：
- 删除生成的文件（会先备份到临时目录）
- 回滚数据库表（如果指定）
- 删除API相关配置（如果指定）
- 删除菜单项（如果指定）

### 回滚命令

```bash
./codegen -rollback <历史记录ID> [-deleteFiles] [-deleteTable] [-deleteAPI] [-deleteMenu]
```

参数说明：
- `<历史记录ID>`: 要回滚的历史记录ID
- `-deleteFiles`: 是否删除生成的文件（默认true）
- `-deleteTable`: 是否删除数据库表（默认false）
- `-deleteAPI`: 是否删除API配置（默认false）
- `-deleteMenu`: 是否删除菜单（默认false）

## 示例

### 基本示例：生成文章模块

```bash
./codegen -struct Article -desc "文章" -fields "Title:string:title:标题:true:false:true:true:true,Content:string:content:内容:true:false:false:false:false,Author:string:author:作者:true:false:true:true:false,Status:uint:status:状态:true:false:false:true:false"
```

### 产品模块

```bash
./codegen -struct Product -desc "产品" -fields "Name:string:name:名称:true:false:true:true:true,Price:float64:price:价格:true:false:false:true:true,Description:string:description:描述:false:false:false:false:false,Status:uint:status:状态:true:false:false:true:false"
```

### 带关系的文章模块

```bash
./codegen -struct Article -desc "文章" -fields "Title:string:title:标题:true:false:true:true:true,Content:string:content:内容:true:false:false:false:false" -relations "User:belongs_to:User:UserID:ID:用户:true:,Comments:has_many:Comment:ArticleID:ID:评论列表:false:,Tags:many_to_many:Tag:article_id:tag_id:标签列表:true:article_tags"
```

## API控制器

代码生成器也可以集成到Web应用中，通过API调用进行代码生成。以下是可用的API端点：

- `GET /admin/codegen/apps` - 获取应用列表
- `GET /admin/codegen/tables?db=xxx` - 获取数据库表列表
- `GET /admin/codegen/columns?db=xxx&table=xxx` - 获取表列信息
- `POST /admin/codegen/generate` - 生成代码
- `GET /admin/codegen/history?page=1&pageSize=10` - 获取生成历史
- `POST /admin/codegen/rollback/:id` - 回滚指定的生成记录

## 常见问题

1. **ID字段重复**: 如果在字段定义中包含了ID字段，会与基础模型的ID字段冲突。建议不要在字段中显式定义ID字段。

2. **路径问题**: 确保路径分隔符在不同操作系统上正确，代码生成器内部使用`filepath.Join`来确保跨平台兼容。

3. **数据库连接**: 确保有可用的数据库连接，否则历史记录功能将不可用。 