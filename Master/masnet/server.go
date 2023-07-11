package masnet

import (
	"errors"
	"fmt"
	"net"
	"netMaster/Master/masiface"
)

// 实现iserver接口
type Server struct {
	Name      string
	IPversion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[start] Server Listener at IP :%s, Port:%d,is starting\n", s.IP, s.Port)
	//1 get  addr
	go func() {
		//"tcp", "tcp4", "tcp6" 默认为tcp
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
				fmt.Println("Accept error: ", err)
				continue
			}
			//将监听到的新连接的业务方法 与 conn 进行绑定 得到我们的连接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			//异步处理启动当前连接业务处理
			go dealConn.Start()
		}
	}()

}

// 定义当前客户端连接所绑定的handle api，写死的，以后修改为路由
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	fmt.Println("[Conn has handled a Callback To Client]...")
	return nil
}
func (s *Server) Stop() {
	//TODO 释放服务器资源、状态，进行停止和回收

}

func (s *Server) Serve() {
	s.Start() //只进行箭监听
	//TODO 做一些启动服务器之后的额外业务
	//阻塞，以免start完成后结束server
	select {}
}

func NewServer(name string) masiface.Iserver {
	s := &Server{
		Name:      name,
		IPversion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
	}
	return s
}
