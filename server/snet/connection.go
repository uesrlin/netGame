package snet

import (
	"errors"
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
	handleAPI siface.HandleFunc
	// 告知当前连接已经退出的channel
	ExitChan chan bool
	// 新增写通道
	msgChan chan []byte
}

func (c Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.ConnID)
	// 启动从当前连接的读数据的业务

	// 启动读协程
	go c.startReader()

	// 启动写协程
	go c.startWriter()

}

func (c Connection) Stop() {
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

func (c Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn

}

func (c Connection) GetConnID() uint32 {
	return c.ConnID

}

func (c Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()

}

func (c Connection) SendMsg(data []byte) error {
	// 将数据写入通道，由写协程处理
	c.msgChan <- data
	return nil
}

func NewConnection(conn *net.TCPConn, connID uint32, handleAPI siface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: handleAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte, 1024), // 新增写通道初始化
	}
	return c
}

func CallBackToClient(conn *net.TCPConn, cid uint32, data []byte, cnt int) error {
	// 回显业务
	fmt.Print("[Conn Handle cid is ]", cid, " CallbackToClient...")
	rece := "收到了,data is："
	byt := []byte(rece)
	if _, err := conn.Write(append(byt, data[:cnt]...)); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	fmt.Println("接收的数据是:", string(data[:cnt]))
	return nil

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
		// 处理读取到的数据（示例调用处理函数）
		if er := c.handleAPI(c.Conn, c.ConnID, buf, n); er != nil {
			fmt.Println("Conn is ", c.ConnID, "handle error", er)
			break
		}

	}
}

// 新增写协程实现
func (c *Connection) startWriter() {
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}
