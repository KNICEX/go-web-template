package conf

type appConf struct {
	Mode          string     `mapstructure:"mode"`
	Name          string     `mapstructure:"name"`
	SessionSecret string     `mapstructure:"session_secret"`
	RemoteConf    remoteConf `mapstructure:"remote"`
	Debug         bool       `mapstructure:"debug"`
	DataCenterId  int64      `mapstructure:"data_center_id"`
	MachineId     int64      `mapstructure:"machine_id"`
}

type serverConf struct {
	Port   string    `mapstructure:"port"`
	Prefix string    `mapstructure:"prefix"`
	Cores  coresConf `mapstructure:"cores"`
}

// 远程配置
type remoteConf struct {
	Provider  string `mapstructure:"provider"`
	Addr      string `mapstructure:"addr"`
	Path      string `mapstructure:"path"`
	Type      string `mapstructure:"type"`
	SecretKey string `mapstructure:"secret_key"`
}

type coresConf struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	SameSite         string   `mapstructure:"same_site"`
	Secure           bool     `mapstructure:"secure"`
}

type databaseConf struct {
	Type     string `mapstructure:"type"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Charset  string `mapstructure:"charset"`
	// for sqlite
	DBFile       string `mapstructure:"db_file"`
	TablePrefix  string `mapstructure:"table_prefix"`
	UnixSocket   bool   `mapstructure:"unix_socket"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type redisConf struct {
	Network  string `mapstructure:"network"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type logConf struct {
	Level      string `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

type emailConf struct {
	Name      string `mapstructure:"name"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Keepalive int    `mapstructure:"keepalive"`
}

type conf struct {
	App      appConf      `mapstructure:"app"`
	Server   serverConf   `mapstructure:"server"`
	Database databaseConf `mapstructure:"database"`
	Redis    redisConf    `mapstructure:"redis"`
	Log      logConf      `mapstructure:"log"`
	Email    emailConf    `mapstructure:"email"`
}

func GetString(key string) string {
	return v.GetString(key)
}

func GetInt(key string) int {
	return v.GetInt(key)
}

func GetBool(key string) bool {
	return v.GetBool(key)
}

func GetFloat64(key string) float64 {
	return v.GetFloat64(key)
}

func GetStringSlice(key string) []string {
	return v.GetStringSlice(key)
}

func GetStringMap(key string) map[string]interface{} {
	return v.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return v.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return v.GetStringMapStringSlice(key)
}

func IsSet(key string) bool {
	return v.IsSet(key)
}

func Set(key string, value interface{}) {
	v.Set(key, value)
}

func SetDefault(key string, value interface{}) {
	v.SetDefault(key, value)
}
