package logic

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"strings"
	"zinx-chatroom-demo/core"
)

// BaseRouter 空的基类，为了让具体的 Router 实现 Handle 方法
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (br *BaseRouter) Handle(req ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}

// TalkRouter 聊天广播路由
type TalkRouter struct {
	znet.BaseRouter
}

func (tr *TalkRouter) Handle(request ziface.IRequest) {
	// 1. 获取 Player 对象
	player, err := request.GetConnection().GetProperty("player")
	if err != nil {
		fmt.Println("Get player property error:", err)
		return
	}
	p := player.(*core.Player)

	// 2. 获取消息内容
	msg := string(request.GetData())

	// 3. 格式化消息并广播
	broadcastMsg := fmt.Sprintf("[%s]: %s", p.Name, msg)
	for _, player := range core.WorldMgrObj.GetAllPlayers() {
		if player.PID == p.PID {
			continue
		}
		player.SendMsg(2, []byte(broadcastMsg))
	}
}

// ChangeNameRouter 修改昵称路由
type ChangeNameRouter struct {
	znet.BaseRouter
}

func (cnr *ChangeNameRouter) Handle(request ziface.IRequest) {
	// 1. 获取 Player 对象
	player, err := request.GetConnection().GetProperty("player")
	if err != nil {
		fmt.Println("Get player property error:", err)
		return
	}
	p := player.(*core.Player)

	// 2. 获取新昵称
	// 消息格式: name|NewName
	parts := strings.Split(string(request.GetData()), "|")
	if len(parts) != 2 || parts[0] != "name" {
		p.SendMsg(1, []byte("Invalid command to change name. Use: name|YourNewName"))
		return
	}
	newName := parts[1]
	oldName := p.Name
	p.Name = newName

	// 3. 广播昵称变更通知
	broadcastMsg := fmt.Sprintf("--> User [%s] changed name to [%s].", oldName, newName)
	for _, player := range core.WorldMgrObj.GetAllPlayers() {
		player.SendMsg(1, []byte(broadcastMsg))
	}
}

// WhoRouter 获取在线列表路由
type WhoRouter struct {
	znet.BaseRouter
}

func (wr *WhoRouter) Handle(request ziface.IRequest) {
	// 1. 获取 Player 对象
	player, err := request.GetConnection().GetProperty("player")
	if err != nil {
		fmt.Println("Get player property error:", err)
		return
	}
	p := player.(*core.Player)

	// 2. 获取所有在线玩家的昵称
	var onlineUsers []string
	for _, player := range core.WorldMgrObj.GetAllPlayers() {
		onlineUsers = append(onlineUsers, player.Name)
	}

	// 3. 将列表发送给当前客户端
	whoMsg := "Online users: " + strings.Join(onlineUsers, ", ")
	p.SendMsg(1, []byte(whoMsg))
}