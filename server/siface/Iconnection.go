package siface

import "net"

/**
 * @Description
 * @Date 2025/3/6 23:33
 **/

// IConnection 定义连接模块的抽象接口
type IConnection interface {
	Start() // 启动连接 让当前的连接开始工作

	Stop() // 停止连接 结束当前连接的工作

	GetTCPConnection() *net.TCPConn // 获取当前连接绑定的socket conn

	GetConnID() uint32 // 获取当前连接模块的连接ID

	RemoteAddr() net.Addr // 获取远程客户端的TCP状态 IP port

	SendMsg(msgId uint32, data []byte) error // 发送数据  将数据发送给远程的客户端

	SetProperty(key string, value interface{})   // 设置连接属性
	GetProperty(key string) (interface{}, error) // 获取连接属性
	RemoveProperty(key string)                   // 删除连接属性

}

// 定义个处理链接业务的方法

//type HandleFunc func(*net.TCPConn, uint32, []byte, int) error
