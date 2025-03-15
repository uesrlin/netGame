package app

import (
	"context"
	"net_game/server/internal/config/app"
	mysql3 "net_game/server/internal/config/mysql"
	redis2 "net_game/server/internal/config/redis"
	"net_game/server/internal/db/mysql"
	"net_game/server/internal/db/redis"
	"net_game/server/util/file"
	"net_game/server/util/path"
)

/**
 * @Description
 * @Date 2025/3/15 19:43
 **/

type App struct {
	*app.ConfigMap
	//logger       logger.CustomLogger
	redisManager *redis.Manager
	mysqlManager *mysql.Manager
}

var appInstance *App = nil

func (a *App) Init(ctx context.Context, configDir string) {

	// 初始化配置表
	appConfigData := file.ReadDataFromPath(path.JoinPath(configDir, "app.yaml"))
	a.ConfigMap = app.InitAppConfigMap(appConfigData)

	// 初始化Mysql
	mysqlConfigData := file.ReadDataFromPath(path.JoinPath(configDir, "mysql.yaml"))
	mysqlDBConfig := mysql3.InitDBConfigMap(mysqlConfigData)
	a.mysqlManager = mysql.NewManager(ctx, mysqlDBConfig)

	//初始化Redis
	redisConfigData := file.ReadDataFromPath(path.GetPath("config/redis.yaml"))
	redisDBConfig := redis2.InitDBConfigMap(redisConfigData)
	a.redisManager = redis.NewManager(ctx, redisDBConfig)
	// 设置一下数据
	appInstance = a
}

// ------------------------------------------MySQL管理器接口------------------------------------------//

func GetClient(name string) *mysql.Client {
	return appInstance.mysqlManager.GetClient(name)
}
