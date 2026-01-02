package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDataBase() *gorm.DB {
	var dsn string
	var db *gorm.DB
	var err error

	switch config.CfgDatabase.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.CfgDatabase.User, config.CfgDatabase.Password, config.CfgDatabase.Host, config.CfgDatabase.Port, config.CfgDatabase.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.CfgDatabase.Host, config.CfgDatabase.User, config.CfgDatabase.Password, config.CfgDatabase.Name, config.CfgDatabase.Port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	default:
		log.Fatal("不支持的数据库类型:", config.CfgDatabase.Type)
	}

	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 设置自动迁移
	db.AutoMigrate(&entity.User{})

	// 设置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("配置数据库失败:", err)
	}

	// 从配置读取连接池参数，如果未配置则使用默认值
	maxOpenConns := config.CfgDatabase.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}
	maxIdleConns := config.CfgDatabase.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	connMaxLifetime := config.CfgDatabase.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 3600
	}
	connMaxIdleTime := config.CfgDatabase.ConnMaxIdleTime
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = 600
	}

	// 设置连接池参数
	// MaxOpenConns: 最大打开连接数
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// MaxIdleConns: 最大空闲连接数
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// ConnMaxLifetime: 连接最大存活时间（秒），超过后会被关闭
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	// ConnMaxIdleTime: 连接最大空闲时间（秒），超过后会被关闭
	sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Second)

	return db
}

func CloseDataBaseConnection(db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("获取底层数据库对象失败:", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Println("关闭数据库连接失败:", err)
	} else {
		log.Println("数据库连接已关闭")
	}
}
