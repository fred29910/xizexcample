package handlers

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	messages "github.com/fred29910/xizexcample/proto"
	"google.golang.org/protobuf/proto"
)

// HandlePush handles message pushing
func HandlePush(req ziface.IRequest) {
	// 1. Parse the protobuf data
	var pushReq messages.PushRequest
	if err := proto.Unmarshal(req.GetData(), &pushReq); err != nil {
		fmt.Println("HandlePush: unmarshal error: ", err)
		// Send error response
		resp := &messages.Response{Code: 1, Message: "Parse error"}
		data, _ := proto.Marshal(resp)
		req.GetConnection().SendMsg(2, data)
		return
	}

	// 2. Business logic: push message (e.g., to a channel or another user)
	fmt.Printf("Pushing message to %s: %s\n", pushReq.TargetUser, pushReq.Content)

	// 3. Success response
	pushResp := &messages.PushResponse{Success: true}
	respData, _ := proto.Marshal(pushResp)
	resp := &messages.Response{Code: 0, Data: respData}
	data, _ := proto.Marshal(resp)
	req.GetConnection().SendMsg(2, data)
}
