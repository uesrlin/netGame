package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	// 从标准输入读取数据
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("输入消息: ")
		message, _ := reader.ReadString('\n')

		// 发送数据到服务器
		conn.Write([]byte(message))

		// 读取服务器的响应
		response, er := bufio.NewReader(conn).ReadString('\n')
		if er != nil {
			fmt.Println("Error reading response:", er.Error())
			return
		}

		// 打印服务器的响应
		fmt.Print("服务器回复: " + response)
	}

}
