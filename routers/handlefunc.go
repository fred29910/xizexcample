package routers

import (
	"context"
	"fmt"

	"github.com/aceld/zinx/ziface"
	"google.golang.org/protobuf/proto"
)

type HandleFunc[T any, R any] func(ctx context.Context, req T) (resp R, err error)

// PreHandle 实现 ziface.IRouter 的 PreHandle 方法（默认空实现）。
func (h HandleFunc[T, R]) PreHandle(request ziface.IRequest) {}

// Handle 实现 ziface.IRouter 的 Handle 方法，调用函数本体处理请求。
func (h HandleFunc[T, R]) Handle(request ziface.IRequest) {

	var req T

	if err := proto.Unmarshal(request.GetData(), &req); err != nil {
		fmt.Println("Handle: unmarshal error: ", err)
		return
	}

	resp, err := h(context.Background(), req)

	// 	proto.Message is not an interface, so we cannot use it as a type constraint
	// 	we need to use a type constraint that is an interface
	// 	we can use a type constraint that is an interface
	// 	we can use a type constraint that is an interface

	data, err := proto.Marshal(resp)
	if err != nil {
		fmt.Println("Handle: marshal error: ", err)
		return
	}
	request.GetConnection().SendMsg(request.GetMsgId(), data)

}

// PostHandle 实现 ziface.IRouter 的 PostHandle 方法（默认空实现）。
func (h HandleFunc[T, R]) PostHandle(request ziface.IRequest) {}
