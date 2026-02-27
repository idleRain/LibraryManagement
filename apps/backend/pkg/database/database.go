package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/library-system/backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlDB  *gorm.DB
	mongoDB  *mongo.Database
	mongoCli *mongo.Client
)

// InitMySQL 初始化 MySQL 连接
func InitMySQL(cfg config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	mysqlDB = db
	log.Println("✅ MySQL connected successfully")

	return db, nil
}

// InitMongoDB 初始化 MongoDB 连接
func InitMongoDB(cfg config.MongoDBConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect MongoDB: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	mongoCli = client
	mongoDB = client.Database(cfg.Database)
	log.Println("✅ MongoDB connected successfully")

	return mongoDB, nil
}

// GetMySQL 获取 MySQL 连接
func GetMySQL() *gorm.DB {
	return mysqlDB
}

// GetMongoDB 获取 MongoDB 连接
func GetMongoDB() *mongo.Database {
	return mongoDB
}

// CloseMySQL 关闭 MySQL 连接
func CloseMySQL() {
	if mysqlDB != nil {
		sqlDB, _ := mysqlDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// CloseMongoDB 关闭 MongoDB 连接
func CloseMongoDB() {
	if mongoCli != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoCli.Disconnect(ctx)
	}
}
