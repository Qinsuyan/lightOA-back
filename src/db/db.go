package db

import (
	"encoding/json"
	"fmt"
	"io"
	"lightOA-end/src/config"
	"lightOA-end/src/entity"
	"lightOA-end/src/log"
	"os"
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
	err = createTables()
	if err != nil {
		return err
	}
	err = createRootResource()
	if err != nil {
		return err
	}
	err = addSuperRole()
	if err != nil {
		return err
	}
	err = addSuperUser()
	if err != nil {
		return err
	}
	err = addDefaultResource()
	if err != nil {
		return err
	}
	return nil
}

// 初始化数据库
func createTables() error {
	err := con.Sync(
		new(entity.UserRaw),
		new(entity.Online),
		new(entity.ResourceRaw),
		new(entity.RoleRaw),
		new(entity.RoleResource),
		new(entity.UserLog),
		new(entity.SystemVariableInts),
		new(entity.SystemVariableTexts),
		new(entity.Department),
	)
	if err != nil {
		log.Err(err).Msg("err while syncing database")
	}
	return err
}

// 创建根资源
func createRootResource() error {
	_, err := con.Exec("insert ignore into resource_raw(id,name,alias,type,createdAt) values (1,'根节点','ROOT',3,'2023-01-01 00:00:00')")
	if err != nil {
		log.Err(err).Msg("err while creating root resource")
	}
	return err
}
func addDefaultResource() error {
	file, err := os.Open("./resources.json")
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	var resources []*entity.Resource
	err = json.Unmarshal(bytes, &resources)
	if err != nil {
		return err
	}
	return insertResources(resources, 1)
}
func insertResources(resources []*entity.Resource, parentId int) error {
	for _, resource := range resources {
		existing := entity.ResourceRaw{
			Alias: resource.Alias,
		}
		has, err := con.Get(&existing)
		if err != nil {
			return err
		}
		var pId int
		if !has {
			newResource := entity.ResourceRaw{
				Name:      resource.Name,
				Alias:     resource.Alias,
				Type:      resource.Type,
				ParentId:  parentId,
				CreatedAt: time.Now(),
			}
			_, err = con.Insert(&newResource)
			if err != nil {
				return err
			}
			pId = newResource.Id
		} else {
			pId = existing.Id
		}
		if len(resource.Children) > 0 {
			err = insertResources(resource.Children, pId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func addSuperUser() error {
	_, err := con.Exec("insert ignore into user_raw(id,username,password,role,createdAt) values (1,'admin','8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',1,'2023-01-01 00:00:00')")
	if err != nil {
		return err
	}
	return nil
}
func addSuperRole() error {
	_, err := con.Exec("insert ignore into role_raw(id,name,description,createdAt) values (1,'admin','超级管理员','2023-01-01 00:00:00')")
	if err != nil {
		return err
	}
	return nil
}
