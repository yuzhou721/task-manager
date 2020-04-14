package models

import (
	"fmt"
	"log"
	"task/conf"

	"github.com/jinzhu/gorm"

	//引入mysql驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	err error
)

func init() {
	db, err = gorm.Open(conf.Config.Database.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Config.Database.User,
		conf.Config.Database.Password,
		conf.Config.Database.Host,
		conf.Config.Database.DbName,
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	//开启debugger模式
	// db.LogMode(true)

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return setting.DatabaseSetting.TablePrefix + defaultTableName
	// }

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

//InitTable 初始化table
func InitTable() {
	log.Println("start init table")
	if db == nil {
		log.Fatal("mysql connect error")
	}
	createTable(&Task{})
	log.Println("init Table complate")
}

func createTable(i ...interface{}) {
	for _, v := range i {
		exist := db.HasTable(v)
		if !exist {
			db.CreateTable(v)
		}
	}
}
