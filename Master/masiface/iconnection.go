package masiface

import "net"

//conn的抽象层

type IConnection interface {

	//启动连接 让当前的连接开始工作
	Start()
	//停止连接 结束当前连接的工作
	Stop()
	//获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前连接模块的id
	GetConnID() uint32
	//获取远程客户端的TCP状态
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程客户端
	Send(data []byte) error
}

// 定义一个处理连接业务的函数
type HandleFunc func(*net.TCPConn, []byte, int) error
