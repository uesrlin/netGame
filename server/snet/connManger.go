package snet

import (
	"net_game/server/siface"
	"sync"
)

/**
 * @Description
 * @Date 2025/3/19 10:20
 **/

type ConnManager struct {
	Connections map[uint32]siface.IConnection
	// 加锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]siface.IConnection),
	}
}

func (c *ConnManager) Add(conn siface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.Connections[conn.GetConnID()] = conn

}

func (c *ConnManager) Remove(conn siface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.Connections, conn.GetConnID())

}

func (c *ConnManager) Get(connId uint32) (siface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.Connections[connId]; ok {
		return conn, nil
	}
	return nil, nil
}

func (c *ConnManager) Len() int {
	return len(c.Connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for connId, conn := range c.Connections {
		conn.Stop()
		delete(c.Connections, connId)
	}

}
