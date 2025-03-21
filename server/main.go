package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net_game/server/internal/logger/zapLog"
	"net_game/server/siface"
	"net_game/server/snet"
)

type PingRouter struct {
	snet.BaseRouter
}

//
//func (p *PingRouter) PreHandle(request siface.IRequest) {
//	zap.S().Debugw("请求预处理", "request MsgID = ", request.GetMsgId())
//
//}

func (p *PingRouter) Handle(request siface.IRequest) {

	zap.S().Debugw("请求正在处理",
		"Addr is", request.GetConnection().RemoteAddr(),
		"connId is ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgId(),
		"data is ", string(request.GetData()))

	// 由于消息结构进行了封装，所以需要进行重写
	err := request.GetConnection().SendMsg(0, []byte("ping ...ping...ping"))
	if err != nil {
		fmt.Println("call back ping error")
		return
	}
}

func (p *PingRouter) PostHandle(request siface.IRequest) {
	zap.S().Debugw("请求追加处理", "request MsgID = ", request.GetMsgId())
}

type Hello struct {
	snet.BaseRouter
}

func (p *Hello) PreHandle(request siface.IRequest) {
	zap.S().Debugw("请求预处理", "request MsgID = ", request.GetMsgId())

}

func (p *Hello) Handle(request siface.IRequest) {

	zap.S().Debugw("请求正在处理",
		"Addr is", request.GetConnection().RemoteAddr(),
		"connId is ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgId(),
		"data is ", string(request.GetData()))

	// 由于消息结构进行了封装，所以需要进行重写
	err := request.GetConnection().SendMsg(1, []byte("hello ...hello ...hello"))
	if err != nil {
		fmt.Println("call back ping error")
		return
	}
}

//func (p *Hello) PostHandle(request siface.IRequest) {
//	zap.S().Debugw("请求处理结束后", "request MsgID = ", request.GetMsgId())
//
//}

// 连接创建之后做的事

func DoConnBegin(conn siface.IConnection) {
	fmt.Println("========>  DoConnBegin is Called ")
	if err := conn.SendMsg(202, []byte("DoConnBegin")); err != nil {
		fmt.Println(err)
	}
	// 给当前的连接设置一些属性
	conn.SetProperty("Name", "张三")
	conn.SetProperty("Home", "北京")

}

// 连接断开之后做的事

func DoConnLost(conn siface.IConnection) {
	fmt.Println("========>  DoConnLost is Called ")
	fmt.Println("conn Id=", conn.GetConnID(), "is Lost ")
	// 获取当前连接的属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name=", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home=", home)
	}

}

func main() {
	configDir := flag.String("ConfigDir", "server/config/", "配置表路径必须设置")
	flag.Parse()
	fmt.Println("configDir: ", *configDir)

	app := snet.Server{}
	server := app.Init(nil, *configDir)
	zapLog.InitZapLog(app.GetLogFilePath(), app.GetLogFileName())

	// 注册钩子函数
	server.SetConnStart(DoConnBegin)
	server.SetConnStop(DoConnLost)

	//// 因为现在server里边是单个的 所以现在只能实现一个路由 如果多个的话就需要使用routerManger实现
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &Hello{})
	server.Serve()

}
