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

//Test PreHandle

func (this *PingRouter) PreHandle(request masiface.IRequest) {
	fmt.Println("call prehandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping 1..."))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

func (this *PingRouter) Handle(request masiface.IRequest) {
	fmt.Println("handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte(" ping 2... "))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// Test PostHandle

func (this *PingRouter) PostHandle(request masiface.IRequest) {
	fmt.Println("POST HANDLE")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping3 ..."))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master V0.3")
	//2 给当前框架添加自定义router
	s.AddRouter(&PingRouter{})
	//启动serve()
	s.Serve()
}
