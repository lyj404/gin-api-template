package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	HttpPort string `yaml:"HttpPort"`
	Mode     string `yaml:"Mode"`
}

type DatabaseConfig struct {
	Type     string `yaml:"Type"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Name     string `yaml:"Name"`
}

type RedisConfig struct {
	Enabled  bool   `yaml:"Enabled"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Password string `yaml:"Password"`
	Database int    `yaml:"DataBase"`
}

type TimeoutConfig struct {
	ContextTimeout int `yaml:"ContextTimeOut"`
}

type TokenConfig struct {
	AccessTokenExpiryHour  int    `yaml:"AccessTokenExpiryHour"`
	RefreshTokenExpiryHour int    `yaml:"RefreshTokenExpiryHour"`
	AccessTokenSecret      string `yaml:"AccessTokenSecret"`
	RefreshTokenSecret     string `yaml:"RefreshTokenSecret"`
	TokenPrefix            string `yaml:"TokenPrefix"`
}

type PasswordConfig struct {
	SaltPrefix string `yaml:"SaltPrefix"`
	SaltSuffix string `yaml:"SaltSuffix"`
	Cost       int    `yaml:"Cost"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Timeout  TimeoutConfig  `yaml:"timeout"`
	Token    TokenConfig    `yaml:"token"`
	Password PasswordConfig `yaml:"password"`
}

var (
	CfgServer   ServerConfig
	CfgDatabase DatabaseConfig
	CfgRedis    RedisConfig
	CfgTimeout  TimeoutConfig
	CfgToken    TokenConfig
	CfgPassword PasswordConfig
)

func InitConfig() {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录时出错: %v", err)
	}

	// 拼接配置文件路径
	configPath := filepath.Join(dir, "config", "config.yml")
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("配置文件读取错误，请检查文件路径:", err)
	}

	var cfg Config
	// 解析配置文件
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("配置文件解析错误:", err)
	}

	CfgServer = cfg.Server
	CfgDatabase = cfg.Database
	CfgRedis = cfg.Redis
	CfgTimeout = cfg.Timeout
	CfgToken = cfg.Token
	CfgPassword = cfg.Password
}
