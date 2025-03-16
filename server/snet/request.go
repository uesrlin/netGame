package snet

import "net_game/server/siface"

/**
 * @Description
 * @Date 2025/3/8 12:28
 **/

type Request struct {
	conn siface.IConnection
	msg  siface.IMessage
}

func (r *Request) GetConnection() siface.IConnection {
	return r.conn
}
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
