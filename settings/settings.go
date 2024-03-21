package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./web_app/")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	return
}

func OnConfigChange(callback func(e fsnotify.Event)) {
	viper.WatchConfig()
	viper.OnConfigChange(callback)
}
