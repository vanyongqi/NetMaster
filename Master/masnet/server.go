package masnet

import (
	"fmt"
	"net"
	"netMaster/Master/masiface"
	"netMaster/Master/utils"
)

// 实现iserver接口
type Server struct {
	Name      string
	IPversion string
	IP        string
	Port      int
	//Router    masiface.IRouter
	//当前server消息管理模块，用来绑定MsgID和对应处理业务API关系
	MsgHandle masiface.IMsgHandle
	ConnMgr   masiface.IConnManager

	OnConnStart func(conn masiface.IConnection)
	OnConnStop  func(conn masiface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[netMaster] Server Name: %s,listener at IP :%s,Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[netMaster] Version %s,MaxConn:%d,MaxPacketSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	//1 get  addr
	go func() {

		//开启工作池以及配套的消息队列
		s.MsgHandle.StartWorkerPool()

		//net 参数为 "ip", "ip4" 或者为"ip6"，net 为空这默认 ip
		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("error happened: resovle tcp add error", err)
			return
		}

		//2 listen addr
		listener, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("error happened:resovle listenner failed: ", s.IPversion, err)
			return
		}
		fmt.Println("start Master successful", s.Name, "Listening :", s.Port)
		var cid uint32
		cid = 0
		// 3阻塞 等待客户端连接，处理客户端连接业务
		for {
			//客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept arr", err)
				continue
			}

			//设置最大链接个数的判断，如果超过当前最大的链接，那么关闭此新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端相应一个超出最大链接的错误包
				fmt.Println("=================>>>  connection is overload")
				conn.Close()
				continue
			}
			//将处理新连接的业务方法 和 conn 进行绑定 得到我们的连接模块
			//dealConn := NewConnection(conn, cid, CallBackToClient)
			//死函数CallBackToClient 替换为路由属性
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++

			//启动当前连接业务处理
			go dealConn.Start()
		}
	}()

}

// // 定义当前客户端连接所绑定的handle api，写死的，以后修改为路由
//
//	func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//		fmt.Println("[Conn Handlel CallbackToClient]...")
//		if _, err := conn.Write(data[:cnt]); err != nil {
//			fmt.Println("write back buf err", err)
//			return errors.New("CallBackToClient error")
//		}
//		return nil
//	}
func (s *Server) Stop() {
	//TODO 释放服务器资源、状态，进行停止和回收
	fmt.Println("[STOP] netMster server name ,", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start() //只进行箭监听
	//TODO 做一些启动服务器之后的额外业务
	//阻塞，以免start完成后结束server
	select {}
}

// 添加一个路由方法
func (s *Server) AddRouter(msgID uint32, router masiface.IRouter) {
	//TODO implement me
	s.MsgHandle.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

// 初始化链接模块方法
func NewServer(name string) masiface.Iserver {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPversion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		//Router:    nil,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}

	return s
}

func (s *Server) SetOnConnStart(hookfunc func(connection masiface.IConnection)) {
	s.OnConnStart = hookfunc
}

func (s *Server) SetOnConnStop(hookfunc func(connection masiface.IConnection)) {
	s.OnConnStop = hookfunc
}
func (s *Server) CallOnConnStart(conn masiface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------------> call on start()")
	}
	s.OnConnStart(conn)
}
func (s *Server) CallOnStop(conn masiface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-------> call on Stop()")
		s.OnConnStop(conn)
	}
}
