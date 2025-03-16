package snet

import (
	"fmt"
	"net_game/server/siface"
)

/**
 * @Description
 * @Date 2025/3/16 22:11
 **/

type MsgHandle struct {
	Apis map[uint32]siface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]siface.IRouter),
	}
}
func (m *MsgHandle) DoMsgHandler(request siface.IRequest) {
	router, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("router is not exist")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router siface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		fmt.Println("router has exist")
		return
	}
	m.Apis[msgId] = router
}
