package conf

var defaultDB = databaseConf{
	Charset:      "utf8mb4",
	MaxIdleConns: 10,
	MaxOpenConns: 100,
}

var defaultApp = appConf{
	Mode:          "debug",
	Name:          "go-starter",
	SessionSecret: "go-starter",
}

var defaultServer = serverConf{
	Port:   "8080",
	Prefix: "/",
	Cores: coresConf{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"*"},
		SameSite:         "Lax",
		Secure:           false,
	},
}

var defaultRedis = redisConf{
	Network: "tcp",
	Port:    6379,
	DB:      0,
}

var defaultLog = logConf{
	Level:      "info",
	Path:       "logs",
	FileName:   "app.log",
	MaxSize:    100,
	MaxAge:     30,
	MaxBackups: 30,
	Compress:   false,
}

var defaultEmail = emailConf{
	Keepalive: 60,
}

func setDefault() {
	setDefaultDB()
	setDefaultApp()
	setDefaultServer()
	setDefaultRedis()
	setDefaultLog()
	setDefaultEmail()
}

func setDefaultDB() {
	cfg.Database = defaultDB
}

func setDefaultApp() {
	cfg.App = defaultApp
}

func setDefaultServer() {
	cfg.Server = defaultServer
}

func setDefaultRedis() {
	cfg.Redis = defaultRedis
}

func setDefaultLog() {
	cfg.Log = defaultLog
}

func setDefaultEmail() {
	cfg.Email = defaultEmail
}
