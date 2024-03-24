package utils

import (
	"Adorable/configs"
	"Adorable/pkg/cache"
	"Adorable/pkg/log"
)

var (
	tokenRedis *cache.RedisCache
	smsRedis   *cache.RedisCache
)

// InitRedisStore 初始化Redis存储.
func InitRedisStore(Config *configs.ServerConfig) (err error) {
	tokenRedis, err = cache.NewRedisCache(Config.TokenRedis)
	if err != nil {
		log.Fatalf("NewRedisCache(%s) failed! err: %v", Config.TokenRedis.DebugStr(), err)
		return
	}
	smsRedis, err = cache.NewRedisCache(Config.SmsRedis)
	if err != nil {
		log.Fatalf("NewRedisCache(%s) failed! err: %v", Config.SmsRedis.DebugStr(), err)
		return
	}

	return
}
