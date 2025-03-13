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

	// 主循环处理用户输入
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\renter message: ") // 使用\r保持在行首
		message, _ := reader.ReadString('\n')
		conn.Write([]byte(message))
		go func() {
			for {
				buf := make([]byte, 512)
				n, er := conn.Read(buf)
				if er != nil {
					fmt.Println("读取错误:", err)
					return
				}
				fmt.Println("收到消息:", string(buf[:n]))
			}
		}()
	}

}
