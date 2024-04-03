package mydb

import (
	"fmt"
	"lightOA-end/src/config"
	"lightOA-end/src/entity"
	"lightOA-end/src/log"
	"time"

	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
)

var con *xorm.Engine

func Init(conf *config.Mysql) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.Dbname)
	engine, err := xorm.NewEngine("mysql", dsn)
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(100)
	engine.SetConnMaxLifetime(time.Minute * 10)
	if err != nil {
		log.Err(err).Msg("err while connecting polardb")
		return err
	}
	if err := engine.Ping(); err != nil {
		log.Err(err).Msg("err while connecting polardb")
		return err
	}
	con = engine
	createTables()
	return nil
}

// 初始化数据库
func createTables() {
	err := con.Sync(new(entity.UserRaw), new(entity.Online), new(entity.ResourceRaw), new(entity.RoleRaw), new(entity.RoleResource), new(entity.UserLog))
	if err != nil {
		log.Err(err).Msg("err while syncing database")
	}
}
