package siface

/**
 * @Description
 * @Date 2025/3/13 10:39
 **/

type IRoom interface {
	GetID() int
	Join(player IPlayer) error // 加入房间
	Leave(player IPlayer)      // 离开房间
	Broadcast(data []byte)     // 广播消息
	IsFull() bool              // 是否满员
	GetPlayerCount() int       // 当前人数
	Close()                    // 关闭房间
	GetPlayers() []IPlayer
	StartGame() // 获取所有玩家
	CheckAndStart()
}
