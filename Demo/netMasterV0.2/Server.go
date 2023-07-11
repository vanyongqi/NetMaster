package main

import "netMaster/Master/masnet"

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master V0.2")
	//启动serve()
	s.Serve()
}
