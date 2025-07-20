package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zhoudm1743/go-web/bootstrap"
)

func main() {
	// 解析命令行参数
	mode := flag.String("mode", "http", "运行模式: http, cli")
	flag.Parse()

	// 初始化应用
	app, err := bootstrap.InitializeApp()
	if err != nil {
		fmt.Printf("初始化应用失败: %v\n", err)
		os.Exit(1)
	}

	// 设置应用模式
	app.SetMode(*mode)

	// 运行应用
	if err := app.Run(); err != nil {
		fmt.Printf("运行应用失败: %v\n", err)
		os.Exit(1)
	}
}
