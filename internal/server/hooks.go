package server

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/server"
)

// OnConnStop is the handler for when a connection is closed
func OnConnStop(conn ziface.IConnection) {
	fmt.Println("OnConnStop: conn stopped")

	// Get player ID from connection property
	playerID, err := conn.GetProperty("playerID")
	if err != nil {
		fmt.Println("OnConnStop: playerID not set on connection")
		return
	}

	// Find the room the player was in
	room := server.RoomMgr.GetRoomByPlayerID(playerID.(int64))
	if room == nil {
		fmt.Printf("OnConnStop: player %d was not in a room\n", playerID)
		return
	}

	// Mark the player as offline in the room
	room.SetPlayerOffline(playerID.(int64))
}