package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

// InitRouter 初始化路由
func InitRouter(server ziface.IServer) {
	// 注册连接创建和销毁的Hook
	server.SetOnConnStart(func(conn ziface.IConnection) {
		// TODO: 实现连接创建时的逻辑，例如将连接与玩家关联
	})

	server.SetOnConnStop(func(conn ziface.IConnection) {
		// TODO: 实现连接断开时的逻辑，例如从房间移除玩家
	})

	// 注册消息路由
	// TODO: 将在后续任务中为每个消息ID添加具体的处理函数
	server.AddRouter(101, &JoinRoomHandler{})      // C2S_JOIN_ROOM_REQ
	server.AddRouter(102, &PlayerReadyHandler{})   // C2S_PLAYER_READY_REQ
	server.AddRouter(103, &BidBankerHandler{})     // C2S_BID_BANKER_REQ
	server.AddRouter(104, &PlaceBetHandler{})      // C2S_PLACE_BET_REQ
	server.AddRouter(105, &ShowdownHandler{})      // C2S_SHOWDOWN_REQ
	server.AddRouter(106, &LeaveRoomHandler{})     // C2S_LEAVE_ROOM_REQ
}

// BaseRouter 基础路由器
type BaseRouter struct {
	znet.BaseRouter
}

// LeaveRoomHandler 离开房间处理器
type LeaveRoomHandler struct {
	BaseRouter
}
