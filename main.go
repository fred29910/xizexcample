package main

import (
	"fmt"
	"github.com/aceld/zinx/zinxserver"
)

func main() {
	// 创建一个Zinx服务器句柄
	s := zinxserver.NewServer("[Zinx Niuniu Server]")

	// 注册路由
	// TODO: 将在后续任务中实现 router.InitRouter(s)

	// 启动服务
	fmt.Println("Starting server...")
	s.Serve()
}