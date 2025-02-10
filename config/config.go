package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	REnabled bool
	RHost     string
	RPort     string
	RPassWord string
	RDataBase int

	ContextTimeOut int

	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
	TokenPrefix            string
)

func init() {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录时出错: %v", err)
	}

	// 构造配置文件的完整路径
	// 从当前目录返回上一级，然后进入config文件夹
	configPath := filepath.Join(dir, "..", "config", "config.ini")
	file, err := ini.Load(configPath)
	if err != nil {
		log.Fatal("配置文件读取错误，请检查文件路径:", err)
	}

	LoadServer(file)
	LoadDataBase(file)
	LoadRedis(file)
	LoadTimeOut(file)
	LoadTokenConfig(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8181")
}

func LoadDataBase(file *ini.File) {
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("123456")
	DbName = file.Section("database").Key("DbName").MustString("study")
}

func LoadRedis(file *ini.File) {
	REnabled = file.Section("redis").Key("REnabled").MustBool(true)
	RHost = file.Section("redis").Key("RHost").MustString("114.115.208.175")
	RPort = file.Section("redis").Key("RPort").MustString("6379")
	RPassWord = file.Section("redis").Key("RPassWord").String()
	RDataBase = file.Section("redis").Key("RDataBase").MustInt(0)
}

func LoadTimeOut(file *ini.File) {
	ContextTimeOut = file.Section("timeout").Key("ContextTimeOut").MustInt(3)
}

func LoadTokenConfig(file *ini.File) {
	AccessTokenExpiryHour = file.Section("token").Key("AccessTokenExpiryHour").MustInt(2)
	RefreshTokenExpiryHour = file.Section("token").Key("AccessTokenExpiryHour").MustInt(168)
	AccessTokenSecret = file.Section("token").Key("AccessTokenExpiryHour").MustString("access_token_secret")
	RefreshTokenSecret = file.Section("token").Key("AccessTokenExpiryHour").MustString("refresh_token_secret")
	TokenPrefix = file.Section("token").Key("TokenPrefix").MustString("Bearer ")
}
