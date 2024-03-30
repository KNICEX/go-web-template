package models

func migration() {
	DB.AutoMigrate(&User{}, &Setting{})

	addDefaultSettings()
}

func addDefaultSettings() {
	for _, setting := range defaultSettings {
		DB.Create(&setting)
	}
}
