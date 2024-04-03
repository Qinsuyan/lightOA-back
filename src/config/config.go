package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 导出的配置
var Log *log
var Http *http
var DBMysql *Mysql

// configure 配置文件
type configure struct {
	Log   log
	Http  http
	Mysql Mysql
}

// LOG 解析日志的配置
type log struct {
	Enable bool
	Level  string
}

type http struct {
	Enable bool
	Dist   string
	Port   string
}
type Mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

// Load 解析配置文件
func Load(dir string) error {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(dir)      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		return fmt.Errorf("exception when ReadInConfig: %s", err)
	}
	//parse
	var conf configure
	err = viper.Unmarshal(&conf)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %v", err)
	}
	//assign
	Log = &conf.Log
	Http = &conf.Http
	DBMysql = &conf.Mysql
	return nil
}
