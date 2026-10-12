package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

// 接口实现
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s,Port %d,is starting\n", s.IP, s.Port)

	go func() {

		//获取TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error :", err)
			return
		}
		//兼听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ,", s.Name, "succ,Listenning...")
		//阻塞等待用户的客户端链接，处理客户端链接业务（读写）
		for {
			//阻塞返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept err", err)
				continue
			}
			//做业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)
					//回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()
}
func (s *Server) Stop() {
	//TODO 将服务器的资源状态或一些链接信息进行停止或者回收
}
func (s *Server) Serve() {
	//启动服务功能
	s.Start()

	//TODO 做一些启动后的额外业务
	//阻塞
	select {}
}

/*
初始化server
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
