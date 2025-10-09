package main

import (
	"github.com/aceld/zinx/znet"
	"github.com/fred29910/xizexcample/routers"
)

func main() {
	// 1. Create a server handle
	s := znet.NewServer()

	// 2. Enable TLS
	s.SetServerCert("server.crt", "server.key")

	// 3. Add custom routers
	s.AddRouter(1, &routers.AuthRouter{})
	s.AddRouter(2, &routers.PushRouter{})

	// 3. Start the server
	s.Serve()
}
