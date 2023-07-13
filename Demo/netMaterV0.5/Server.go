package main

import (
	"fmt"
	"netMaster/Master/masiface"
	"netMaster/Master/masnet"
)

// ping test 自定义路由
type PingRouter struct {
	masnet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request masiface.IRequest) {
	fmt.Println("handle")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client ;msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping----ping----ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master V0.5")
	//2 给当前框架添加自定义router
	s.AddRouter(&PingRouter{})
	//启动serve()
	s.Serve()
}
