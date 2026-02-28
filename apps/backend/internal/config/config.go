package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server  ServerConfig
	MySQL   MySQLConfig
	MongoDB MongoDBConfig
	Redis   RedisConfig
	JWT     JWTConfig
}

type ServerConfig struct {
	Port int
}

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Enabled  bool
}

type JWTConfig struct {
	Secret string
	Expire int
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnvInt("MYSQL_PORT", 3306),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", ""),
			Database: getEnv("MYSQL_DATABASE", "library_system"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "library_borrow"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			Enabled:  getEnvBool("REDIS_ENABLED", true),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "library-system-secret-key"),
			Expire: getEnvInt("JWT_EXPIRE", 86400),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
