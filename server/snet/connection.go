package snet

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
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
	msgChan   chan []byte
	PlayerRef *Player // ⚠️ 弱引用（非必须）
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

func NewConnection(conn *net.TCPConn, connID uint32) *Connection {
	c := &Connection{
		Conn:   conn,
		ConnID: connID,
		//handleAPI: handleAPI,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		msgChan:  make(chan []byte), // 新增写通道初始化
	}
	return c
}

// 新增读协程实现
func (c *Connection) startReader() {

	logrus.WithFields(logrus.Fields{
		"conn_id": c.ConnID,
		"remote":  c.RemoteAddr().String(),
	}).Debug("Reader goroutine started")

	defer func() {
		logrus.WithFields(logrus.Fields{
			"conn_id": c.ConnID,
		}).Warn("Reader exiting")
		c.Stop()
	}()
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		n, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read error", err)
			return
		}
		// 新增命令处理
		cmd := string(buf[:n])

		c.handleCommand(cmd)

	}
}

// 新增指令处理方法
func (c *Connection) handleCommand(cmd string) {
	// 示例：简单解析指令
	suffix := strings.TrimSuffix(cmd, "\n")
	if suffix == "cj" {
		if player := c.PlayerRef; player != nil {
			room, err := player.CreateRoom(2) // 创建房间，最多2个玩家
			if err != nil {
				logrus.Info("创建房间失败")
				return
			}

			c.SendMsg([]byte("创建房间成功"))
			room.CheckAndStart()
			logrus.WithFields(logrus.Fields{"房间ID是:": room.GetID(), "玩家id是": player.GetID()}).Info("创建房间成功")

		}
	}
	if suffix == "jr" {
		if player := c.PlayerRef; player != nil {
			// 尝试加入房间
			for _, room := range rooms {
				if !room.IsFull() {
					err := player.JoinRoom(room)
					// jion 里边有判断 如果慢的话就自动开始游戏
					if err == nil {
						c.SendMsg([]byte("加入房间成功"))
						room.CheckAndStart()
						logrus.WithFields(logrus.Fields{"房间ID是:": room.GetID(), "玩家id是": player.GetID()}).Info("加入房间成功")
						return
					}

				}
			}
			logrus.Warn("没有可加入的房间")
			c.SendMsg([]byte("没有可加入的房间"))
		}
	}
}

// 新增写协程实现
func (c *Connection) startWriter() {
	logrus.WithField("conn_id", c.ConnID).Debug("Writer goroutine started")

	defer logrus.WithField("conn_id", c.ConnID).Warn("Writer exiting")

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
