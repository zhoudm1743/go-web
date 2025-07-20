package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// generateModel 生成模型文件
func (g *Generator) generateModel() error {
	// 模型模板
	const modelTemplate = `package models

import (
	"time"
	"github.com/zhoudm1743/go-web/core/facades"
	"gorm.io/gorm"
)

// {{.StructName}} {{.Description}}
type {{.StructName}} struct {
	ID        uint           ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `                // 主键ID
	CreatedAt time.Time      ` + "`" + `json:"createdAt"` + "`" + `                           // 创建时间
	UpdatedAt time.Time      ` + "`" + `json:"updatedAt"` + "`" + `                           // 更新时间
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `                      // 删除时间
{{range .Fields}}
	{{.FieldName}} {{.FieldType}} ` + "`" + `{{.GormTag}}` + "`" + `{{if .FieldDesc}} // {{.FieldDesc}}{{end}}
{{end}}
}

// TableName 设置表名
func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

{{if .HasRelations}}
// 关系预加载
func (m *{{.StructName}}) LoadRelations(db *gorm.DB) *gorm.DB {
	query := db
	{{range .PreloadFields}}
	query = query.Preload("{{.FieldName}}")
	{{end}}
	return query
}
{{end}}
`

	// 准备模板数据
	type FieldData struct {
		FieldName string
		FieldType string
		GormTag   string
		FieldDesc string
	}

	type TemplateData struct {
		*Config
		Fields        []FieldData
		HasRelations  bool
		PreloadFields []struct {
			FieldName string
		}
	}

	data := TemplateData{
		Config:        g.Config,
		Fields:        make([]FieldData, 0),
		HasRelations:  false,
		PreloadFields: make([]struct{ FieldName string }, 0),
	}

	// 准备字段数据
	for _, field := range g.Config.Fields {
		var gormTag string
		var fieldType string

		// 如果是关系字段
		if field.IsRelation {
			// 根据关系类型生成字段
			switch field.RelationType {
			case BelongsTo:
				// belongsTo关系: Post belongs to User
				// Post有一个UserID外键，引用User的ID
				fieldType = field.RelatedModel
				gormTag = fmt.Sprintf(`json:"%s"`, ToLowerCamel(field.FieldName))

				// 添加外键字段
				if field.ForeignKey != "" {
					// 使用自定义外键
					fkFieldName := field.ForeignKey
					data.Fields = append(data.Fields, FieldData{
						FieldName: fkFieldName,
						FieldType: "uint",
						GormTag:   fmt.Sprintf(`gorm:"column:%s" json:"%s"`, ToSnakeCase(fkFieldName), ToLowerCamel(fkFieldName)),
						FieldDesc: fmt.Sprintf("%s外键", field.RelatedModel),
					})
				} else {
					// 使用默认外键命名
					fkFieldName := field.FieldName + "ID"
					data.Fields = append(data.Fields, FieldData{
						FieldName: fkFieldName,
						FieldType: "uint",
						GormTag:   fmt.Sprintf(`gorm:"column:%s" json:"%s"`, ToSnakeCase(fkFieldName), ToLowerCamel(fkFieldName)),
						FieldDesc: fmt.Sprintf("%s外键", field.RelatedModel),
					})
				}

				// 添加关系配置
				if field.ForeignKey != "" && field.References != "" {
					gormTag = fmt.Sprintf(`gorm:"foreignKey:%s;references:%s" json:"%s"`, field.ForeignKey, field.References, ToLowerCamel(field.FieldName))
				} else if field.ForeignKey != "" {
					gormTag = fmt.Sprintf(`gorm:"foreignKey:%s" json:"%s"`, field.ForeignKey, ToLowerCamel(field.FieldName))
				}

			case HasOne:
				// hasOne关系: User has one Profile
				// Profile有一个UserID外键，引用User的ID
				fieldType = field.RelatedModel
				gormTag = fmt.Sprintf(`json:"%s"`, ToLowerCamel(field.FieldName))

				// 添加关系配置
				if field.ForeignKey != "" && field.References != "" {
					gormTag = fmt.Sprintf(`gorm:"foreignKey:%s;references:%s" json:"%s"`, field.ForeignKey, field.References, ToLowerCamel(field.FieldName))
				} else if field.ForeignKey != "" {
					gormTag = fmt.Sprintf(`gorm:"foreignKey:%s" json:"%s"`, field.ForeignKey, ToLowerCamel(field.FieldName))
				}

			case HasMany:
				// hasMany关系: User has many Post
				// Post有一个UserID外键，引用User的ID
				fieldType = fmt.Sprintf("[]%s", field.RelatedModel)
				gormTag = fmt.Sprintf(`json:"%s"`, ToLowerCamel(field.FieldName))

				// 添加关系配置
				if field.ForeignKey != "" && field.References != "" {
					gormTag = fmt.Sprintf(`gorm:"foreignKey:%s;references:%s" json:"%s"`, field.ForeignKey, field.References, ToLowerCamel(field.FieldName))
				} else if field.ForeignKey != "" {
					gormTag = fmt.Sprintf(`gorm:"foreignKey:%s" json:"%s"`, field.ForeignKey, ToLowerCamel(field.FieldName))
				}

			case ManyToMany:
				// manyToMany关系: User has many Role through UserRole
				fieldType = fmt.Sprintf("[]%s", field.RelatedModel)

				// 添加关系配置
				if field.JoinTable != "" {
					if field.ForeignKey != "" && field.References != "" {
						gormTag = fmt.Sprintf(`gorm:"many2many:%s;foreignKey:%s;references:%s" json:"%s"`, field.JoinTable, field.ForeignKey, field.References, ToLowerCamel(field.FieldName))
					} else if field.ForeignKey != "" {
						gormTag = fmt.Sprintf(`gorm:"many2many:%s;foreignKey:%s" json:"%s"`, field.JoinTable, field.ForeignKey, ToLowerCamel(field.FieldName))
					} else {
						gormTag = fmt.Sprintf(`gorm:"many2many:%s" json:"%s"`, field.JoinTable, ToLowerCamel(field.FieldName))
					}
				} else {
					gormTag = fmt.Sprintf(`json:"%s"`, ToLowerCamel(field.FieldName))
				}
			}

			// 添加到字段列表
			data.Fields = append(data.Fields, FieldData{
				FieldName: field.FieldName,
				FieldType: fieldType,
				GormTag:   gormTag,
				FieldDesc: field.FieldDesc,
			})

			// 记录这是一个关系字段
			data.HasRelations = true

			// 如果需要预加载
			if field.Preload {
				data.PreloadFields = append(data.PreloadFields, struct{ FieldName string }{
					FieldName: field.FieldName,
				})
			}
		} else {
			// 普通字段
			if field.IsPrimaryKey {
				gormTag = fmt.Sprintf(`gorm:"primaryKey;column:%s" json:"%s"`, field.ColumnName, ToLowerCamel(field.FieldName))
			} else {
				gormTag = fmt.Sprintf(`gorm:"column:%s" json:"%s"`, field.ColumnName, ToLowerCamel(field.FieldName))
			}

			data.Fields = append(data.Fields, FieldData{
				FieldName: field.FieldName,
				FieldType: field.FieldType,
				GormTag:   gormTag,
				FieldDesc: field.FieldDesc,
			})
		}
	}

	// 解析模板
	t, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("解析模型模板失败: %w", err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染模型模板失败: %w", err)
	}

	// 确保目录存在
	dir := filepath.Join(g.RootPath, "server/apps/admin/models")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	filename := filepath.Join(dir, ToLowerCamel(g.Config.StructName)+".go")
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入模型文件失败: %w", err)
	}

	// 记录生成的文件
	g.AddGeneratedFile(filename, "model")

	fmt.Printf("生成模型文件: %s\n", filename)
	return nil
}
