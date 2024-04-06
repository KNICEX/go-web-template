package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-web-template/pkg/conf"
	"go-web-template/pkg/logger"
)

var Store Driver = NewLocalStore()

func Init() {
	if conf.RedisConf.Host != "" {
		rc := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.RedisConf.Host, conf.RedisConf.Port),
			Password: conf.RedisConf.Password,
			DB:       conf.RedisConf.DB,
			PoolSize: conf.RedisConf.PoolSize,
			Username: conf.RedisConf.User,
		})
		Store = NewRedisStore(rc, rc.Close)
	}
	err := Store.Ping()
	if err != nil {
		logger.L().Panic("cache ping error ", err)
	}
	if err = Store.Restore(DefaultCacheFile); err != nil {
		logger.L().Warn(err)
	}
}
