package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("客户端连接成功!")

	// 创建消息接收协程
	go receiveMessages(conn)

	// 主循环处理用户输入
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\renter message: ") // 使用\r保持在行首
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(message))
	}
}

// 持续接收服务器消息的独立函数
func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\n连接已关闭:", err.Error())
			os.Exit(0)
		}
		// 加入这一行
		fmt.Printf("\r\x1b[K") // \x1b[K 清除当前行
		fmt.Printf("服务器回复: %s\n", response)
		fmt.Print("enter message: ") // 保持输入提示可见
	}
}
