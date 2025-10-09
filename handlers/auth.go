package handlers

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/fred29910/xizexcample/enhanced_router"
	messages "github.com/fred29910/xizexcample/proto"
	"google.golang.org/protobuf/proto"
)

// HandleAuth handles user authentication
func HandleAuth(req ziface.IRequest) {
	// 1. Parse the protobuf data
	var authReq messages.AuthRequest
	if err := proto.Unmarshal(req.GetData(), &authReq); err != nil {
		fmt.Println("HandleAuth: unmarshal error: ", err)
		// Send error response
		resp := &messages.Response{Code: 1, Message: "Parse error"}
		data, _ := proto.Marshal(resp)
		req.GetConnection().SendMsg(1, data)
		return
	}

	// 2. Business logic: validate username and password
	if authReq.Username == "user" && authReq.Password == "pass" {
		// 3. Generate JWT
		tokenString, err := enhanced_router.GenerateJWT(authReq.Username)
		if err != nil {
			// Handle error
			return
		}

		// 4. Success response
		authResp := &messages.AuthResponse{Token: tokenString}
		respData, _ := proto.Marshal(authResp)
		resp := &messages.Response{Code: 0, Data: respData}
		data, _ := proto.Marshal(resp)
		req.GetConnection().SendMsg(1, data)
	} else {
		// 4. Failure response
		resp := &messages.Response{Code: 2, Message: "Authentication failed"}
		data, _ := proto.Marshal(resp)
		req.GetConnection().SendMsg(1, data)
	}
}
