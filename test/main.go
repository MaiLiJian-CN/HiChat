package main

import (
	"HiChat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Welcome12@tcp(127.0.0.1:3306)/HiChat?charset=utf8mb4&parseTime=True&loc=Local"
	db,err:=gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err!=nil{
		panic(err)
	}
	err=db.AutoMigrate(&models.UserBasic{})
	if err!=nil{
		panic(err)
	}
	err=db.AutoMigrate(&models.Relation{})
	if err!=nil{
		panic(err)
	}
	err=db.AutoMigrate(&models.Community{})
	if err!=nil{
		panic(err)
	}
	err=db.AutoMigrate(&models.Message{})
	if err!=nil {
		panic(err)
	}
}	
