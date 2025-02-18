package config

import (
	"flag"

	"go.uber.org/config"
)

func ProvideConfig() (Cfg, error) {
	var (
		c       Cfg
		cfgPath string
	)

	flag.StringVar(&cfgPath, "cfg", "config.yaml", "load config path")
	flag.Parse()

	cfg, err := config.NewYAML(config.File(cfgPath))
	if err != nil {
		return c, err
	}

	if err := cfg.Get("").Populate(&c); err != nil {
		return c, err
	}

	return c, nil
}

type Cfg struct {
	App      App                 `yaml:"app"`
	Otel     Otel                `yaml:"otel"`
	Server   map[string]Server   `yaml:"server"`
	SMTP     map[string]SMTP     `yaml:"smtp"`
	Template map[string]string   `yaml:"template"`
	DB       map[string]DB       `yaml:"db"`
	Cache    map[string]Cache    `yaml:"cache"`
	Log      Log                 `yaml:"log"`
	Provider map[string]Provider `yaml:"provider"`
}

type App struct {
	Env      string `yaml:"env"`
	Project  string `yaml:"project"`
	Timezone string `yaml:"timezone"`
}

type Otel struct {
	Tracer bool       `yaml:"tracer"`
	Metric bool       `yaml:"metric"`
	Server OtelServer `yaml:"server"`
}

type OtelServer struct {
	GrpcHost string `yaml:"grpc.host"`
	GrpcPort int    `yaml:"grpc.port"`
}

type Server struct {
	Name       string            `yaml:"name"`
	Host       string            `yaml:"host"`
	Port       int               `yaml:"port"`
	Address    string            `yaml:"address"`
	Domain     string            `yaml:"domain"`
	Additional map[string]string `yaml:"additional"`
}

type SMTP struct {
	Host       string         `yaml:"host"`
	Port       int            `yaml:"port"`
	Credential SMTPCredential `yaml:"credential"`
}

type SMTPCredential struct {
	Name     string `yaml:"name"`
	Email    string `yaml:"email"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Cache struct {
	Disabled   bool            `yaml:"disabled"`
	DBName     string          `yaml:"dbname"`
	Port       int             `yaml:"port"`
	Address    string          `yaml:"address"`
	Credential CacheCredential `yaml:"credential"`
	Options    CacheOption     `yaml:"options"`
}

type CacheOption struct {
	DialTimeout  int `yaml:"dial.timeout"`
	ReadTimeout  int `yaml:"read.timeout"`
	WriteTimeout int `yaml:"write.timeout"`
}

type CacheCredential struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DB struct {
	Disabled   bool         `yaml:"disabled"`
	Driver     string       `yaml:"driver"`
	DBName     string       `yaml:"dbname"`
	Port       int          `yaml:"port"`
	Address    string       `yaml:"address"`
	Credential DBCredential `yaml:"credential"`
	Options    DBOptions    `yaml:"options"`
}

type DBOptions struct {
	Timezone              string `yaml:"timezone"`
	Sslmode               string `yaml:"sslmode"`
	ConnectionTimeout     int    `yaml:"connection.timeout"`
	MaxConnectionLifetime int    `yaml:"max.connection.lifetime"`
	MaxOpenConnection     int    `yaml:"max.open.connection"`
	MaxIdleConnection     int    `yaml:"max.idle.connection"`
}

type DBCredential struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Log struct {
	BasePath  string             `yaml:"base.path"`
	TrxClient []string           `yaml:"trx.client"`
	LogType   map[string]LogType `yaml:"log.type"`
}

type LogType struct {
	Disabled bool       `yaml:"disabled"`
	Notify   LogNotify  `yaml:"notify"`
	Console  LogConsole `yaml:"console"`
	File     LogFile    `yaml:"file"`
}

type LogConsole struct {
	Disabled bool `yaml:"disabled"`
	Level    int  `yaml:"level"`
}

type LogNotify struct {
	Enabled   bool `yaml:"enabled"`
	Debug     bool `yaml:"debug"`
	Retention int  `yaml:"retention"`
}

type LogFile struct {
	Disabled bool        `yaml:"disabled"`
	Level    int         `yaml:"level"`
	Rotation LogRotation `yaml:"rotation"`
}

type LogRotation struct {
	Filename  string `yaml:"filename"`
	MaxBackup int    `yaml:"max.backup"`
	MaxSize   int    `yaml:"max.size"`
	MaxAge    int    `yaml:"max.age"`
	LocalTime bool   `yaml:"local.time"`
	Compress  bool   `yaml:"compress"`
}
type Provider struct {
	BaseURL    string            `yaml:"base.url"`
	Additional map[string]string `yaml:"additional"`
}
