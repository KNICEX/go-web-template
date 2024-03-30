package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-web-template/pkg/conf"
	"time"
)

var DB *gorm.DB

func Init() {
	var dsn string
	switch conf.DatabaseConf.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			conf.DatabaseConf.User,
			conf.DatabaseConf.Password,
			conf.DatabaseConf.Host,
			conf.DatabaseConf.Port,
			conf.DatabaseConf.Name,
		)
	default:
		panic("unknown database type")
	}

	var err error
	DB, err = gorm.Open(conf.DatabaseConf.Type, dsn)
	if err != nil {
		panic(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.DatabaseConf.TablePrefix + defaultTableName
	}

	if conf.AppConf.Debug {
		DB.LogMode(true)
	} else {
		DB.LogMode(false)
	}

	DB.DB().SetMaxIdleConns(conf.DatabaseConf.MaxIdleConns)
	DB.DB().SetMaxOpenConns(conf.DatabaseConf.MaxOpenConns)

	DB.DB().SetConnMaxLifetime(time.Second * 30)

	migration()
}
