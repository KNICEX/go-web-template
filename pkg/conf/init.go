package conf

import (
	"flag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"strings"
)

var confPath string

var v *viper.Viper

func init() {
	flag.StringVar(&confPath, "conf", "remote.yaml", "config file path")
}

var cfg conf

var AppConf *appConf
var ServerConf *serverConf
var DatabaseConf *databaseConf
var RedisConf *redisConf
var LogConf *logConf

func setExpose() {
	AppConf = &cfg.App
	ServerConf = &cfg.Server
	DatabaseConf = &cfg.Database
	RedisConf = &cfg.Redis
	LogConf = &cfg.Log
}

func Init(path ...string) {
	v = viper.New()
	setDefault()
	// 读取本地基础配置
	if len(path) > 0 {
		confPath = path[0]
	}
	v.SetConfigFile(confPath)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("read config file error: %v", err)
	}
	if v.GetString("app.remote.addr") != "" {
		// 加载远程配置
		key := v.GetString("app.remote.path")
		if err = loadRemote(key); err != nil {
			log.Fatalf("load remote config error: %v", err)
		}
		// 加载远程配置的模式配置
		if v.GetString("app.mode") != "" {
			if err = loadRemote(key + "-" + v.GetString("app.mode")); err != nil {
				// 这里不需要 panic，模式对应的配置文件不是必须的
				log.Println("load remote mode config error: ", err)
			}
		}
	} else {
		// 加载本地模式配置
		if v.GetString("app.mode") != "" {
			dotIndex := strings.LastIndex(confPath, ".")
			modePath := confPath[:dotIndex] + "-" + v.GetString("app.mode") + confPath[dotIndex:]
			if err = loadLocal(modePath); err != nil {
				// 这里不需要 panic，模式对应的配置文件不是必须的
				log.Println("load local mode config error: ", err)
			}
		}
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unmarshal config error: %v", err)
	}
	setExpose()

}

func loadRemote(key string) (err error) {
	remoteV := viper.New()
	remoteV.SetConfigType(v.GetString("app.remote.type"))
	useSecure := v.GetString("app.remote.secret_key") != ""
	if err = setRemoteConf(remoteV, key, useSecure); err != nil {
		return
	}
	if err = remoteV.ReadRemoteConfig(); err != nil {
		return
	}
	return v.MergeConfigMap(remoteV.AllSettings())
}

func loadLocal(path string) (err error) {
	localV := viper.New()
	localV.SetConfigFile(path)
	if err = localV.ReadInConfig(); err != nil {
		return
	}
	return v.MergeConfigMap(localV.AllSettings())
}

func setRemoteConf(remoteV *viper.Viper, key string, useSecure bool) (err error) {
	if useSecure {
		err = remoteV.AddSecureRemoteProvider(
			v.GetString("app.remote.provider"),
			v.GetString("app.remote.addr"),
			key,
			v.GetString("app.remote.secret_key"),
		)
	} else {
		err = remoteV.AddRemoteProvider(
			v.GetString("app.remote.provider"),
			v.GetString("app.remote.addr"),
			key,
		)
	}
	return
}
