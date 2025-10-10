# Zinx 聊天室项目代码实现规划

本文档包含了 Zinx 聊天室示例项目所有核心 Go 源代码的规划内容。

## 1. 项目入口: `main.go`

该文件是服务器的启动入口，负责初始化 Zinx 服务、注册全局钩子函数和业务路由。

```go
// file: zinx-chatroom-demo/main.go
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
```

## 2. 核心模块: `core/player.go`

该文件定义了 `Player` 对象，用于封装客户端连接的各种状态和业务。

```go
// file: zinx-chatroom-demo/core/player.go
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
```

## 3. 核心模块: `core/world_manager.go`

该文件定义了世界管理器，用于管理所有在线的 `Player` 对象。

```go
// file: zinx-chatroom-demo/core/world_manager.go
package core

import (
	"sync"
)

// WorldManager 当前世界总管理模块
type WorldManager struct {
	Players map[uint32]*Player // 在线玩家集合
	pLock   sync.RWMutex       // 保护 Players 的读写锁
}

// WorldMgrObj 提供一个对外的世界管理模块句柄 (单例)
var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[uint32]*Player),
	}
}

// AddPlayer 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	wm.Players[player.PID] = player
}

// RemovePlayer 删除一个玩家
func (wm *WorldManager) RemovePlayer(pid uint32) {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	delete(wm.Players, pid)
}

// GetPlayerByPID 通过玩家ID获取玩家对象
func (wm *WorldManager) GetPlayerByPID(pid uint32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

// GetAllPlayers 获取所有在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players := make([]*Player, 0)
	for _, p := range wm.Players {
		players = append(players, p)
	}
	return players
}
```
## 4. 业务逻辑模块: `logic/router.go`

该文件包含了所有业务逻辑的 Router 处理器。

```go
// file: zinx-chatroom-demo/logic/router.go
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
```
## 5. Go Modules 文件: `go.mod`

```mod
// file: zinx-chatroom-demo/go.mod
module zinx-chatroom-demo

go 1.18

require github.com/aceld/zinx v1.0.1
```