package snet

import (
	"errors"
	"github.com/sirupsen/logrus"
	"userlin/netGame/server/siface"
)

/**
 * @Description
 * @Date 2025/3/13 10:31
 **/

type Room struct {
	RoomId    int
	players   []*Player
	isFull    bool
	isWaiting bool
	// 房间最大人数
	maxNum int
	// 棋盘
	board [][]int
}

func (r *Room) GetID() int {
	return r.RoomId
}

func (r *Room) Leave(player siface.IPlayer) {
	// 实现离开房间逻辑
	for i, p := range r.players {
		if p == player {
			r.players = append(r.players[:i], r.players[i+1:]...)
			break
		}
	}
}

func (r *Room) Broadcast(message []byte) {
	for _, player := range r.players {
		player.GetConnection().SendMsg(message)
	}
}

func (r *Room) IsFull() bool {
	return len(r.players) >= r.maxNum
}

func (r *Room) GetPlayerCount() int {
	return len(r.players)
}

func (r *Room) Close() {
	// 关闭房间，释放资源
	r.isWaiting = false
	r.isFull = true
	r.players = nil
}

func (r *Room) GetPlayers() []siface.IPlayer {
	//获取玩家的列表
	players := make([]siface.IPlayer, len(r.players))
	for i, player := range r.players {
		players[i] = player
	}
	return players
}

func (r *Room) Join(player siface.IPlayer) error {
	if r.isFull {
		return errors.New("room is full")
	}
	r.players = append(r.players, player.(*Player))
	r.isFull = true
	return nil
}

// StartGame 在Room结构体添加方法
func (r *Room) StartGame() {
	if len(r.players) < 2 {
		return // 最少需要2名玩家
	}

	// 初始化棋盘
	for i := range r.board {
		r.board[i] = make([]int, 15)
	}
	// 通知所有玩家游戏开始
	r.Broadcast([]byte("游戏开始！"))
	logrus.Info("游戏开始了")

	r.isWaiting = false
}

// CheckAndStart 检查房间人数是否足够并开始游戏
func (r *Room) CheckAndStart() {
	if r.isFull {
		r.StartGame()
	}

}
