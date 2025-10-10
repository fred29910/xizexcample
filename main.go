package main

import (
	"fmt"
	"github.com/aceld/zinx/znet"
	"xizexcample/internal/server"
)

var (
	version = "dev"
	commit  = "none"
)

func init() {
	fmt.Printf("Version: %s, Commit: %s\n", version, commit)
}

func main() {
	// 创建一个Zinx服务器句柄
	s := znet.NewServer()

	// 注册路由
	// TODO: 将在后续任务中实现 router.InitRouter(s)

	// 设置钩子
	s.SetOnConnStop(server.OnConnStop)

	// 启动服务
	fmt.Println("Starting server...")
	s.Serve()
}
