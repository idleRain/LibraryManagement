package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// InitRedis 初始化 Redis 连接
func InitRedis(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect Redis: %w", err)
	}

	redisClient = client
	log.Println("✅ Redis connected successfully")

	return client, nil
}

// GetRedis 获取 Redis 客户端
func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() {
	if redisClient != nil {
		redisClient.Close()
	}
}

// ============================================
// 缓存操作
// ============================================

// CacheSet 设置缓存
func CacheSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if redisClient == nil {
		return nil // Redis 未启用，静默忽略
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, data, expiration).Err()
}

// CacheGet 获取缓存
func CacheGet(ctx context.Context, key string, dest interface{}) error {
	if redisClient == nil {
		return redis.Nil
	}

	data, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// CacheDel 删除缓存
func CacheDel(ctx context.Context, keys ...string) error {
	if redisClient == nil {
		return nil
	}

	return redisClient.Del(ctx, keys...).Err()
}

// CacheExists 检查缓存是否存在
func CacheExists(ctx context.Context, key string) (bool, error) {
	if redisClient == nil {
		return false, nil
	}

	n, err := redisClient.Exists(ctx, key).Result()
	return n > 0, err
}

// ============================================
// JWT 黑名单操作
// ============================================

const (
	JWTBlacklistPrefix = "jwt:blacklist:"
	JWTBlacklistTTL    = 24 * time.Hour
)

// JWTBlacklistAdd 将 JWT 加入黑名单
func JWTBlacklistAdd(ctx context.Context, token string, expiration time.Duration) error {
	if redisClient == nil {
		return nil
	}

	key := JWTBlacklistPrefix + token
	return redisClient.Set(ctx, key, "1", expiration).Err()
}

// JWTBlacklistExists 检查 JWT 是否在黑名单中
func JWTBlacklistExists(ctx context.Context, token string) (bool, error) {
	if redisClient == nil {
		return false, nil
	}

	key := JWTBlacklistPrefix + token
	return redisClient.Exists(ctx, key).Result()
}

// ============================================
// 分布式锁
// ============================================

const (
	LockPrefix = "lock:"
	LockTTL    = 10 * time.Second
)

// DistributedLock 分布式锁
type DistributedLock struct {
	key   string
	value string
}

// AcquireLock 获取分布式锁
func AcquireLock(ctx context.Context, key string, ttl time.Duration) (*DistributedLock, error) {
	if redisClient == nil {
		return &DistributedLock{key: key, value: "local"}, nil
	}

	lockKey := LockPrefix + key
	value := fmt.Sprintf("%d", time.Now().UnixNano())

	// 使用 SET NX EX 原子操作
	ok, err := redisClient.SetNX(ctx, lockKey, value, ttl).Result()
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("failed to acquire lock: %s", key)
	}

	return &DistributedLock{key: lockKey, value: value}, nil
}

// Release 释放锁
func (l *DistributedLock) Release(ctx context.Context) error {
	if redisClient == nil {
		return nil
	}

	// 使用 Lua 脚本确保原子性
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`

	return redisClient.Eval(ctx, script, []string{l.key}, l.value).Err()
}

// ============================================
// 计数器操作
// ============================================

// IncrCounter 增加计数器
func IncrCounter(ctx context.Context, key string) (int64, error) {
	if redisClient == nil {
		return 0, nil
	}

	return redisClient.Incr(ctx, key).Result()
}

// GetCounter 获取计数器值
func GetCounter(ctx context.Context, key string) (int64, error) {
	if redisClient == nil {
		return 0, nil
	}

	return redisClient.Get(ctx, key).Int64()
}

// ============================================
// 限流器
// ============================================

// RateLimiter 限流器
type RateLimiter struct {
	key   string
	limit int
	window time.Duration
}

// NewRateLimiter 创建限流器
func NewRateLimiter(key string, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		key:    key,
		limit:  limit,
		window: window,
	}
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow(ctx context.Context) (bool, error) {
	if redisClient == nil {
		return true, nil
	}

	key := "rate:" + r.key

	// 使用滑动窗口算法
	now := time.Now().UnixMilli()
	windowStart := now - r.window.Milliseconds()

	// Lua 脚本实现滑动窗口
	script := `
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local windowStart = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		local window = tonumber(ARGV[4])

		-- 移除过期的请求
		redis.call("ZREMRANGEBYSCORE", key, 0, windowStart)

		-- 获取当前窗口内的请求数
		local count = redis.call("ZCARD", key)

		if count < limit then
			-- 添加当前请求
			redis.call("ZADD", key, now, now)
			redis.call("PEXPIRE", key, window)
			return 1
		else
			return 0
		end
	`

	result, err := redisClient.Eval(ctx, script, []string{key},
		now, windowStart, r.limit, r.window.Milliseconds()).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
