package initialize

import (
	"HiChat/global"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root",
	// 	"Welcome12", "127.0.0.1", 3306, "HiChat")
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.ServiceConfig.DB.User,
		global.ServiceConfig.DB.Password, global.ServiceConfig.DB.Host, global.ServiceConfig.DB.Port, global.ServiceConfig.DB.Name)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}

func InitRedis(){
	opt:=redis.Options{
		Addr: fmt.Sprintf("%s:%d",global.ServiceConfig.RedisDB.Host,global.ServiceConfig.RedisDB.Port),
		Password: "",
		DB:0,
	}
	global.RedisDB=redis.NewClient(&opt)
}