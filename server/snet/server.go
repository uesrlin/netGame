package snet

import (
	"context"
	"fmt"
	"net"
	"net_game/server/internal/config/app"
	"net_game/server/internal/db/mysql"
	"net_game/server/internal/db/redis"
	"net_game/server/siface"
	"net_game/server/util/file"
	"net_game/server/util/path"
)

/**
 * @Description
 * @Date 2025/3/6 22:40
 **/

// Server 服务器配置
type Server struct {
	*app.ConfigMap
	//logger       logger.CustomLogger
	redisManager *redis.Manager
	mysqlManager *mysql.Manager
	Router       siface.IRouter // 服务器路由
}

var appInstance *Server = nil

func (a *Server) Init(ctx context.Context, configDir string) *Server {

	// 初始化配置表
	appConfigData := file.ReadDataFromPath(path.JoinPath(configDir, "app.yaml"))
	a.ConfigMap = app.InitAppConfigMap(appConfigData)

	// 初始化Mysql
	//mysqlConfigData := file.ReadDataFromPath(path.JoinPath(configDir, "mysql.yaml"))
	//mysqlDBConfig := mysql3.InitDBConfigMap(mysqlConfigData)
	//a.mysqlManager = mysql.NewManager(ctx, mysqlDBConfig)

	//初始化Redis
	//redisConfigData := file.ReadDataFromPath(path.GetPath("config/redis.yaml"))
	//redisDBConfig := redis2.InitDBConfigMap(redisConfigData)
	//a.redisManager = redis.NewManager(ctx, redisDBConfig)
	// 设置一下数据
	appInstance = a
	return a
}

func (s *Server) Start() {
	// 这里写tcp服务器启动的逻辑
	var cid uint32 = 0
	//监听本地的端口
	listenAddr := fmt.Sprintf("%s:%s", s.GetHost(), s.GetListenPoint())
	fmt.Println(listenAddr)

	go func() {
		addr, er := net.ResolveTCPAddr("tcp", listenAddr)
		if er != nil {
			fmt.Println("解析本地地址失败", er.Error())
			return
		}
		listen, err := net.ListenTCP("tcp", addr)
		if err != nil {
			fmt.Println("启动监听失败", err.Error())
			return
		}
		fmt.Println("服务器启动成功....")
		// 记得关闭监听
		defer listen.Close()
		// 这里有多个连接，需要使用goroutine来处理

		for {
			// 建立连接
			conn, er0 := listen.AcceptTCP()
			if er0 != nil {
				fmt.Println("建立连接失败", er.Error())
				continue
			}
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()

		}

	}()

}

func (s *Server) Stop() {
	// 这里写tcp服务器停止的逻辑

}

func (s *Server) Serve() {
	// 这里写tcp服务器运行的逻辑
	s.Start()
	// 阻塞在这里
	select {}
}
func (s *Server) AddRouter(router siface.IRouter) {
	s.Router = router
}

func GetClient(name string) *mysql.Client {
	return appInstance.mysqlManager.GetClient(name)
}
