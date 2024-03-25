package email

import (
	"go-web-template/pkg/conf"
	"go-web-template/pkg/logger"
)

var Client Driver

func Init() {
	if conf.EmailConf.User == "" {
		return
	}

	logger.Debug("Initializing email sending queue...")

	Client = NewSMTPClient()
}
