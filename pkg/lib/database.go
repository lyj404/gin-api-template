package lib

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lyj404/gin-api-template/config"
	"github.com/lyj404/gin-api-template/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDataBase() *gorm.DB {
	db, err := openDatabase()
	if err != nil {
		// 数据库不存在时自动创建后重连
		if isDatabaseNotExistError(err) {
			log.Printf("数据库 %s 不存在，尝试自动创建...", config.CfgDatabase.Name)
			if createErr := createDatabase(); createErr != nil {
				log.Fatal("自动创建数据库失败:", createErr)
			}
			log.Printf("数据库 %s 创建成功", config.CfgDatabase.Name)
			db, err = openDatabase()
		}
		if err != nil {
			log.Fatal("连接数据库失败:", err)
		}
	}

	// 设置自动迁移
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Resource{},
		&entity.Role{},
		&entity.RoleResource{},
		&entity.OrgUnit{},
		&entity.RoleOrgScope{},
		&entity.UserRole{},
		&entity.OrgEntityBinding{},
		&entity.AuditLog{},
		&entity.Menu{},
		&entity.RoleMenu{},
		&entity.MenuResource{},
		&entity.SysDictionary{},
		&entity.SysDictionaryDetail{},
	); err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}

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

// openDatabase 根据配置打开目标数据库连接
func openDatabase() (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	switch config.CfgDatabase.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.CfgDatabase.User, config.CfgDatabase.Password, config.CfgDatabase.Host, config.CfgDatabase.Port, config.CfgDatabase.Name)
		return gorm.Open(mysql.Open(dsn), gormCfg)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.CfgDatabase.Host, config.CfgDatabase.User, config.CfgDatabase.Password, config.CfgDatabase.Name, config.CfgDatabase.Port)
		return gorm.Open(postgres.Open(dsn), gormCfg)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.CfgDatabase.Type)
	}
}

// createDatabase 连接到管理库并创建目标数据库
func createDatabase() error {
	var adminDB *gorm.DB
	var err error
	var createSQL string

	switch config.CfgDatabase.Type {
	case "mysql":
		// MySQL 连接时可不指定数据库
		adminDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&parseTime=True&loc=Local",
			config.CfgDatabase.User, config.CfgDatabase.Password, config.CfgDatabase.Host, config.CfgDatabase.Port)
		adminDB, err = gorm.Open(mysql.Open(adminDSN), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("连接 MySQL 服务失败: %w", err)
		}
		createSQL = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci",
			config.CfgDatabase.Name)
	case "postgres":
		// PostgreSQL 必须连接到一个已存在的库（默认 postgres）
		adminDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.CfgDatabase.Host, config.CfgDatabase.User, config.CfgDatabase.Password, config.CfgDatabase.Port)
		adminDB, err = gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("连接 postgres 管理库失败: %w", err)
		}
		// PostgreSQL 不支持 IF NOT EXISTS 在所有版本上，这里依赖外层先判断不存在
		createSQL = fmt.Sprintf(`CREATE DATABASE "%s"`, config.CfgDatabase.Name)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.CfgDatabase.Type)
	}

	defer func() {
		if sqlDB, dbErr := adminDB.DB(); dbErr == nil {
			_ = sqlDB.Close()
		}
	}()

	if err := adminDB.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("执行建库语句失败: %w", err)
	}
	return nil
}

// isDatabaseNotExistError 判断错误是否为目标数据库不存在
func isDatabaseNotExistError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// PostgreSQL: SQLSTATE 3D000 (invalid_catalog_name)
	// MySQL: Error 1049 (42000): Unknown database 'xxx'
	return strings.Contains(msg, "SQLSTATE 3D000") ||
		strings.Contains(msg, "Error 1049") ||
		strings.Contains(msg, "Unknown database")
}
