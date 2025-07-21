package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// generateController 生成控制器文件
func (g *Generator) generateController() error {
	// 控制器模板
	const controllerTemplate = `package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-web/apps/{{.PackageName}}/dto"
	"github.com/zhoudm1743/go-web/apps/{{.PackageName}}/models"
	"github.com/zhoudm1743/go-web/core/facades"
	"github.com/zhoudm1743/go-web/core/response"
	"strconv"
)

// {{.StructName}}Controller {{.Description}}控制器
type {{.StructName}}Controller struct{}

// New{{.StructName}}Controller 创建{{.Description}}控制器
func New{{.StructName}}Controller() *{{.StructName}}Controller {
	return &{{.StructName}}Controller{}
}
{{if .HasList}}
// Get{{.PluralName}} 获取{{.Description}}列表
func (c *{{.StructName}}Controller) Get{{.PluralName}}(ctx *gin.Context) {
	var params dto.{{.StructName}}QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	db := facades.DB()
	var total int64
	var items []*models.{{.StructName}}
	
	query := db.Model(&models.{{.StructName}}{})
	
	// 应用查询条件
	{{range .QueryFields}}
	if params.{{.FieldName}} != {{.ZeroValue}} {
		query = query.Where("{{.ColumnName}} = ?", params.{{.FieldName}})
	}
	{{end}}

	{{if .HasRelations}}
	// 应用关联表查询
	{{range .JoinFields}}
	if params.{{.RelatedFieldName}}Filter != "" {
		query = query.Joins("JOIN {{.JoinTable}} ON {{.JoinCondition}}").
			Where("{{.FilterCondition}} = ?", params.{{.RelatedFieldName}}Filter)
	}
	{{end}}
	
	// 应用预加载
	if params.WithRelations {
		query = (&models.{{.StructName}}{}).LoadRelations(query)
	} else {
		// 默认预加载关键关联
		{{range .PreloadFields}}
		query = query.Preload("{{.FieldName}}")
		{{end}}
	}
	{{end}}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 分页查询
	if err := query.Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&items).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	// 构造响应数据
	result := &dto.{{.StructName}}ListResponse{
		Total: total,
		List:  make([]*dto.{{.StructName}}Response, len(items)),
	}

	for i, item := range items {
		resp := &dto.{{.StructName}}Response{
			ID:        item.ID,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		
		// 复制其他字段
		response.Copy(resp, item)
		result.List[i] = resp
	}

	response.OkWithData(ctx, result)
}
{{end}}

{{if .HasDetail}}
// Get{{.StructName}} 获取{{.Description}}详情
func (c *{{.StructName}}Controller) Get{{.StructName}}(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	db := facades.DB()
	var item models.{{.StructName}}
	
	query := db{{range .PreloadFields}}.Preload("{{.FieldName}}"){{end}}
	{{if .HasRelations}}
	// 预加载关联数据
	withRelations := ctx.Query("withRelations")
	if withRelations == "true" {
		query = (&models.{{.StructName}}{}).LoadRelations(query)
	}
	{{end}}

	if err := query.First(&item, itemID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "{{.Description}}不存在")
		return
	}

	// 构造响应数据
	resp := &dto.{{.StructName}}Response{
		ID:        item.ID,
		CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// 复制其他字段
	response.Copy(resp, item)

	response.OkWithData(ctx, resp)
}
{{end}}

{{if .HasCreate}}
// Create{{.StructName}} 创建{{.Description}}
func (c *{{.StructName}}Controller) Create{{.StructName}}(ctx *gin.Context) {
	var req dto.{{.StructName}}CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	// 创建{{.Description}}
	item := &models.{{.StructName}}{}
	response.Copy(item, req)

	db := facades.DB()
	if err := db.Create(item).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "创建成功")
}
{{end}}

{{if .HasUpdate}}
// Update{{.StructName}} 更新{{.Description}}
func (c *{{.StructName}}Controller) Update{{.StructName}}(ctx *gin.Context) {
	var req dto.{{.StructName}}UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "请求参数有误")
		return
	}

	db := facades.DB()
	var item models.{{.StructName}}
	if err := db.First(&item, req.ID).Error; err != nil {
		response.FailWithMsg(ctx, response.Failed, "{{.Description}}不存在")
		return
	}

	// 只更新提供的字段
	updates := map[string]interface{}{}

	// 创建一个临时对象，用于复制非空字段
	temp{{.StructName}} := &models.{{.StructName}}{}
	response.Copy(temp{{.StructName}}, req)

	{{range .UpdateFields}}
	if req.{{.FieldName}} != {{.ZeroValue}} {
		updates["{{.ColumnName}}"] = temp{{$.StructName}}.{{.FieldName}}
	}
	{{end}}

	if err := db.Model(&item).Updates(updates).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "更新成功")
}
{{end}}

{{if .HasDelete}}
// Delete{{.StructName}} 删除{{.Description}}
func (c *{{.StructName}}Controller) Delete{{.StructName}}(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMsg(ctx, response.ParamsValidError, "ID无效")
		return
	}

	db := facades.DB()
	if err := db.Delete(&models.{{.StructName}}{}, itemID).Error; err != nil {
		response.Fail(ctx, response.SystemError)
		return
	}

	response.OkWithMsg(ctx, "删除成功")
}
{{end}}
`

	// 准备模板数据
	type FieldData struct {
		FieldName  string
		ColumnName string
		ZeroValue  string
	}

	type TemplateData struct {
		*Config
		PluralName    string
		QueryFields   []FieldData
		UpdateFields  []FieldData
		HasRelations  bool
		PreloadFields []struct {
			FieldName string
		}
		JoinFields []struct {
			RelatedFieldName string
			JoinTable        string
			JoinCondition    string
			FilterCondition  string
		}
	}

	data := TemplateData{
		Config:        g.Config,
		PluralName:    ToPlural(g.Config.StructName),
		QueryFields:   make([]FieldData, 0),
		UpdateFields:  make([]FieldData, 0),
		HasRelations:  false,
		PreloadFields: make([]struct{ FieldName string }, 0),
		JoinFields: make([]struct {
			RelatedFieldName string
			JoinTable        string
			JoinCondition    string
			FilterCondition  string
		}, 0),
	}

	// 处理字段
	for _, field := range g.Config.Fields {
		// 检查是否有关系字段
		if field.IsRelation {
			data.HasRelations = true

			// 添加到预加载字段列表
			// 仅当字段设置了预加载选项时才添加到预加载列表
			if field.Preload {
				// 预加载的字段名应该是关系字段名(FieldName)，不是外键名称
				data.PreloadFields = append(data.PreloadFields, struct{ FieldName string }{
					FieldName: field.FieldName,
				})
			}

			// 为BelongsTo和HasOne关系添加JOIN支持
			if (field.RelationType == BelongsTo || field.RelationType == HasOne) && field.Joinable {
				// 确定外键和引用字段
				foreignKey := field.ForeignKey
				references := field.References

				if foreignKey == "" {
					// 默认外键命名：关联模型名+ID
					foreignKey = field.RelatedModel + "ID"
				}

				if references == "" {
					// 默认引用字段：ID
					references = "ID"
				}

				// 获取表名
				mainTable := g.Config.TableName
				relatedTable := ToSnakeCase(field.RelatedModel) + "s" // 假设表名是模型名的复数形式

				// 生成JOIN条件
				joinCondition := fmt.Sprintf("%s.%s = %s.%s", mainTable, ToSnakeCase(foreignKey), relatedTable, ToSnakeCase(references))
				if field.JoinCondition != "" {
					joinCondition = field.JoinCondition
				}

				// 生成过滤条件
				filterCondition := fmt.Sprintf("%s.name", relatedTable) // 默认过滤字段为name
				if field.FilterCondition != "" {
					filterCondition = field.FilterCondition
				}

				// 添加JOIN查询
				data.JoinFields = append(data.JoinFields, struct {
					RelatedFieldName string
					JoinTable        string
					JoinCondition    string
					FilterCondition  string
				}{
					RelatedFieldName: field.FieldName,
					JoinTable:        relatedTable,
					JoinCondition:    joinCondition,
					FilterCondition:  filterCondition,
				})
			}
		}

		// 如果是主键字段，跳过
		if field.IsPrimaryKey {
			continue
		}

		// 确定字段的零值
		zeroValue := "0"
		if field.FieldType == "string" {
			zeroValue = `""`
		} else if field.FieldType == "bool" {
			zeroValue = "false"
		}

		// 查询参数字段
		if field.IsSearchable || field.IsFilterable {
			// 只对非关系字段添加查询条件
			if !field.IsRelation {
				queryField := FieldData{
					FieldName:  field.FieldName,
					ColumnName: field.ColumnName,
					ZeroValue:  zeroValue,
				}
				data.QueryFields = append(data.QueryFields, queryField)
			}
		}

		// 更新字段
		if !field.IsRelation {
			updateField := FieldData{
				FieldName:  field.FieldName,
				ColumnName: field.ColumnName,
				ZeroValue:  zeroValue,
			}
			data.UpdateFields = append(data.UpdateFields, updateField)
		}
	}

	// 解析模板
	t, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return fmt.Errorf("解析控制器模板失败: %w", err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染控制器模板失败: %w", err)
	}

	// 确保目录存在
	dir := filepath.Join(g.RootPath, "server/apps", g.Config.PackageName, "controllers")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	filename := filepath.Join(dir, strings.ToLower(g.Config.StructName)+"_controller.go")
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入控制器文件失败: %w", err)
	}

	// 记录生成的文件
	g.AddGeneratedFile(filename, "controller")

	fmt.Printf("生成控制器文件: %s\n", filename)
	return nil
}
