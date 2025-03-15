package definition

import "net_game/server/internal/db/mysql"

// 配置读取器

type AppConfigReader interface {
	GetSalt() string
	GetDebug() bool
	GetLogFilePath() string
	GetLogFileName() string
	IsInternalService() bool
}

// Mysql管理器接口

type DBClientManager interface {
	GetClient(tableName string) *mysql.Client
}
