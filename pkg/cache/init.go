package cache

import (
	"fmt"
	"my-bluebell/pkg/conf"
	"my-bluebell/pkg/logger"
)

var Store Driver = NewLocalStore()

func Init() {
	if conf.RedisConf.Host != "" {
		Store = NewRedisStore(
			conf.RedisConf.Network,
			fmt.Sprintf("%s:%d", conf.RedisConf.Host, conf.RedisConf.Port),
			conf.RedisConf.User,
			conf.RedisConf.Password,
			conf.RedisConf.DB,
			conf.RedisConf.PoolSize,
		)
	}
	err := Store.Ping()
	if err != nil {
		logger.Panic(err)
	}
	if err = Store.Restore(DefaultCacheFile); err != nil {
		logger.Warn(err)
	}
}
