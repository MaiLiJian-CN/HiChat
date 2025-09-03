package global

import (
	"HiChat/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	RedisDB *redis.Client
	ServiceConfig *config.ServiceConfig
)
