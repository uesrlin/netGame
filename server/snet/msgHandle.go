package snet

import (
	"fmt"
	"go.uber.org/zap"
	"net_game/server/siface"
)

/**
 * @Description
 * @Date 2025/3/16 22:11
 **/

type MsgHandle struct {
	Apis map[uint32]siface.IRouter
	// 加入线程池 线程池获取消息队列然后处理
	WorkPoolSize uint32
	//加入消息队列, 防止阻塞  消息队列存储请求
	TaskQueue []chan siface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]siface.IRouter),
		// 初始化线程池  只有10个
		WorkPoolSize: WorkerPoolSize,
		// 初始化消息队列  定为10
		TaskQueue: make([]chan siface.IRequest, WorkerPoolSize),
	}
}
func (m *MsgHandle) DoMsgHandler(workerId int, request siface.IRequest) {
	router, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("router is not exist")
		return
	}
	zap.S().Debugw("线程开始处理请求", "request MsgID = ", request.GetMsgId(), "To WorkerID = ", workerId)
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
	zap.S().Debugw("线程处理完毕", "request MsgID = ", request.GetMsgId(), "To WorkerID = ", workerId)
}

func (m *MsgHandle) AddRouter(msgId uint32, router siface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		zap.S().Debug("router has exist")
		return
	}
	m.Apis[msgId] = router
}

// 根据workPoolSize创建线程池，启动每一个worker线程

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkPoolSize); i++ {
		m.TaskQueue[i] = make(chan siface.IRequest, MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])

	}
}

// 启动一个worker的工作流程

func (m *MsgHandle) StartOneWorker(workerId int, taskQueue chan siface.IRequest) {
	fmt.Println("workerId = ", workerId, "is started...")
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(workerId, request)
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request siface.IRequest) {
	// 得到连接的id
	connId := request.GetConnection().GetConnID()
	// 得到连接的id 然后根据连接id 得到线程池的id
	workerId := connId % m.WorkPoolSize
	zap.S().Debugw("请求开始分配线程池", "conn_id ", connId, "request MsgID = ", request.GetMsgId(), "To WorkerID = ", workerId)
	// 将请求加入到线程池的消息队列中
	m.TaskQueue[workerId] <- request

}

// 判断线程池是否开启

func (m *MsgHandle) IsWorkerPoolStarted() bool {
	for _, queue := range m.TaskQueue {
		if queue == nil {
			return false
		}
	}
	return true
}
