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
	App      App               `yaml:"app"`
	Server   map[string]Server `yaml:"server"`
	SMTP     map[string]SMTP   `yaml:"smtp"`
	Template map[string]string `yaml:"template"`
	DB       map[string]DB     `yaml:"db"`
	Cache    map[string]Cache  `yaml:"cache"`
	Log      map[string]Log    `yaml:"log"`
}

type App struct {
	Env      string `yaml:"env"`
	Name     string `yaml:"name"`
	Domain   string `yaml:"domain"`
	Timezone string `yaml:"timezone"`
}

type Server struct {
	Host       string            `yaml:"host"`
	Port       int               `yaml:"port"`
	Address    string            `yaml:"address"`
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
	Timezone          string `yaml:"timezone"`
	Sslmode           string `yaml:"sslmode"`
	ConnectionTimeout int    `yaml:"connectionTimeout"`
}

type DBCredential struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Log struct {
	Disabled bool       `yaml:"disabled"`
	Console  LogConsole `yaml:"console"`
	File     LogFile    `yaml:"file"`
}

type LogConsole struct {
	Disabled bool `yaml:"disabled"`
	Level    int  `yaml:"level"`
}

type LogFile struct {
	Disabled bool        `yaml:"disabled"`
	Level    int         `yaml:"level"`
	Rotation LogRotation `yaml:"rotation"`
}

type LogRotation struct {
	BasePath  string `yaml:"basePath"`
	Filename  string `yaml:"filename"`
	MaxBackup int    `yaml:"maxBackup"`
	MaxSize   int    `yaml:"maxSize"`
	MaxAge    int    `yaml:"maxAge"`
	LocalTime bool   `yaml:"localTime"`
	Compress  bool   `yaml:"compress"`
}
