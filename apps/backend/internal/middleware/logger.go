package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"gorm.io/gorm"
)

// OperationLogger 操作日志中间件
func OperationLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录 POST 请求
		if c.Request.Method != "POST" {
			c.Next()
			return
		}

		// 排除不需要记录的路径
		path := c.Request.URL.Path
		if strings.Contains(path, "/login") ||
			strings.Contains(path, "/register") ||
			strings.Contains(path, "/health") ||
			strings.Contains(path, "/profile") ||
			strings.Contains(path, "/list") ||
			strings.Contains(path, "/stats") {
			c.Next()
			return
		}

		// 记录开始时间
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 使用响应捕获器
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 获取用户信息
		var userID uint
		var username string
		if v, exists := c.Get("user_id"); exists {
			userID = v.(uint)
		}
		if v, exists := c.Get("username"); exists {
			username = v.(string)
		}

		// 解析模块和操作
		module, action := parseModuleAndAction(path)

		// 创建日志记录
		log := models.OperationLog{
			UserID:     userID,
			Username:   username,
			Module:     module,
			Action:     action,
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Content:    string(requestBody),
			CreatedAt:  start,
		}

		// 异步保存日志
		go func() {
			db.Create(&log)
		}()
	}
}

// bodyLogWriter 响应体捕获器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// parseModuleAndAction 解析模块和操作
func parseModuleAndAction(path string) (module, action string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// /api/users/create -> module: users, action: create
	if len(parts) >= 3 {
		module = parts[1]
		action = parts[2]
	}

	// 映射操作名称
	actionMap := map[string]string{
		"create":         "创建",
		"update":         "更新",
		"delete":         "删除",
		"assign-roles":   "分配角色",
		"change-password": "修改密码",
		"reset-password": "重置密码",
		"in":             "入库",
		"out":            "出库",
		"borrow":         "借阅",
		"return":         "归还",
		"renew":          "续借",
		"pay-fine":       "支付罚款",
		"cancel":         "取消",
		"add":            "添加",
		"remove":         "移除",
		"clear":          "清空",
	}

	if v, ok := actionMap[action]; ok {
		action = v
	}

	return module, action
}

// PrettyJSON 格式化 JSON
func PrettyJSON(data []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		return string(data)
	}
	return buf.String()
}
