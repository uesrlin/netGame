package snet

import (
	"errors"
	"sync"
	"userlin/netGame/server/siface"
	rd "userlin/netGame/server/util/random"
)

/**
 * @Description
 * @Date 2025/3/13 10:31
 **/

type Coordinate struct {
	X float32
	Y float32
}

type Player struct {
	conn         siface.IConnection
	Room         siface.IRoom
	playerId     int32
	specialPosit [2]int
	playerName   string
	historyCoord map[Coordinate]bool // 新增坐标记录

}

var (
	rooms         = make(map[int]*Room)
	roomLock      sync.Mutex
	roomIDCounter int
)

func (p *Player) GetID() int32 {
	return p.playerId

}

func (p *Player) Send(dar []byte) error {
	err := p.GetConnection().SendMsg(dar)
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) GetRoom() siface.IRoom {
	return p.Room
}

func (p *Player) LeaveRoom() {
	if p.Room == nil {
		return
	}
	p.Room.Leave(p)
	p.Room = nil

}

func (p *Player) Disconnect() {
	// 断开连接
	p.conn.Stop()
	// 离开房间
	p.LeaveRoom()
}

func (p *Player) SetName(name string) {
	p.playerName = name

}

func (p *Player) GetConnection() siface.IConnection {
	return p.conn
}

var pid int32 = 0

func NewPlayer(conn siface.IConnection) *Player {
	pid++
	return &Player{
		conn: conn,
		// 在玩家加入room后才初始化
		Room: nil,
		// 这里需要玩家随机生成一个id
		playerId:     pid,
		playerName:   rd.GenerateRandomName(3),
		historyCoord: make(map[Coordinate]bool),
	}
}

// CreateRoom 创建房间方法（通过Player调用）
func (p *Player) CreateRoom(maxPlayers int) (*Room, error) {
	roomLock.Lock()
	defer roomLock.Unlock()
	roomIDCounter++
	if p.Room != nil {
		p.GetConnection().SendMsg([]byte("您已在其他房间，不能在创建房间"))
		return nil, errors.New("创建房间失败 已有房间")
	}

	// 检查是否已经有房间
	if len(rooms) > 0 {
		for _, room := range rooms {
			if !room.IsFull() {
				// 如果已经有房间存在，尝试加入未满的房间
				if err := room.Join(p); err == nil {
					p.Room = room
					// 发送消息通知玩家
					if err := p.GetConnection().SendMsg([]byte("已加入未满的房间")); err != nil {
						return nil, err // 发送消息失败
					}
					return room, nil
				}
				return nil, errors.New("加入房间失败") // 加入房间失败
			}
		}
	}

	// 创建新房间
	newRoom := &Room{
		RoomId:    roomIDCounter,
		maxNum:    maxPlayers,
		players:   []*Player{p},
		isWaiting: true,
		board:     make([][]int, 15), // 15x15棋盘
	}

	p.Room = newRoom
	rooms[newRoom.RoomId] = newRoom
	return newRoom, nil

}

func (p *Player) JoinRoom(room siface.IRoom) error {
	if p.Room != nil {
		p.GetConnection().SendMsg([]byte("您已经在房间里了"))
		return errors.New("already in a room")
	}
	err := room.Join(p)
	if err != nil {
		return err
	}
	p.Room = room
	// 检查房间是否已满，若已满则开始游戏
	return nil
}
