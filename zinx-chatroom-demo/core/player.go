package core

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"sync"
)

// Player 玩家对象
type Player struct {
	PID  uint32          // 玩家ID (就是 ConnID)
	Conn ziface.IConnection // 当前玩家的连接
	Name string          // 玩家昵称
}

var PIDGen uint32 = 1 // 用于生成玩家ID的计数器
var IDLock sync.Mutex // 保护 PIDGen 的互斥锁

// NewPlayer 创建一个玩家对象
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个唯一的ID
	IDLock.Lock()
	defer IDLock.Unlock()
	pid := PIDGen
	PIDGen++

	return &Player{
		PID:  pid,
		Conn: conn,
		Name: fmt.Sprintf("user-%d", pid), // 默认昵称
	}
}

// SyncOnline 广播上线消息
func (p *Player) SyncOnline() {
	// 1. 获取所有在线玩家
	players := WorldMgrObj.GetAllPlayers()

	// 2. 遍历所有玩家，发送上线通知
	for _, player := range players {
		// 排除自己
		if player.PID == p.PID {
			continue
		}
		// 发送消息 (MsgID=1, 预留为系统消息)
		onlineMsg := fmt.Sprintf("--> User [%s] is online.", p.Name)
		player.SendMsg(1, []byte(onlineMsg))
	}
}

// SyncOffline 广播下线消息
func (p *Player) SyncOffline() {
	// 1. 获取所有在线玩家
	players := WorldMgrObj.GetAllPlayers()

	// 2. 遍历所有玩家，发送下线通知
	for _, player := range players {
		// 排除自己
		if player.PID == p.PID {
			continue
		}
		// 发送消息 (MsgID=1, 预留为系统消息)
		offlineMsg := fmt.Sprintf("--> User [%s] is offline.", p.Name)
		player.SendMsg(1, []byte(offlineMsg))
	}
}

// SendMsg 发送消息给客户端
func (p *Player) SendMsg(msgID uint32, data []byte) {
	if p.Conn == nil {
		fmt.Println("Connection in player is nil")
		return
	}
	if err := p.Conn.SendMsg(msgID, data); err != nil {
		fmt.Println("Player SendMsg error !", err)
		return
	}
}