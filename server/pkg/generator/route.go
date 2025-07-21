package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// generateRoute 生成路由文件
func (g *Generator) generateRoute() error {
	// 判断是否需要更新现有路由文件
	routesFilePath := filepath.Join(g.RootPath, "server/apps/admin/routes/routes.go")

	fileContent, err := os.ReadFile(routesFilePath)
	if err != nil {
		return fmt.Errorf("读取路由文件失败: %w", err)
	}

	// 检查是否需要引入控制器
	controllerImport := fmt.Sprintf(`controllers.New%sController()`, g.Config.StructName)
	if !strings.Contains(string(fileContent), controllerImport) {
		// 需要添加控制器变量定义
		newRouteContent := updateRoutesFile(string(fileContent), g.Config.StructName)

		// 写回文件
		if err := os.WriteFile(routesFilePath, []byte(newRouteContent), 0644); err != nil {
			return fmt.Errorf("更新路由文件失败: %w", err)
		}
		fmt.Printf("更新路由文件: %s\n", routesFilePath)

		// 记录修改的文件
		g.AddGeneratedFile(routesFilePath, "route")
	}

	// 生成路由组代码片段
	routeGroupTemplate := `
	// {{.Description}}路由
	{{.VarName}}Group := {{.RouterGroup}}.Group("/{{.ApiPrefix}}")
	{{if .NeedAuth}}{{.VarName}}Group.Use(middlewares.PermissionAuth()){{end}}
	{
		{{if .HasList}}{{.VarName}}Group.GET("/list", {{.VarName}}Controller.Get{{.PluralName}}){{end}}
		{{if .HasDetail}}{{.VarName}}Group.GET("/detail/:id", {{.VarName}}Controller.Get{{.StructName}}){{end}}
		{{if .HasCreate}}{{.VarName}}Group.POST("/create", {{.VarName}}Controller.Create{{.StructName}}){{end}}
		{{if .HasUpdate}}{{.VarName}}Group.PUT("/update", {{.VarName}}Controller.Update{{.StructName}}){{end}}
		{{if .HasDelete}}{{.VarName}}Group.DELETE("/delete/:id", {{.VarName}}Controller.Delete{{.StructName}}){{end}}
	}`

	// 准备模板数据
	type TemplateData struct {
		*Config
		VarName    string
		PluralName string
		NeedAuth   bool
	}

	data := TemplateData{
		Config:     g.Config,
		VarName:    strings.ToLower(g.Config.StructName),
		PluralName: ToPlural(g.Config.StructName),
		NeedAuth:   true, // 默认需要认证
	}

	// 生成路由组代码片段
	t, err := template.New("routeGroup").Parse(routeGroupTemplate)
	if err != nil {
		return fmt.Errorf("解析路由组模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染路由组模板失败: %w", err)
	}

	routeGroup := buf.String()
	fmt.Println("生成路由组代码片段成功，请将以下代码添加到对应的路由分组中:")
	fmt.Println("=============================================")
	fmt.Println(routeGroup)
	fmt.Println("=============================================")

	// 创建路由代码片段文件，方便后续使用
	snippetDir := filepath.Join(g.RootPath, "server/temp/snippets")
	if err := os.MkdirAll(snippetDir, 0755); err != nil {
		return fmt.Errorf("创建代码片段目录失败: %w", err)
	}

	snippetFilePath := filepath.Join(snippetDir, strings.ToLower(g.Config.StructName)+"_route_snippet.txt")
	if err := os.WriteFile(snippetFilePath, []byte(routeGroup), 0644); err != nil {
		return fmt.Errorf("写入路由代码片段文件失败: %w", err)
	}

	// 记录生成的代码片段文件
	g.AddGeneratedFile(snippetFilePath, "route_snippet")

	return nil
}

