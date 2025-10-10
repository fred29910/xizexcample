package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"zinx-chatroom-demo/core"
	"zinx-chatroom-demo/logic"
)

// OnConnectionAdd 当客户端连接建立后的 Hook 函数
func OnConnectionAdd(conn ziface.IConnection) {
	fmt.Printf("ConnID=%d is Connected...\n", conn.GetConnID())

	// 创建一个 Player 对象
	player := core.NewPlayer(conn)

	// 上线广播
	player.SyncOnline()

	// 将 player 添加到在线管理器中
	core.WorldMgrObj.AddPlayer(player)

	// 将 conn 和 player 绑定
	conn.SetProperty("player", player)
}

// OnConnectionLost 当客户端连接断开时的 Hook 函数
func OnConnectionLost(conn ziface.IConnection) {
	fmt.Printf("ConnID=%d is Lost...\n", conn.GetConnID())

	// 获取 Player 对象
	player, err := conn.GetProperty("player")
	if err != nil {
		fmt.Println("Get player property error:", err)
		return
	}

	// 触发下线业务
	player.(*core.Player).SyncOffline()

	// 从在线管理器中删除
	core.WorldMgrObj.RemovePlayer(player.(*core.Player).PID)
}


func main() {
	// 1. 创建一个 Zinx Server
	s := znet.NewServer()

	// 2. 设置连接建立和断开的 Hook 函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 3. 注册业务路由
	s.AddRouter(2, &logic.TalkRouter{})      // 聊天消息
	s.AddRouter(3, &logic.ChangeNameRouter{}) // 修改昵称
	s.AddRouter(4, &logic.WhoRouter{})        // 获取在线列表

	// 4. 启动服务
	fmt.Println("Zinx Chatroom Server is running...")
	s.Serve()
}