package main

import (
	"userlin/netGame/server/internal/logger/logrusLog"
	"userlin/netGame/server/snet"
)

func main() {
	// 日志二选一都行
	logrusLog.InitLogrus()
	// zapLog.InitZapLog()
	server := snet.NewServer("127.0.0.1", 8080, "", "")
	server.Serve()

}