// updateRoutesFile 更新路由文件，添加控制器变量和路由注册
func updateRoutesFile(content string, structName string) string {
	lines := strings.Split(content, "\n")
	varName := strings.ToLower(structName) + "Controller"

	// 检查文件中是否已存在控制器变量定义
	controllerVarPattern := fmt.Sprintf("%s := controllers.New%sController()", varName, structName)
	if strings.Contains(content, controllerVarPattern) {
		// 已存在控制器变量，检查是否存在重复定义
		count := strings.Count(content, controllerVarPattern)
		if count > 1 {
			// 存在重复定义，移除所有实例然后只添加一次
			newContent := strings.ReplaceAll(content, controllerVarPattern, "")
			content = newContent
			lines = strings.Split(content, "\n")
		} else {
			// 只有一个实例，无需添加
		}
	}

	// 1. 添加控制器变量（如果不存在）
	controllerInitIndex := -1
	// 查找合适的插入位置，优先找控制器初始化注释
	for i, line := range lines {
		if strings.Contains(line, "// 初始化控制器") || strings.Contains(line, "Controller :=") {
			controllerInitIndex = i
			break
		}
	}

	// 如果没找到明确的控制器初始化区域，找其他控制器定义
	if controllerInitIndex == -1 {
		for i, line := range lines {
			if strings.Contains(line, "Controller :=") {
				controllerInitIndex = i
				break
			}
		}
	}

	// 如果依然没找到，尝试找一个合理位置
	if controllerInitIndex == -1 {
		for i, line := range lines {
			if strings.Contains(line, "func main()") || strings.Contains(line, "func RegisterRoutes") {
				controllerInitIndex = i + 1
				break
			}
		}
	}

	// 如果找到了插入位置，添加控制器变量
	if controllerInitIndex >= 0 && !strings.Contains(content, controllerVarPattern) {
		newLines := append(
			lines[:controllerInitIndex+1],
			fmt.Sprintf("\t%s := controllers.New%sController()", varName, structName),
		)
		newLines = append(newLines, lines[controllerInitIndex+1:]...)
		lines = newLines
	}

	// 2. 添加路由注册
	// 寻找私有路由分组的位置
	privateRoutesIndex := -1
	privateRoutesBlockStartIndex := -1
	privateRoutesBlockEndIndex := -1

	for i, line := range lines {
		if strings.Contains(line, "privateRoutes := r.Group") {
			privateRoutesIndex = i
		}
		// 找到路由分组的开始位置
		if privateRoutesIndex != -1 && strings.Contains(line, "{") {
			privateRoutesBlockStartIndex = i
		}
		// 找到路由分组的结束位置
		if privateRoutesBlockStartIndex != -1 && strings.TrimSpace(line) == "}" {
			privateRoutesBlockEndIndex = i
			break
		}
	}

	// 生成路由注册代码
	pluralName := structName
	if strings.HasSuffix(pluralName, "y") {
		pluralName = pluralName[:len(pluralName)-1] + "ies"
	} else {
		pluralName = pluralName + "s"
	}

	// 检查是否已经存在该结构体的路由
	routeMarker := fmt.Sprintf("// %s路由", structName)
	routeExists := false

	// 重新连接为完整字符串以便进行检查
	contentString := strings.Join(lines, "\n")
	if strings.Contains(contentString, routeMarker) {
		routeExists = true
	}

	// 插入路由代码到私有路由分组
	if privateRoutesBlockStartIndex >= 0 && privateRoutesBlockEndIndex >= 0 && !routeExists {
		// 构建完整的路由代码块
		routeCode := fmt.Sprintf("\n\t\t// %s路由", structName)

		// 使用一致的命名规范
		varName := ToLowerCamel(structName) + "Controller" // 确保控制器变量名使用小驼峰

		routeCode += fmt.Sprintf("\n\t\tprivateRoutes.GET(\"/%ss\", %s.Get%s)", strings.ToLower(structName), varName, pluralName)
		routeCode += fmt.Sprintf("\n\t\tprivateRoutes.GET(\"/%s/:id\", %s.Get%s)", strings.ToLower(structName), varName, structName)
		routeCode += fmt.Sprintf("\n\t\tprivateRoutes.POST(\"/%s\", %s.Create%s)", strings.ToLower(structName), varName, structName)
		routeCode += fmt.Sprintf("\n\t\tprivateRoutes.PUT(\"/%s\", %s.Update%s)", strings.ToLower(structName), varName, structName)
		routeCode += fmt.Sprintf("\n\t\tprivateRoutes.DELETE(\"/%s/:id\", %s.Delete%s)", strings.ToLower(structName), varName, structName)

		// 在私有路由组的结尾大括号前插入代码
		line := lines[privateRoutesBlockEndIndex]
		lines[privateRoutesBlockEndIndex] = routeCode
		newLines := append(lines[:privateRoutesBlockEndIndex+1], line)
		lines = newLines
	}

	return strings.Join(lines, "\n")
}
