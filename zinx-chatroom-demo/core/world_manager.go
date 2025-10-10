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