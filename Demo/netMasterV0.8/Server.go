package main

import (
	"fmt"
	"netMaster/Master/masiface"
	"netMaster/Master/masnet"
)

// ping test 自定义路由 继承基类路由 3个方法
type PingRouter struct {
	masnet.BaseRouter
}

// Test Handle 重写一个
func (this *PingRouter) Handle(request masiface.IRequest) {
	fmt.Println("handle ping!!!!!!!!!")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client ;msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(200, []byte("ping----ping----ping 0 "))
	if err != nil {
		fmt.Println(err)
	}
}

// ping test 自定义路由 继承基类路由 3个方法
type HelloRouter struct {
	masnet.BaseRouter
}

// Test Handle 重写一个
func (this *HelloRouter) Handle(request masiface.IRequest) {
	fmt.Println("handle hello!!!!!!!")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(201, []byte("HELLO ---HELLO ----HELLO 1"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master V0.8")
	//2 给当前框架添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	//启动serve()
	s.Serve()
}
