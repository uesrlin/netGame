package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	redis2 "net_game/server/internal/config/redis"
	"time"
)

type Manager struct {
	redisDBMap map[string]*Client
}

func (m *Manager) addClientWithName(name string, client *Client) {
	_, ok := m.redisDBMap[name]
	if ok {
		panic(fmt.Sprintf("duplication register client. name is %s.", name))
	}
	m.redisDBMap[name] = client
}

func (m *Manager) GetClient(name string) *Client {
	value, ok := m.redisDBMap[name]
	if !ok {
		return nil
	}
	return value
}

func NewManager(ctx context.Context, configMap *redis2.ConfigMap) *Manager {
	manager := &Manager{
		redisDBMap: make(map[string]*Client, 5),
	}
	for key, value := range configMap.Dbs {
		rdbClient := redis.NewClient(&redis.Options{
			Addr:         value.ApplyURL(),
			Password:     value.Passwd,
			PoolSize:     int(value.MaxPoolSize),
			MinIdleConns: int(value.MinPoolSize),
			IdleTimeout:  time.Second * time.Duration(value.Timeout),
		}).WithContext(ctx)
		result := rdbClient.Ping(ctx)
		if result.Err() != nil {
			panic(errors.Wrap(result.Err(), "NewManager rdbClient.Ping()"))
		}
		manager.addClientWithName(key, &Client{rdbClient})
	}
	return manager
}
