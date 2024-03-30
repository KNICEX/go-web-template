package auth

import "go-web-template/models"

func Init() {
	HMACSecret := models.GetSettingByName(models.SecretKey)
	if HMACSecret == "" {
		panic("secret key is missing")
	}
	General = HMACAuth{SecretKey: []byte(HMACSecret)}
}
