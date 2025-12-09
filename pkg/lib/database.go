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
	// 设置连接池中最大闲置连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置数据库连接池最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可重用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)

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
