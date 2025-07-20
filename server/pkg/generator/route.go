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

// updateRoutesFile 更新路由文件，添加控制器变量
func updateRoutesFile(content string, structName string) string {
	lines := strings.Split(content, "\n")

	// 查找初始化控制器的位置
	controllerInitIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "// 初始化控制器") {
			controllerInitIndex = i
			break
		}
	}

	// 找到控制器初始化代码块的结束位置
	controllerBlockEndIndex := -1
	if controllerInitIndex >= 0 {
		for i := controllerInitIndex + 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "" {
				controllerBlockEndIndex = i
				break
			}
		}
	}

	// 如果找到了控制器初始化代码块
	if controllerInitIndex >= 0 && controllerBlockEndIndex >= 0 {
		// 添加新的控制器变量
		varName := strings.ToLower(structName) + "Controller"
		newControllerLine := fmt.Sprintf("\t%s := controllers.New%sController()", varName, structName)

		// 插入新行
		newLines := append(
			lines[:controllerBlockEndIndex],
			newControllerLine,
		)
		newLines = append(newLines, lines[controllerBlockEndIndex:]...)

		return strings.Join(newLines, "\n")
	}

	// 如果找不到控制器初始化代码块，返回原内容
	return content
}
