package main

import (
	"flag"
	"fmt"
	"net_game/server/internal/logger/zapLog"
	"net_game/server/siface"
	"net_game/server/snet"
)

type PingRouter struct {
	snet.BaseRouter
}

func (p *PingRouter) PreHandle(request siface.IRequest) {
	fmt.Println("before")
	err := request.GetConnection().SendMsg([]byte(" received message  before\n"))
	if err != nil {
		fmt.Println("call back ping error")
		return
	}
}

func (p *PingRouter) Handle(request siface.IRequest) {
	fmt.Printf("Addr :%s,receive message conn is:%d ,data is %s", request.GetConnection().RemoteAddr(), request.GetConnection().GetConnID(), string(request.GetData()))
	err := request.GetConnection().SendMsg([]byte("server received message\n"))
	if err != nil {
		fmt.Println("call back ping error")
		return
	}
}

func (p *PingRouter) PostHandle(request siface.IRequest) {
	fmt.Println("after")
	err := request.GetConnection().SendMsg([]byte(" received message  after\n"))
	if err != nil {
		fmt.Println("call back ping error")
		return
	}
}

func main() {
	configDir := flag.String("ConfigDir", "server/config/", "配置表路径必须设置")
	flag.Parse()
	fmt.Println("configDir: ", *configDir)

	app := snet.Server{}
	server := app.Init(nil, *configDir)
	//logrusLog.InitLogrus(app.GetLogFilePath(), app.GetLogFileName())
	zapLog.InitZapLog(app.GetLogFilePath(), app.GetLogFileName())

	//// 因为现在server里边是单个的 所以现在只能实现一个路由 如果多个的话就需要使用routerManger实现
	server.AddRouter(&PingRouter{})
	server.Serve()

}
