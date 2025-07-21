package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// generateModel 生成模型文件
func (g *Generator) generateModel() error {
	// 模型模板
	const modelTemplate = `package models

import (
	"time"
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
// LoadRelations 关系预加载
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
				// belongsTo关系: Post belongs to Category
				// Post有一个CategoryID外键，引用Category的ID
				// 构建字段名和标签

				// 获取关联模型字段和外键字段名
				relationFieldName := field.FieldName // 关联模型字段名，例如"Category"
				foreignKeyFieldName := ""

				if field.ForeignKey != "" {
					foreignKeyFieldName = field.ForeignKey // 自定义外键字段名
				} else {
					foreignKeyFieldName = relationFieldName + "ID" // 默认外键字段名，如"CategoryID"
				}

				// 生成外键列名，特殊处理ID结尾的情况
				var foreignKeyColumnName string

				// 检查是否以ID结尾，特殊处理
				if strings.HasSuffix(foreignKeyFieldName, "ID") && foreignKeyFieldName != "ID" {
					// 移除ID后缀
					base := foreignKeyFieldName[:len(foreignKeyFieldName)-2]

					// 处理基础部分转为蛇形命名
					for i, c := range base {
						if i > 0 && unicode.IsUpper(c) {
							foreignKeyColumnName += "_"
						}
						foreignKeyColumnName += string(unicode.ToLower(c))
					}

					// 添加_id后缀
					foreignKeyColumnName += "_id"
				} else if foreignKeyFieldName == "ID" {
					// ID直接转换为id
					foreignKeyColumnName = "id"
				} else {
					// 常规处理
					for i, c := range foreignKeyFieldName {
						if i > 0 && unicode.IsUpper(c) {
							foreignKeyColumnName += "_"
						}
						foreignKeyColumnName += string(unicode.ToLower(c))
					}
				}

				// 关联模型字段
				relationFieldType := "*" + field.RelatedModel // 使用指针类型
				relationFieldGormTag := ""

				// 添加关系配置
				if field.References != "" {
					relationFieldGormTag = fmt.Sprintf(`gorm:"foreignKey:%s;references:%s" json:"%s"`,
						foreignKeyFieldName, field.References, ToLowerCamel(relationFieldName))
				} else {
					relationFieldGormTag = fmt.Sprintf(`json:"%s"`, ToLowerCamel(relationFieldName))
				}

				// 添加关联模型字段
				data.Fields = append(data.Fields, FieldData{
					FieldName: relationFieldName, // 这里是关系字段名，例如"Category"
					FieldType: relationFieldType, // 类型是指针类型，例如"*Category"
					GormTag:   relationFieldGormTag,
					FieldDesc: field.FieldDesc,
				})

				// 添加外键字段
				data.Fields = append(data.Fields, FieldData{
					FieldName: foreignKeyFieldName, // 这里是外键字段名，例如"CategoryID"
					FieldType: "uint",              // 使用uint类型，应与被引用字段类型一致(通常是主键ID)
					GormTag:   fmt.Sprintf(`gorm:"column:%s" json:"%s"`, foreignKeyColumnName, ToLowerCamel(foreignKeyFieldName)),
					FieldDesc: fmt.Sprintf("%s外键", field.RelatedModel),
				})

				// 记录这是一个关系字段
				data.HasRelations = true

				// 如果需要预加载
				if field.Preload {
					data.PreloadFields = append(data.PreloadFields, struct{ FieldName string }{
						FieldName: relationFieldName, // 预加载字段名是关系字段名，不是外键名
					})
				}

				// 已经处理了belongsTo关系，直接跳到下一个字段
				continue

			case HasOne:
				// hasOne关系: User has one Profile
				// Profile有一个UserID外键，引用User的ID
				fieldType = "*" + field.RelatedModel // 使用指针类型
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
				// 跳过名为ID的主键字段，因为已经在模板中默认定义了
				if field.FieldName == "ID" {
					continue
				}
				// 确保column名称不为空
				columnName := field.ColumnName
				if columnName == "" {
					columnName = ToSnakeCase(field.FieldName)
				}
				gormTag = fmt.Sprintf(`gorm:"primaryKey;column:%s" json:"%s"`, columnName, ToLowerCamel(field.FieldName))
			} else {
				// 确保column名称不为空
				columnName := field.ColumnName
				if columnName == "" {
					columnName = ToSnakeCase(field.FieldName)
				}
				gormTag = fmt.Sprintf(`gorm:"column:%s" json:"%s"`, columnName, ToLowerCamel(field.FieldName))
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
	appDir := filepath.Join(g.RootPath, "server/apps", g.Config.PackageName)
	dir := filepath.Join(appDir, "models")
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
