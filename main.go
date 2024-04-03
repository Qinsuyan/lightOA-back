package main

import (
	"lightOA-end/src/config"
	"lightOA-end/src/log"
	"lightOA-end/src/mydb"
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
	err := mydb.Init(config.DBMysql)
	if err != nil {
		log.Err(err).Msg("failed to establish db connection")
		return
	}
	chanQuit := make(chan os.Signal, 1)
	signal.Notify(chanQuit, os.Interrupt)
	<-chanQuit
}
