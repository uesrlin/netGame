package siface

/**
 * @Description
 * @Date 2025/3/13 10:39
 **/

type IPlayer interface {
	GetRoom() IRoom            // 当前所在房间
	JoinRoom(room IRoom) error // 加入房间
	LeaveRoom()

	GetID() int32               // 玩家唯一标识
	Send(data []byte) error     // 发送消息// 离开房间
	Disconnect()                // 断开连接
	SetName(name string)        // 设置玩家名
	GetConnection() IConnection // 获取网络连接
}
