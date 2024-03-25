package bootstrap

import (
	"go-web-template/pkg/cache"
	"go-web-template/pkg/conf"
	"go-web-template/pkg/logger"
)

func Init() {
	conf.Init()
	logger.Init()
	cache.Init()
}
