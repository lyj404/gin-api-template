package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	HttpPort       string   `yaml:"HttpPort"`
	Mode           string   `yaml:"Mode"`
	AllowedOrigins []string `yaml:"AllowedOrigins"`
	RateLimit      int      `yaml:"RateLimit"`
}

type DatabaseConfig struct {
	Type            string `yaml:"Type"`
	Host            string `yaml:"Host"`
	Port            string `yaml:"Port"`
	User            string `yaml:"User"`
	Password        string `yaml:"Password"`
	Name            string `yaml:"Name"`
	MaxOpenConns    int    `yaml:"MaxOpenConns"`
	MaxIdleConns    int    `yaml:"MaxIdleConns"`
	ConnMaxLifetime int    `yaml:"ConnMaxLifetime"`
	ConnMaxIdleTime int    `yaml:"ConnMaxIdleTime"`
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

type LogConfig struct {
	Level      string `yaml:"Level"`
	Format     string `yaml:"Format"`
	Filename   string `yaml:"Filename"`
	MaxSize    int    `yaml:"MaxSize"`    // MB
	MaxBackups int    `yaml:"MaxBackups"` // number of backups
	MaxAge     int    `yaml:"MaxAge"`     // days
	Compress   bool   `yaml:"Compress"`   // compress backups
}

type SessionConfig struct {
	SessionSecret string `yaml:"SessionSecret"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Timeout  TimeoutConfig  `yaml:"timeout"`
	Token    TokenConfig    `yaml:"token"`
	Password PasswordConfig `yaml:"password"`
	Session  SessionConfig  `yaml:"session"`
	Log      LogConfig      `yaml:"log"`
}

var (
	CfgServer   ServerConfig
	CfgDatabase DatabaseConfig
	CfgRedis    RedisConfig
	CfgTimeout  TimeoutConfig
	CfgToken    TokenConfig
	CfgPassword PasswordConfig
	CfgSession  SessionConfig
	CfgLog      LogConfig
)

func InitConfig() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，将使用配置文件默认值")
	}

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("无法获取当前工作目录: %v", err)
	}
	log.Printf("INFO: 当前工作目录: %s", wd)

	// 尝试在当前工作目录下查找配置文件
	configPath := filepath.Join(wd, "config", "config.yml")

	// 如果找不到配置文件，尝试向上查找两级目录
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 尝试项目根目录 (向上一级)
		configPath = filepath.Join(wd, "..", "config", "config.yml")
	}

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

	// 从环境变量覆盖配置（优先级高于 config.yml）
	// 这样可以将敏感信息放在 .env 文件中，不提交到版本控制

	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		cfg.Server.HttpPort = httpPort
	}
	if mode := os.Getenv("MODE"); mode != "" {
		cfg.Server.Mode = mode
	}
	if rateLimit := os.Getenv("RATE_LIMIT"); rateLimit != "" {
		if val, err := strconv.Atoi(rateLimit); err == nil {
			cfg.Server.RateLimit = val
		}
	}

	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		cfg.Database.Type = dbType
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		cfg.Database.Port = dbPort
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.Name = dbName
	}
	if maxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConns != "" {
		if val, err := strconv.Atoi(maxOpenConns); err == nil {
			cfg.Database.MaxOpenConns = val
		}
	}
	if maxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConns != "" {
		if val, err := strconv.Atoi(maxIdleConns); err == nil {
			cfg.Database.MaxIdleConns = val
		}
	}
	if connMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); connMaxLifetime != "" {
		if val, err := strconv.Atoi(connMaxLifetime); err == nil {
			cfg.Database.ConnMaxLifetime = val
		}
	}
	if connMaxIdleTime := os.Getenv("DB_CONN_MAX_IDLE_TIME"); connMaxIdleTime != "" {
		if val, err := strconv.Atoi(connMaxIdleTime); err == nil {
			cfg.Database.ConnMaxIdleTime = val
		}
	}

	if redisEnabled := os.Getenv("REDIS_ENABLED"); redisEnabled != "" {
		if val, err := strconv.ParseBool(redisEnabled); err == nil {
			cfg.Redis.Enabled = val
		}
	}
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		cfg.Redis.Host = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		cfg.Redis.Port = redisPort
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		cfg.Redis.Password = redisPassword
	}
	if redisDatabase := os.Getenv("REDIS_DATABASE"); redisDatabase != "" {
		if val, err := strconv.Atoi(redisDatabase); err == nil {
			cfg.Redis.Database = val
		}
	}

	if contextTimeout := os.Getenv("CONTEXT_TIMEOUT"); contextTimeout != "" {
		if val, err := strconv.Atoi(contextTimeout); err == nil {
			cfg.Timeout.ContextTimeout = val
		}
	}

	if accessTokenExpiryHour := os.Getenv("ACCESS_TOKEN_EXPIRY_HOUR"); accessTokenExpiryHour != "" {
		if val, err := strconv.Atoi(accessTokenExpiryHour); err == nil {
			cfg.Token.AccessTokenExpiryHour = val
		}
	}
	if refreshTokenExpiryHour := os.Getenv("REFRESH_TOKEN_EXPIRY_HOUR"); refreshTokenExpiryHour != "" {
		if val, err := strconv.Atoi(refreshTokenExpiryHour); err == nil {
			cfg.Token.RefreshTokenExpiryHour = val
		}
	}
	if accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET"); accessTokenSecret != "" {
		cfg.Token.AccessTokenSecret = accessTokenSecret
	}
	if refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET"); refreshTokenSecret != "" {
		cfg.Token.RefreshTokenSecret = refreshTokenSecret
	}
	if tokenPrefix := os.Getenv("TOKEN_PREFIX"); tokenPrefix != "" {
		cfg.Token.TokenPrefix = tokenPrefix
	}

	if sessionSecret := os.Getenv("SESSION_SECRET"); sessionSecret != "" {
		cfg.Session.SessionSecret = sessionSecret
	}

	if saltPrefix := os.Getenv("SALT_PREFIX"); saltPrefix != "" {
		cfg.Password.SaltPrefix = saltPrefix
	}
	if saltSuffix := os.Getenv("SALT_SUFFIX"); saltSuffix != "" {
		cfg.Password.SaltSuffix = saltSuffix
	}
	if cost := os.Getenv("PASSWORD_COST"); cost != "" {
		if val, err := strconv.Atoi(cost); err == nil {
			cfg.Password.Cost = val
		}
	}

	CfgServer = cfg.Server
	CfgDatabase = cfg.Database
	CfgRedis = cfg.Redis
	CfgTimeout = cfg.Timeout
	CfgToken = cfg.Token
	CfgPassword = cfg.Password
	CfgSession = cfg.Session
	CfgLog = cfg.Log
}
