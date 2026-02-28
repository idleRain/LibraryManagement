package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type LogController struct {
	db *gorm.DB
}

func NewLogController(db *gorm.DB) *LogController {
	return &LogController{db: db}
}

// GetOperationLogs 获取操作日志列表
func (lc *LogController) GetOperationLogs(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		UserID   int    `json:"user_id"`
		Module   string `json:"module"`
		Action   string `json:"action"`
		Username string `json:"username"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Page = 1
		req.PageSize = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	query := lc.db.Model(&models.OperationLog{})

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var logs []models.OperationLog
	query.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&logs)

	utils.PageSuccess(c, logs, total, req.Page, req.PageSize)
}

// GetLogModules 获取日志模块列表
func (lc *LogController) GetLogModules(c *gin.Context) {
	var modules []string
	lc.db.Model(&models.OperationLog{}).
		Distinct().
		Pluck("module", &modules)

	utils.Success(c, modules)
}

// GetLogActions 获取日志操作列表
func (lc *LogController) GetLogActions(c *gin.Context) {
	var actions []string
	lc.db.Model(&models.OperationLog{}).
		Distinct().
		Pluck("action", &actions)

	utils.Success(c, actions)
}

// GetUserLoginHistory 获取用户登录历史
func (lc *LogController) GetUserLoginHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var logs []models.OperationLog
	lc.db.Where("user_id = ? AND action = ?", userID, "登录").
		Order("created_at DESC").
		Limit(10).
		Find(&logs)

	utils.Success(c, logs)
}

// GetDashboardStats 获取仪表盘统计
func (lc *LogController) GetDashboardStats(c *gin.Context) {
	// 今日操作数
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	lc.db.Model(&models.OperationLog{}).
		Where("created_at >= ?", today+" 00:00:00").
		Count(&todayCount)

	// 本周操作数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format("2006-01-02")
	var weekCount int64
	lc.db.Model(&models.OperationLog{}).
		Where("created_at >= ?", weekStart+" 00:00:00").
		Count(&weekCount)

	// 本月操作数
	monthStart := time.Now().Format("2006-01") + "-01"
	var monthCount int64
	lc.db.Model(&models.OperationLog{}).
		Where("created_at >= ?", monthStart+" 00:00:00").
		Count(&monthCount)

	// 活跃用户数（最近7天）
	var activeUsers int64
	lc.db.Model(&models.OperationLog{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7).Format("2006-01-02")+" 00:00:00").
		Distinct("user_id").
		Count(&activeUsers)

	// 操作类型统计
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	lc.db.Model(&models.OperationLog{}).
		Select("action, count(*) as count").
		Where("created_at >= ?", today+" 00:00:00").
		Group("action").
		Order("count DESC").
		Limit(10).
		Find(&actionStats)

	utils.Success(c, gin.H{
		"today_count":    todayCount,
		"week_count":     weekCount,
		"month_count":    monthCount,
		"active_users":   activeUsers,
		"action_stats":   actionStats,
	})
}

// CleanOldLogs 清理旧日志
func (lc *LogController) CleanOldLogs(c *gin.Context) {
	var req struct {
		Days int `json:"days" binding:"required,min=30"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误，至少保留30天")
		return
	}

	cutoffDate := time.Now().AddDate(0, 0, -req.Days).Format("2006-01-02") + " 00:00:00"

	result := lc.db.Where("created_at < ?", cutoffDate).Delete(&models.OperationLog{})

	utils.SuccessWithMsg(c, "清理完成", gin.H{
		"deleted_count": result.RowsAffected,
	})
}
