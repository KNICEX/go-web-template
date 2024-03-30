package models

import "go-web-template/pkg/utils"

var defaultSettings = []Setting{
	{Name: SiteURLKey, Value: "http://localhost:8080", Type: "basic"},
	{Name: SecretKey, Value: utils.RandomString(10), Type: "basic"},
	{Name: JWTSecret, Value: utils.RandomString(10), Type: "basic"},
}
