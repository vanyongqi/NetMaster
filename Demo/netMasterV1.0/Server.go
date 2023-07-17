package main

import (
	"fmt"
	"netMaster/Master/masiface"
	"netMaster/Master/masnet"
)

/**
0 鲁班七号 不得不承认，有时候肌肉比头脑管用。				  200
1 娜可露露 玛玛哈哈					 				  201
1 周瑜 	  这种前赴后继送死的勇气令我钦佩！ 				  202
2 孙策     你以为自己是谁?灾祸吗?命运吗?可终究会被我踩在脚下！  203
4 赵云     会安然无恙的,我保证。  						  204
*/

// ping test 自定义路由 继承基类路由 3个方法
type LubanRouter struct {
	masnet.BaseRouter
}

func (this *LubanRouter) Handle(request masiface.IRequest) {
	fmt.Println("接入鲁班数据")
	//先读取客户端语句，再回写
	fmt.Println("recv from client ;msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(200, []byte("鲁班上线状态：不得不承认，有时候肌肉比头脑管用。"))
	if err != nil {
		fmt.Println(err)
	}
}

// ping test 自定义路由 继承基类路由 3个方法
type ZhouYuRouter struct {
	masnet.BaseRouter
}

func (this *ZhouYuRouter) Handle(request masiface.IRequest) {
	fmt.Println("接入周瑜数据")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(201, []byte("这种前赴后继送死的勇气令我钦佩！"))
	if err != nil {
		fmt.Println(err)
	}
}

// ping test 自定义路由 继承基类路由 3个方法
type SunCeRouter struct {
	masnet.BaseRouter
}

func (this *SunCeRouter) Handle(request masiface.IRequest) {
	fmt.Println("接入孙策数据")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(202, []byte("你以为自己是谁?灾祸吗?命运吗?可终究会被我踩在脚下！"))
	if err != nil {
		fmt.Println(err)
	}
}

// ping test 自定义路由 继承基类路由 3个方法
type ZhaoYunRouter struct {
	masnet.BaseRouter
}

func (this *ZhaoYunRouter) Handle(request masiface.IRequest) {
	fmt.Println("接入赵云数据")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(203, []byte("会安然无恙的,我保证。 "))
	if err != nil {
		fmt.Println(err)
	}
}

// ping test 自定义路由 继承基类路由 3个方法
type NakeLuLuRouter struct {
	masnet.BaseRouter
}

func (this *NakeLuLuRouter) Handle(request masiface.IRequest) {
	fmt.Println("接入娜可露露数据")
	//先读取客户端语句，再回写ping... ping... ping...
	fmt.Println("recv from client msgID = ", request.GetMsgId(),
		"data = ", string(request.GetData()))
	//send 会打包 然后client解包
	err := request.GetConnection().SendMsg(204, []byte("玛玛哈哈"))
	if err != nil {
		fmt.Println(err)
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
// 链接开始之前需要执行的hook函数
func DoConnectionBegin(conn masiface.IConnection) {
	fmt.Println("=====Connetcin is Called 广播 ，XXX进入房间")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
	//链接断开前 需要执行的函数
	fmt.Println("Set Conn Property....")
	conn.SetProperty("Type", "排位")
	conn.SetProperty("Level", "青铜")
	conn.SetProperty("mallocSize", "10人房间")

}

// 链接断开之前需要执行的hook函数
func DoConnectionLost(conn masiface.IConnection) {
	fmt.Println("====> Connection Lost is Called... 广播 , xxx退出房间")
	fmt.Println("conn ID = ", conn.GetConnID(), "is Lost...")

	//获取链接属性
	if Type, err := conn.GetProperty("Type"); err == nil {
		fmt.Println("比赛类型", Type)
	}
	if Level, err := conn.GetProperty("Level"); err == nil {
		fmt.Println("Level", Level)
	}
	if mallocSize, err := conn.GetProperty("mallocSize"); err == nil {
		fmt.Println("mallocSize", mallocSize)
	}
}

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master V1.0")
	//1.1 注册链接hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//2 给当前框架添加自定义router
	s.AddRouter(0, &LubanRouter{})    //鲁班
	s.AddRouter(1, &NakeLuLuRouter{}) //娜可露露
	s.AddRouter(2, &ZhouYuRouter{})   //周瑜
	s.AddRouter(3, &SunCeRouter{})    //孙策
	s.AddRouter(4, &ZhaoYunRouter{})  //赵云
	//启动serve()
	s.Serve()
}
