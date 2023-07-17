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

func DoConnectionBegin(conn masiface.IConnection) {
	fmt.Println("=====DoConnetcin")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
	//链接断开前 需要执行的函数
	fmt.Println("Set Conn Property....")
	conn.SetProperty("Name", "fanyongqi")
	conn.SetProperty("Home", "anhui")
	conn.SetProperty("github", "https://github.com/vanyongqi")

}

// 链接断开之前需要执行的函数
func DoConnectionLost(conn masiface.IConnection) {
	fmt.Println("====> Do ConnectionLost is Called...")
	fmt.Println("conn ID = ", conn.GetConnID(), "is Lost...")

	//获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("name", name)
	}
	if Home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home", Home)
	}
	if github, err := conn.GetProperty("github"); err == nil {
		fmt.Println("github", github)
	}
}

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master V0.8")
	//1.1 注册链接hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//2 给当前框架添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	//启动serve()
	s.Serve()
}
