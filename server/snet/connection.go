package snet

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"net_game/server/siface"
	"sync"
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
	// 告知当前连接已经退出的channel
	ExitChan chan bool
	// 新增写通道
	msgChan chan []byte
	// 处理该链接的方法router
	msgHandel siface.IMsgHandle
	// 由于每创建一个连接就需要添加到连接管理器这里集成了server
	TcpServer siface.IServer
	//连接属性
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.RWMutex
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.ConnID)
	// 启动从当前连接的读数据的业务

	// 启动读协程
	go c.startReader()

	// 启动写协程
	go c.startWriter()
	// 调用hook函数
	c.TcpServer.CallConnStart(c)

}

func (c *Connection) Stop() {
	fmt.Println("conn Stop()...ConnID=", c.ConnID)
	if c.isClosed == true {

		return
	}

	c.isClosed = true
	// 关闭socket连接
	c.TcpServer.CallConnStop(c)

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从连接管理器中删除
	c.TcpServer.GetManger().Remove(c)
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
	c.Conn.Close()
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

func (c *Connection) SendMsg(msgid uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	pack := NewDataPack()
	// 将data进行封包
	msg, err := pack.Pack(NewMsgPackage(msgid, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgid)
		return errors.New("pack error msg")
	}
	// 将数据写入通道，由写协程处理
	c.msgChan <- msg
	return nil
}

func NewConnection(server siface.IServer, conn *net.TCPConn, connID uint32, router siface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		//handleAPI: handleAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte), // 新增写通道初始化
		msgHandel: router,
		property:  make(map[string]interface{}),
	}

	// 将当前连接添加到连接管理中
	c.TcpServer.GetManger().Add(c)
	return c
}

// 新增读协程实现
func (c *Connection) startReader() {

	// logrus.WithFields(logrus.Fields{"conn_id": c.ConnID, "remote":  c.RemoteAddr().String(),}).Debug("Reader goroutine started")
	zap.S().Debugw("Reader goroutine started", "conn_id ", c.ConnID, "remote addr", c.RemoteAddr().String())

	defer func() {
		//logrus.WithFields(logrus.Fields{ "conn_id": c.ConnID,}).Warn("Reader exiting")
		zap.S().Debugw("Reader exiting", "conn_id ", c.ConnID, "remote addr", c.RemoteAddr().String())
		c.Stop()
	}()
	for {
		pack := NewDataPack()
		// 读取客户端的Msg Head
		headData := make([]byte, pack.GetHeadLen())
		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Println("")
			break
		}
		// 将读到的头部数据进行拆包到msg中
		msg, err := pack.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据dataLen再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, er := io.ReadFull(c.Conn, data)
			if er != nil {
				fmt.Println("read msg data error", er)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		// 判断线程池是否开启
		if c.msgHandel.IsWorkerPoolStarted() {
			// 将消息发送给消息队列
			c.msgHandel.SendMsgToTaskQueue(&req)
		} else {
			// 直接执行  创建一个特有的workerId 表示没有走线程处理池
			c.msgHandel.DoMsgHandler(ErrWorkerId, &req)
		}

	}
}

// 新增写协程实现
func (c *Connection) startWriter() {
	//logrus.WithField("conn_id", c.ConnID).Debug("Writer goroutine started")
	zap.S().Debugw("Writer goroutine started", "conn_id ", c.ConnID, "remote addr", c.RemoteAddr().String())
	//defer logrus.WithField("conn_id", c.ConnID).Warn("Writer exiting")
	defer zap.S().Debugw("Writer exiting ", "conn_id ", c.ConnID, "remote addr", c.RemoteAddr().String())

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
			//由于stop时触发的信号 这里就不用管了
			for data := range c.msgChan {
				c.Conn.Write(data)
			}
			return
		}
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
