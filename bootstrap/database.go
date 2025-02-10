package bootstrap

import (
	"fmt"
	"gin-api-template/config"
	"gin-api-template/domain"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDataBase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUser, config.DbPassWord, config.DbHost, config.DbPort, config.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	// 设置自动迁移
	db.AutoMigrate(&domain.User{})

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

func CloseMySQLConnection(db *gorm.DB) {
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
