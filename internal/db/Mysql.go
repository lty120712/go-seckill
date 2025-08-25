package db

import (
	"database/sql"
	"fmt"
	"go-chat/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Mysql *gorm.DB // 全局变量，供其他包使用

// Init 初始化数据库连接
func InitMysql() {
	// 数据库连接字符串（DSN）
	// 根据配置文件构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.AppConfig.Database.Username,
		configs.AppConfig.Database.Password,
		configs.AppConfig.Database.Host,
		configs.AppConfig.Database.Port,
		configs.AppConfig.Database.Dbname,
	)
	// 连接数据库
	var err error
	Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出目标（终端）
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // 日志级别（Info = 打印 SQL）
				Colorful:      true,        // 彩色打印
			},
		),
	})

	// 错误检查
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 设置连接池
	var sqlDB *sql.DB
	sqlDB, err = Mysql.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
		return
	}

	// 最小连接数
	sqlDB.SetMaxIdleConns(2)
	// 最大连接数
	sqlDB.SetMaxOpenConns(10)
	// 单条连接最长存活时间，到期后会被关闭并新建
	sqlDB.SetConnMaxLifetime(time.Hour)
	//打印数据库初始化 信息
	log.Println("Database initialized successfully")
}

// 获取 GORM DB 实例，优先使用事务，如果没有事务则使用默认的 db.Mysql
func GetGormDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 {
		return tx[0] // 使用事务对象
	}
	return Mysql // 使用默认数据库连接
}
