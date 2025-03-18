package main

import (
	"fmt"
	"io"
	"net"
	"net_game/server/snet"
)

// 由于消息进行了封装 所以客户端也需要进行封装
func main() {

	// 链接远程服务器得到一个conn包
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client run fail")
		return
	}
	defer conn.Close()
	fmt.Println("客户端连接成功!")

	// 发送封包的消息
	dp := snet.NewDataPack()
	msgData, er := dp.Pack(snet.NewMsgPackage(0, []byte("hello")))
	if er != nil {
		fmt.Println("Pack error:", err)
		return
	}
	if _, er1 := conn.Write(msgData); er1 != nil {
		fmt.Println("write error ", er1)
		return
	}

	go func() {

		// 接下来是要进行读数据
		headData := make([]byte, dp.GetHeadLen())
		_, err2 := io.ReadFull(conn, headData)
		if err2 != nil {
			fmt.Println("read head error")
			return
		}
		// 这里只是从数据中达到了头部的前8个数据，故我们要进行拆包
		// 这里的headData是只包含了消息的长度和id,不包含数据，对headData进行读取
		msgHeader, er2 := dp.Unpack(headData)
		if er2 != nil {
			fmt.Println("server unpack error")
			return
		}
		if msgHeader.GetMsgLen() > 0 {
			// msg 是有数据的 需要进行二次读取
			// 2.第二次从conn中读，根据head中的datalen在读取data的内容
			msg := msgHeader.(*snet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			// 这里有一个疑问就是第二次读取的时候，不知道偏移量怎么进行读8个字节后边的数据,readFull函数默认是有偏移量的
			_, er3 := io.ReadFull(conn, msg.Data)
			if er3 != nil {
				fmt.Println("server unpack fail")
				return
			}
			fmt.Println("Recv------ MsgId:", msg.Id, "dataLen=", msg.DataLen, "data=", string(msg.Data))

		}

	}()

	select {}

}
