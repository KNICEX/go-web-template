package bootstrap

import (
	"go-web-template/models"
	"go-web-template/pkg/auth"
	"go-web-template/pkg/authstorage"
	"go-web-template/pkg/cache"
	"go-web-template/pkg/conf"
	"go-web-template/pkg/email"
	"go-web-template/pkg/logger"
	"go-web-template/pkg/snowflake"
)

func Init() {
	conf.Init()
	logger.Init()
	models.Init()
	cache.Init()
	auth.Init()
	email.Init()
	snowflake.Init()
	authstorage.Init()

}

func Shutdown() error {
	return cache.Store.Close()
}
