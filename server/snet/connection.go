package snet

import (
	"fmt"
	"net"
	"net_game/server/siface"
)

/**
 * @Description
 * @Date 2025/3/6 23:34
 **/

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	// 当前连接的状态
	isClosed bool
	// 处理该连接的方法Router
	//handleAPI siface.HandleFunc
	// 告知当前连接已经退出的channel
	ExitChan chan bool
	// 新增写通道
	msgChan chan []byte
	// 处理该链接的方法router
	Router siface.IRouter
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.ConnID)
	// 启动从当前连接的读数据的业务

	// 启动读协程
	go c.startReader()

	// 启动写协程
	go c.startWriter()

}

func (c *Connection) Stop() {
	fmt.Println("conn Stop()...ConnID=", c.ConnID)
	if c.isClosed == true {

		return
	}
	c.isClosed = true
	// 关闭socket连接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true
	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn

}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID

}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()

}

func (c *Connection) SendMsg(data []byte) error {
	// 将数据写入通道，由写协程处理
	c.msgChan <- data
	return nil
}

func NewConnection(conn *net.TCPConn, connID uint32, router siface.IRouter) *Connection {
	c := &Connection{
		Conn:   conn,
		ConnID: connID,
		//handleAPI: handleAPI,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		msgChan:  make(chan []byte), // 新增写通道初始化
		Router:   router,
	}
	return c
}

// 新增读协程实现
func (c *Connection) startReader() {

	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		n, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read error", err)
			return
		}
		req := Request{
			conn: c,
			// 注意这里的n   如果不采用n的话  会出现一定的乱码   因为buf的长度是512  而n是实际读取的长度
			data: buf[:n],
		}
		// 执行路由处理方法
		go func(request siface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// 新增写协程实现
func (c *Connection) startWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("Writer exit connID=", c.ConnID) // 添加退出日志

	for {
		select {
		case data, ok := <-c.msgChan: // 添加通道状态检测
			if !ok { // 通道已关闭
				fmt.Println("msgChan is closed")
				return
			}
			if _, err := c.GetTCPConnection().Write(data); err != nil {
				fmt.Println("send data error", err)
				c.ExitChan <- true // 触发连接关闭
				return
			}
		case <-c.ExitChan:
			// 关闭前发送剩余消息
			close(c.msgChan)
			for data := range c.msgChan {
				c.Conn.Write(data)
			}
			return
		}
	}
}
