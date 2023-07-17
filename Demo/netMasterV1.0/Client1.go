package main

import (
	"fmt"
	"io"
	"net"
	"netMaster/Master/masnet"
	"time"
)

// Client
func main() {
	fmt.Println("Client V0.8 start...")
	time.Sleep(1 * time.Second)
	//conn 连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start err,exit!")
	}
	for {
		//发送封包消息
		dp := masnet.NewDataPack()
		binaryMsg, _ := dp.Pack(masnet.NewMessage(1, []byte("玩家1号发送请求：娜可露露--->")))
		if err != nil {
			fmt.Println("Pack error", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器回复一个message msgId 1 pingping。。。
		//先读取流中的head部分，粘包处理  ID datalen

		binaryHead := make([]byte, dp.GetHeadLen()) //读头
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error ", err)
			break
		}

		msgHead, err := dp.UnPack(binaryHead) //解头
		if err != nil {
			fmt.Println("client unpack masgHead error", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//msg 里是有数据的
			msg := msgHead.(*masnet.Message) //根据头创消息缓冲
			msg.Data = make([]byte, msg.GetMsgLen())
			//读消息缓冲，接受到了回复的消息了
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("recv Server Msg: ID", msg.Id, ",len=", msg.DataLen,
				"data = ", string(msg.Data))
		}
		//		//读取data
		//_, err := conn.Write([]byte("hello,mt"))
		//if err != nil {
		//	fmt.Println("write conn err", err)
		//	return
		//}
		//
		//buf := make([]byte, 512)
		//cnt, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf err")
		//	return
		//}
		//fmt.Printf("server call back %s,cnt =%d \n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
	//调用 Write写数据kk

}
