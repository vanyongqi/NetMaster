package masnet

import (
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

		//阻塞 等待客户端连接，处理客户端连接业务
		for {
			//客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept arr", err)
				continue
			}
			//client 以及建立连接，可以开始进行业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf) //阻塞
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					//echo
					fmt.Printf("recv : %s,cnt is %d \n", buf, cnt)
					if _, err := conn.Write(buf[0:cnt]); err != nil {
						fmt.Println("write back buf err")
						continue
					}
				}
				fmt.Println()
			}()
		}
	}()

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
