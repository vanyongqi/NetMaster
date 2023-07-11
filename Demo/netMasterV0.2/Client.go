package main

import (
	"fmt"
	"net"
	"time"
)

// Client
func main() {
	fmt.Println("Client V0.2 start...")
	time.Sleep(1 * time.Second)
	//conn 连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start err,exit!")
	}
	for {
		_, err := conn.Write([]byte("hello,mt"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err")
			return
		}
		fmt.Printf("server call back %s,cnt =%d \n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
	//调用 Write写数据

}
