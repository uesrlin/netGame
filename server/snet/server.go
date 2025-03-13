package snet

import (
	"fmt"
	"net"
)

/**
 * @Description
 * @Date 2025/3/6 22:40
 **/

// Server 服务器配置
type Server struct {
	Addr    string // 服务器地址
	Port    int    // 服务器端口
	Name    string // 服务器名称
	Version string // 服务器版本
}

func (s *Server) Start() {
	// 这里写tcp服务器启动的逻辑
	var cid uint32 = 0
	//监听本地的端口
	listenAddr := fmt.Sprintf("%s:%d", s.Addr, s.Port)

	go func() {
		addr, er := net.ResolveTCPAddr("tcp", listenAddr)
		if er != nil {
			fmt.Println("解析本地地址失败", er.Error())
			return
		}
		listen, err := net.ListenTCP("tcp", addr)
		if err != nil {
			fmt.Println("启动监听失败", err.Error())
			return
		}
		fmt.Println("服务器启动成功....")
		// 记得关闭监听
		defer listen.Close()
		// 这里有多个连接，需要使用goroutine来处理

		for {
			// 建立连接
			conn, er0 := listen.AcceptTCP()
			if er0 != nil {
				fmt.Println("建立连接失败", er.Error())
				continue
			}
			dealConn := NewConnection(conn, cid)
			player := NewPlayer(dealConn)
			dealConn.PlayerRef = player
			cid++
			go dealConn.Start()

		}

	}()

}

func (s *Server) Stop() {
	// 这里写tcp服务器停止的逻辑

}

func (s *Server) Serve() {
	// 这里写tcp服务器运行的逻辑
	s.Start()
	// 阻塞在这里
	select {}
}

func NewServer(addr string, port int, name string, version string) *Server {
	if name == "" {
		name = "Server"
	}
	if version == "" {
		version = "1.0.0"
	}
	return &Server{Addr: addr, Port: port, Name: name, Version: version}
}
