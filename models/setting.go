package models

import (
	"encoding/gob"
	"github.com/jinzhu/gorm"
	"go-web-template/pkg/cache"
	"net/url"
)

const (
	SiteURLKey    = "site_url"
	SettingPrefix = "setting_"
	SecretKey     = "secret_key"
	JWTSecret     = "jwt_secret"
)

type Setting struct {
	gorm.Model
	Type  string `gorm:"not null"`
	Name  string `gorm:"unique;not null;index:setting_key"`
	Value string `gorm:"size:65535"`
}

func init() {
	gob.Register(Setting{})
}

func GetSettingByNames(names ...string) map[string]string {
	var queryRes []Setting
	res, miss := cache.GetStrs(names, SettingPrefix)
	if len(miss) > 0 {
		DB.Where("name in (?)", miss).Find(&queryRes)
		for _, setting := range queryRes {
			res[setting.Name] = setting.Value
		}
	}
	_ = cache.SetStrs(res, SettingPrefix, 0)
	return res
}

func GetSettingByName(name string) string {
	cacheV, ok := cache.Get(SettingPrefix + name)
	if ok {
		return cacheV.(string)
	}
	var setting Setting
	DB.Where("name = ?", name).First(&setting)
	_ = cache.Set("setting_"+name, setting.Value, 0)
	return setting.Value
}

func GetSiteURL() *url.URL {
	siteURL, err := url.Parse(GetSettingByName(SiteURLKey))
	if err != nil {
		return &url.URL{}
	}
	return siteURL
}
