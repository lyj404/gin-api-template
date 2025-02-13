package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type ServerConfig struct {
	HttpPort string
}

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	Database int
}

type TimeoutConfig struct {
	ContextTimeout int
}

type TokenConfig struct {
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
	TokenPrefix            string
}

type PasswordConfig struct {
	SaltPrefix string
	SaltSuffix string
	Cost       int
}

// 全局配置变量
var (
	CfgServer   ServerConfig
	CfgDatabase DatabaseConfig
	CfgRedis    RedisConfig
	CfgTimeout  TimeoutConfig
	CfgToken    TokenConfig
	CfgPassword PasswordConfig
)

func init() {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录时出错: %v", err)
	}

	// 构造配置文件的完整路径
	configPath := filepath.Join(dir, "config", "config.ini")
	file, err := ini.Load(configPath)
	if err != nil {
		log.Fatal("配置文件读取错误，请检查文件路径:", err)
	}

	// 加载配置
	LoadConfig(file)
}

// LoadConfig 通用加载配置方法
func LoadConfig(file *ini.File) {
	LoadSection(file, "server", &CfgServer)
	LoadSection(file, "database", &CfgDatabase)
	LoadSection(file, "redis", &CfgRedis)
	LoadSection(file, "timeout", &CfgTimeout)
	LoadSection(file, "token", &CfgToken)
	LoadSection(file, "password", &CfgPassword)
}

// LoadSection 加载具体配置
func LoadSection(file *ini.File, section string, config interface{}) {
	err := file.Section(section).MapTo(config)
	if err != nil {
		log.Fatalf("加载 %s 配置失败: %v", section, err)
	}
}
