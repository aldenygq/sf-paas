package config

import "time"

type Config struct {
	Log    LogConfig  `mapstructure:"log"`
	Mysql  DataBase   `mapstructure:"mysql"`
	Util   UtilParam  `mapstructure:"util"`
	Base   BaseInfo   `mapstructure:"base"`
	Redis  RedisInfo  `mapstructure:"redis"`
	Server ServerInfo `mapstructure:"server"`
	Oncall OncallInfo `mapstructure:"oncall"`
}
type OncallInfo struct {
	Info string `yaml:"info"`
}
type UtilParam struct {
	UuapUrl   string `yaml:"uuapurl"`
	Uuaptoken string `yaml:"uuaptoken"`
	FerryUrl string `yaml:"ferryurl"`
	PermDomain string `yaml:"permdomain"`
	PermUrl string `yaml:"permurl"`
	SvcUrl string `yaml:"svcurl"`
}
type LogConfig struct {
	Loglevel  string        `yaml:"loglevel"`
	Logfile   string        `yaml:"logfile"`
	Logmaxage time.Duration `yaml:"logmaxage"`
}
type DataBase struct {
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
	Dbname          string `yaml:"dbname"`
	Pwd             string `yaml:"pwd"`
	Port            int    `yaml:"port"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxConnLifeTime int    `yaml:"maxConnLifeTime"`
	Type            string `yaml:"type"`
	Dbcharset       string `yaml:"dbcharset"`
}
type BaseInfo struct {
	Opgroups    []string `yaml:"opgroups"`
	ObjManager  []string `yaml:"objmanager"`
	TypePath    string   `yaml:"typepath"`
	SupportPath string   `yaml:"supportpath"`
	UCenterPath string   `yaml:"ucenterpath"`
	FbTypes     []string `yaml:"fbtypes"`
}
type ServerInfo struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type RedisInfo struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}
