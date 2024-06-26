package main

import (
	"lightOA-end/src/api"
	"lightOA-end/src/config"
	"lightOA-end/src/db"
	"lightOA-end/src/log"
	"lightOA-end/src/storage"
	"os"
	"os/signal"
	"time"
)

func main() {
	//固定时区
	time.Local = time.FixedZone("utc8", 8*3600)
	//加载配置项
	if err := config.Load("./"); err != nil {
		panic(err)
	}
	//log配置
	if config.Log.Enable {
		log.Setup(config.Log.Level)
	}
	//连接数据库
	err := db.Init(config.DBMysql)
	if err != nil {
		log.Err(err).Msg("failed to establish db connection")
		return
	}
	//初始化file
	if config.Storage.Path == "" {
		log.Err(nil).Msg("storage file path is empty")
	}
	err = storage.Init(config.Storage.Path, config.Storage.Password)
	if err != nil {
		log.Err(err).Msg("failed to init file storage")
		return
	}
	// http
	if config.Http.Enable {
		err = api.Start(config.Http.Port, config.Http.Dist)
		if err == nil {
			log.Info().Msgf("start to listen http at %s", config.Http.Port)
		} else {
			log.Err(err).Msg("err while starting http")
		}
	}
	chanQuit := make(chan os.Signal, 1)
	signal.Notify(chanQuit, os.Interrupt)
	<-chanQuit
}
