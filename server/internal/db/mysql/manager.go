package mysql

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mysql3 "net_game/server/internal/config/mysql"
	"time"
)

type Manager struct {
	mysqlDBMap map[string]*Client
}

func (m *Manager) addClientWithName(name string, client *Client) {
	_, ok := m.mysqlDBMap[name]
	if ok {
		panic(fmt.Sprintf("duplication register client. name is %s.", name))
	}
	m.mysqlDBMap[name] = client
}

func (m *Manager) GetClient(name string) *Client {
	value, ok := m.mysqlDBMap[name]
	if !ok {
		return nil
	}
	return value
}

func NewManager(ctx context.Context, mysqlMysqlConfigMap *mysql3.ConfigMap) *Manager {
	manager := &Manager{
		mysqlDBMap: make(map[string]*Client, 5),
	}
	for key, value := range mysqlMysqlConfigMap.Dbs {

		db, err := gorm.Open(mysql.Open(value.ApplyURL()))
		if err != nil {
			panic(errors.New(fmt.Sprintf("[Mysql][NewManager] Error %+v", err)))
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(errors.New(fmt.Sprintf("[Mysql][NewManager] DB Error %+v", err)))
		}

		// SetMaxIdleCons 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(value.MinPoolSize)

		// SetMaxOpenCons 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(value.MaxPoolSize)

		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Minute)

		manager.addClientWithName(key, &Client{db})
	}
	return manager
}
