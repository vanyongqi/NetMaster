package main

import "netMaster/Master/masnet"

func main() {
	//1 创建serve服务句柄，使用netmatser api
	s := masnet.NewServer("master arrived, TROUBLE  DISAPPEAR!")
	//启动serve()
	s.Serve()
}
