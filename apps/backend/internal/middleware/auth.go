package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/config"
	"github.com/library-system/backend/pkg/database"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string, cfg *config.JWTConfig) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.Expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "library-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string, cfg *config.JWTConfig) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// JWTAuth JWT 认证中间件
func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "Token 格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 检查黑名单
		inBlacklist, err := database.JWTBlacklistExists(c.Request.Context(), tokenString)
		if err != nil {
			// Redis 错误，继续验证 token
			inBlacklist = false
		}
		if inBlacklist {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "Token 已失效，请重新登录",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString, cfg)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 存储用户信息
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token", tokenString)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（不强制要求登录）
func OptionalAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := ParseToken(parts[1], cfg)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
			}
		}

		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s %d %v",
			start.Format("2006/01/02 15:04:05"),
			method,
			path,
			statusCode,
			latency,
		)
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// RateLimit 限流中间件
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 IP 作为限流 key
		key := c.ClientIP()
		limiter := database.NewRateLimiter(key, limit, window)

		allowed, err := limiter.Allow(c.Request.Context())
		if err != nil {
			// Redis 错误，允许请求通过
			c.Next()
			return
		}

		if !allowed {
			c.JSON(429, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
