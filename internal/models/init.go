package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"op-panel-go/define"
)

var DB *gorm.DB

func NewDB() {
	dsn := define.MysqlDNS + "/op_panel?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm Init Error : ", err)
	}
	err = db.AutoMigrate(&ConfigBasic{}, &TaskBasic{}, &SoftBasic{}, &WebBasic{}, &ExecuteQueue{})
	if err != nil {
		panic("[MIGRATE ERROR] : " + err.Error())
	}
	DB = db
}
