package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// generateDTO 生成DTO文件
func (g *Generator) generateDTO() error {
	// DTO模板
	const dtoTemplate = `package dto

{{if .HasTimeImport}}
import "time"
{{end}}

// {{.StructName}}CreateRequest 创建{{.Description}}请求
type {{.StructName}}CreateRequest struct {
{{range .CreateFields}}
	{{.FieldName}} {{.FieldType}} ` + "`" + `json:"{{.JsonName}}" {{if .Binding}}binding:"{{.Binding}}"{{end}}` + "`" + `{{if .FieldDesc}} // {{.FieldDesc}}{{end}}
{{end}}
}

// {{.StructName}}UpdateRequest 更新{{.Description}}请求
type {{.StructName}}UpdateRequest struct {
	ID uint ` + "`" + `json:"id" binding:"required"` + "`" + ` // ID
{{range .UpdateFields}}
	{{.FieldName}} {{.FieldType}} ` + "`" + `json:"{{.JsonName}}" {{if .Binding}}binding:"{{.Binding}}"{{end}}` + "`" + `{{if .FieldDesc}} // {{.FieldDesc}}{{end}}
{{end}}
}

// {{.StructName}}QueryParams {{.Description}}查询参数
type {{.StructName}}QueryParams struct {
	Page     int ` + "`" + `form:"page"` + "`" + `      // 页码
	PageSize int ` + "`" + `form:"pageSize"` + "`" + `  // 每页条数
{{if .HasRelations}}
	WithRelations bool ` + "`" + `form:"withRelations"` + "`" + ` // 是否加载关联
{{end}}
{{range .QueryFields}}
	{{.FieldName}} {{.FieldType}} ` + "`" + `form:"{{.JsonName}}"` + "`" + `{{if .FieldDesc}} // {{.FieldDesc}}{{end}}
{{end}}
}

// {{.StructName}}Response {{.Description}}响应
type {{.StructName}}Response struct {
	ID        uint   ` + "`" + `json:"id"` + "`" + `         // ID
	CreatedAt string ` + "`" + `json:"createdAt"` + "`" + `  // 创建时间
	UpdatedAt string ` + "`" + `json:"updatedAt"` + "`" + `  // 更新时间
{{range .ResponseFields}}
	{{.FieldName}} {{.FieldType}} ` + "`" + `json:"{{.JsonName}}"` + "`" + `{{if .FieldDesc}} // {{.FieldDesc}}{{end}}
{{end}}
}

// {{.StructName}}ListResponse {{.Description}}列表响应
type {{.StructName}}ListResponse struct {
	Total int64                  ` + "`" + `json:"total"` + "`" + `  // 总数
	List  []*{{.StructName}}Response ` + "`" + `json:"list"` + "`" + `    // 列表
}
`

	// 准备模板数据
	type FieldData struct {
		FieldName string
		FieldType string
		JsonName  string
		Binding   string
		FieldDesc string
	}

	type TemplateData struct {
		*Config
		HasTimeImport  bool
		CreateFields   []FieldData
		UpdateFields   []FieldData
		QueryFields    []FieldData
		ResponseFields []FieldData
		HasRelations   bool
	}

	data := TemplateData{
		Config:         g.Config,
		HasTimeImport:  false,
		CreateFields:   make([]FieldData, 0),
		UpdateFields:   make([]FieldData, 0),
		QueryFields:    make([]FieldData, 0),
		ResponseFields: make([]FieldData, 0),
		HasRelations:   false,
	}

	// 处理字段
	for _, field := range g.Config.Fields {
		// 如果是主键字段，跳过
		if field.IsPrimaryKey {
			continue
		}

		jsonName := ToLowerCamel(field.FieldName)

		// 检查是否需要导入time包
		if strings.Contains(field.FieldType, "time.Time") {
			data.HasTimeImport = true
		}

		// 检查是否有关系字段
		if field.IsRelation {
			data.HasRelations = true

			// 创建响应字段，关系字段也需要添加到响应中
			var responseType string
			switch field.RelationType {
			case HasMany, ManyToMany:
				responseType = fmt.Sprintf("[]*%sResponse", field.RelatedModel)
			case BelongsTo, HasOne:
				responseType = fmt.Sprintf("*%sResponse", field.RelatedModel)
			}

			responseField := FieldData{
				FieldName: field.FieldName,
				FieldType: responseType,
				JsonName:  jsonName,
				FieldDesc: field.FieldDesc,
			}
			data.ResponseFields = append(data.ResponseFields, responseField)

			// 外键字段添加到创建和更新请求中
			if field.RelationType == BelongsTo && field.ForeignKey != "" {
				foreignKeyJsonName := ToLowerCamel(field.ForeignKey)

				// 创建请求字段
				binding := ""
				if field.Required {
					binding = "required"
				}

				createField := FieldData{
					FieldName: field.ForeignKey,
					FieldType: "uint",
					JsonName:  foreignKeyJsonName,
					Binding:   binding,
					FieldDesc: fmt.Sprintf("%s ID", field.FieldDesc),
				}
				data.CreateFields = append(data.CreateFields, createField)

				// 更新请求字段
				updateField := FieldData{
					FieldName: field.ForeignKey,
					FieldType: "uint",
					JsonName:  foreignKeyJsonName,
					FieldDesc: fmt.Sprintf("%s ID", field.FieldDesc),
				}
				data.UpdateFields = append(data.UpdateFields, updateField)

				// 查询参数字段
				if field.IsSearchable || field.IsFilterable {
					queryField := FieldData{
						FieldName: field.ForeignKey,
						FieldType: "uint",
						JsonName:  foreignKeyJsonName,
						FieldDesc: fmt.Sprintf("%s ID", field.FieldDesc),
					}
					data.QueryFields = append(data.QueryFields, queryField)
				}
			}
			continue
		}

		// 创建请求字段
		binding := ""
		if field.Required {
			binding = "required"
		}

		createField := FieldData{
			FieldName: field.FieldName,
			FieldType: field.FieldType,
			JsonName:  jsonName,
			Binding:   binding,
			FieldDesc: field.FieldDesc,
		}
		data.CreateFields = append(data.CreateFields, createField)

		// 更新请求字段
		updateField := FieldData{
			FieldName: field.FieldName,
			FieldType: field.FieldType,
			JsonName:  jsonName,
			FieldDesc: field.FieldDesc,
		}
		data.UpdateFields = append(data.UpdateFields, updateField)

		// 查询参数字段
		if field.IsSearchable || field.IsFilterable {
			queryField := FieldData{
				FieldName: field.FieldName,
				FieldType: field.FieldType,
				JsonName:  jsonName,
				FieldDesc: field.FieldDesc,
			}
			data.QueryFields = append(data.QueryFields, queryField)
		}

		// 响应字段
		responseField := FieldData{
			FieldName: field.FieldName,
			FieldType: field.FieldType,
			JsonName:  jsonName,
			FieldDesc: field.FieldDesc,
		}
		data.ResponseFields = append(data.ResponseFields, responseField)
	}

	// 解析模板
	t, err := template.New("dto").Parse(dtoTemplate)
	if err != nil {
		return fmt.Errorf("解析DTO模板失败: %w", err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染DTO模板失败: %w", err)
	}

	// 确保目录存在
	dir := filepath.Join(g.RootPath, "server/apps/admin/dto")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	filename := filepath.Join(dir, strings.ToLower(g.Config.StructName)+".go")
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入DTO文件失败: %w", err)
	}

	// 记录生成的文件
	g.AddGeneratedFile(filename, "dto")

	fmt.Printf("生成DTO文件: %s\n", filename)
	return nil
}
