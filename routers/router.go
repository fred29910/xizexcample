package routers

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/fred29910/xizexcample/handlers"
)

// AuthRouter handles authentication requests
type AuthRouter struct {
	znet.BaseRouter
}

// Handle processes the authentication request
func (ar *AuthRouter) Handle(request ziface.IRequest) {
	handlers.HandleAuth(request)
}

// PushRouter handles message push requests
type PushRouter struct {
	znet.BaseRouter
}

// Handle processes the message push request
func (pr *PushRouter) Handle(request ziface.IRequest) {
	handlers.HandlePush(request)
}
