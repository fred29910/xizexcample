package main

import (
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/znet"
	"xizexcample/internal/conf"
	"xizexcample/internal/pkg/logger"
	"xizexcample/internal/router"
	"xizexcample/internal/server"
)

var (
	version = "dev"
	commit  = "none"
)

func init() {
	logger.InfoLogger.Printf("Version: %s, Commit: %s", version, commit)
}

func main() {
	// 在服务器启动前，通过 zconf.GlobalObject 配置全局设置
	zconf.GlobalObject.Host = conf.AppConfig.ServerHost
	zconf.GlobalObject.TCPPort = conf.AppConfig.ServerPort

	// 创建一个Zinx服务器句柄
	s := znet.NewServer()

	// 注册路由
	router.InitRouter(s)

	// 设置钩子
	s.SetOnConnStop(server.OnConnStop)

	// 启动服务
	logger.InfoLogger.Printf("Starting server at %s:%d...", zconf.GlobalObject.Host, zconf.GlobalObject.TCPPort)
	s.Serve()
	logger.InfoLogger.Println("Server stopped.")
}
