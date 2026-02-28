package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DashboardController struct {
	mysqlDB *gorm.DB
	mongoDB *mongo.Database
}

func NewDashboardController(mysqlDB *gorm.DB, mongoDB *mongo.Database) *DashboardController {
	return &DashboardController{mysqlDB: mysqlDB, mongoDB: mongoDB}
}

// GetOverview 获取仪表盘概览数据
func (dc *DashboardController) GetOverview(c *gin.Context) {
	// 图书统计
	var bookStats struct {
		Total    int64
		Active   int64
		Inactive int64
	}
	dc.mysqlDB.Model(&models.Book{}).Count(&bookStats.Total)
	dc.mysqlDB.Model(&models.Book{}).Where("status = ?", 1).Count(&bookStats.Active)
	dc.mysqlDB.Model(&models.Book{}).Where("status = ?", 0).Count(&bookStats.Inactive)

	// 用户统计
	var userStats struct {
		Total   int64
		Active  int64
		NewToday int64
	}
	dc.mysqlDB.Model(&models.User{}).Count(&userStats.Total)
	dc.mysqlDB.Model(&models.User{}).Where("status = ?", 1).Count(&userStats.Active)
	today := time.Now().Format("2006-01-02")
	dc.mysqlDB.Model(&models.User{}).Where("created_at >= ?", today).Count(&userStats.NewToday)

	// 库存统计
	var stockStats struct {
		Total      int64
		LowStock   int64
		OutOfStock int64
	}
	dc.mysqlDB.Model(&models.Stock{}).Count(&stockStats.Total)
	dc.mysqlDB.Model(&models.Stock{}).Where("available_quantity <= warning_threshold AND available_quantity > 0").Count(&stockStats.LowStock)
	dc.mysqlDB.Model(&models.Stock{}).Where("available_quantity = 0").Count(&stockStats.OutOfStock)

	// 借阅统计
	borrowCollection := dc.mongoDB.Collection("borrow_records")
	borrowedCount, _ := borrowCollection.CountDocuments(context.Background(), bson.M{"status": "borrowed"})
	overdueCount, _ := borrowCollection.CountDocuments(context.Background(), bson.M{"status": "overdue"})

	// 今日销售
	var todaySales struct {
		Orders int64
		Amount float64
	}
	dc.mysqlDB.Model(&models.SaleOrder{}).
		Where("status = 'paid' AND created_at >= ?", today+" 00:00:00").
		Count(&todaySales.Orders)
	dc.mysqlDB.Model(&models.SaleOrder{}).
		Where("status = 'paid' AND created_at >= ?", today+" 00:00:00").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&todaySales.Amount)

	utils.Success(c, gin.H{
		"books": gin.H{
			"total":    bookStats.Total,
			"active":   bookStats.Active,
			"inactive": bookStats.Inactive,
		},
		"users": gin.H{
			"total":     userStats.Total,
			"active":    userStats.Active,
			"new_today": userStats.NewToday,
		},
		"stocks": gin.H{
			"total":       stockStats.Total,
			"low_stock":   stockStats.LowStock,
			"out_of_stock": stockStats.OutOfStock,
		},
		"borrows": gin.H{
			"borrowed": borrowedCount,
			"overdue":  overdueCount,
		},
		"sales": gin.H{
			"today_orders": todaySales.Orders,
			"today_amount": todaySales.Amount,
		},
	})
}

// GetBorrowTrend 获取借阅趋势（最近7天/30天）
func (dc *DashboardController) GetBorrowTrend(c *gin.Context) {
	var req struct {
		Days int `json:"days"`
	}
	c.ShouldBindJSON(&req)
	if req.Days <= 0 {
		req.Days = 7
	}

	collection := dc.mongoDB.Collection("borrow_records")

	// 计算日期范围
	startDate := time.Now().AddDate(0, 0, -req.Days)

	// 按日期聚合
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"borrow_date": bson.M{"$gte": startDate},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$borrow_date",
					},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	defer cursor.Close(context.Background())

	var results []struct {
		Date  string `bson:"_id"`
		Count int    `bson:"count"`
	}
	cursor.All(context.Background(), &results)

	// 填充缺失的日期
	trend := make([]map[string]interface{}, 0)
	for i := 0; i < req.Days; i++ {
		date := time.Now().AddDate(0, 0, -req.Days+i+1).Format("2006-01-02")
		count := 0
		for _, r := range results {
			if r.Date == date {
				count = r.Count
				break
			}
		}
		trend = append(trend, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}

	utils.Success(c, trend)
}

// GetSalesTrend 获取销售趋势
func (dc *DashboardController) GetSalesTrend(c *gin.Context) {
	var req struct {
		Days int `json:"days"`
	}
	c.ShouldBindJSON(&req)
	if req.Days <= 0 {
		req.Days = 7
	}

	type DayStat struct {
		Date   string
		Orders int64
		Amount float64
	}

	var stats []DayStat

	for i := 0; i < req.Days; i++ {
		date := time.Now().AddDate(0, 0, -req.Days+i+1)
		dateStr := date.Format("2006-01-02")

		var orders int64
		var amount float64

		dc.mysqlDB.Model(&models.SaleOrder{}).
			Where("status = 'paid' AND DATE(created_at) = ?", dateStr).
			Count(&orders)
		dc.mysqlDB.Model(&models.SaleOrder{}).
			Where("status = 'paid' AND DATE(created_at) = ?", dateStr).
			Select("COALESCE(SUM(total_amount), 0)").Scan(&amount)

		stats = append(stats, DayStat{
			Date:   dateStr,
			Orders: orders,
			Amount: amount,
		})
	}

	utils.Success(c, stats)
}

// GetCategoryStats 获取图书分类统计
func (dc *DashboardController) GetCategoryStats(c *gin.Context) {
	type CategoryStat struct {
		Category string
		Count    int64
	}

	var stats []CategoryStat
	dc.mysqlDB.Model(&models.Book{}).
		Select("category, count(*) as count").
		Where("category != ''").
		Group("category").
		Order("count DESC").
		Limit(10).
		Find(&stats)

	utils.Success(c, stats)
}

// GetTopBorrowedBooks 获取热门借阅图书
func (dc *DashboardController) GetTopBorrowedBooks(c *gin.Context) {
	var req struct {
		Limit int `json:"limit"`
	}
	c.ShouldBindJSON(&req)
	if req.Limit <= 0 {
		req.Limit = 10
	}

	collection := dc.mongoDB.Collection("borrow_records")

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$book_id",
				"title": bson.M{"$first": "$book_title"},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
		{
			"$limit": req.Limit,
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	defer cursor.Close(context.Background())

	var results []struct {
		BookID int    `bson:"_id"`
		Title  string `bson:"title"`
		Count  int    `bson:"count"`
	}
	cursor.All(context.Background(), &results)

	utils.Success(c, results)
}

// GetTopSoldBooks 获取热销图书
func (dc *DashboardController) GetTopSoldBooks(c *gin.Context) {
	var req struct {
		Limit int `json:"limit"`
	}
	c.ShouldBindJSON(&req)
	if req.Limit <= 0 {
		req.Limit = 10
	}

	type BookSale struct {
		BookID   int
		Title    string
		Quantity int64
		Amount   float64
	}

	var results []BookSale
	dc.mysqlDB.Table("sale_order_items").
		Select("book_id, books.title, SUM(sale_order_items.quantity) as quantity, SUM(sale_order_items.total_price) as amount").
		Joins("JOIN books ON books.id = sale_order_items.book_id").
		Joins("JOIN sale_orders ON sale_orders.id = sale_order_items.order_id").
		Where("sale_orders.status = 'paid'").
		Group("book_id, books.title").
		Order("quantity DESC").
		Limit(req.Limit).
		Find(&results)

	utils.Success(c, results)
}

// GetRecentActivities 获取最近活动
func (dc *DashboardController) GetRecentActivities(c *gin.Context) {
	var req struct {
		Limit int `json:"limit"`
	}
	c.ShouldBindJSON(&req)
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 最近借阅
	var recentBorrows []models.OperationLog
	dc.mysqlDB.Where("module = ?", "borrows").
		Order("created_at DESC").
		Limit(req.Limit/2).
		Find(&recentBorrows)

	// 最近销售
	var recentSales []models.OperationLog
	dc.mysqlDB.Where("module = ?", "sales").
		Order("created_at DESC").
		Limit(req.Limit/2).
		Find(&recentSales)

	utils.Success(c, gin.H{
		"borrows": recentBorrows,
		"sales":   recentSales,
	})
}

// GetAlerts 获取预警信息
func (dc *DashboardController) GetAlerts(c *gin.Context) {
	alerts := make([]map[string]interface{}, 0)

	// 库存预警
	var lowStocks []models.Stock
	dc.mysqlDB.Preload("Book").
		Where("available_quantity <= warning_threshold").
		Find(&lowStocks)

	for _, s := range lowStocks {
		alertType := "warning"
		if s.AvailableQty == 0 {
			alertType = "error"
		}
		alerts = append(alerts, map[string]interface{}{
			"type":    alertType,
			"module":  "stock",
			"title":   "库存预警",
			"message": s.Book.Title + " 库存不足",
			"detail":  s.AvailableQty,
		})
	}

	// 逾期借阅
	borrowCollection := dc.mongoDB.Collection("borrow_records")
	overdueCount, _ := borrowCollection.CountDocuments(context.Background(), bson.M{
		"status":    "borrowed",
		"due_date": bson.M{"$lt": time.Now()},
	})

	if overdueCount > 0 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "warning",
			"module":  "borrow",
			"title":   "逾期提醒",
			"message": "有图书逾期未还",
			"detail":  overdueCount,
		})
	}

	utils.Success(c, alerts)
}

// Search 全局搜索
func (dc *DashboardController) Search(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword" binding:"required"`
		Limit   int    `json:"limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入搜索关键词")
		return
	}
	if req.Limit <= 0 {
		req.Limit = 5
	}

	results := make(map[string]interface{})

	// 搜索图书
	var books []models.Book
	dc.mysqlDB.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%").
		Limit(req.Limit).
		Find(&books)
	results["books"] = books

	// 搜索用户
	var users []models.User
	dc.mysqlDB.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ?",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%").
		Limit(req.Limit).
		Find(&users)
	results["users"] = users

	// 搜索供应商
	var suppliers []models.Supplier
	dc.mysqlDB.Where("name LIKE ? OR contact LIKE ?",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%").
		Limit(req.Limit).
		Find(&suppliers)
	results["suppliers"] = suppliers

	// 搜索借阅记录
	borrowCollection := dc.mongoDB.Collection("borrow_records")
	cursor, _ := borrowCollection.Find(context.Background(), bson.M{
		"$or": []bson.M{
			{"book_title": primitive.Regex{Pattern: req.Keyword, Options: "i"}},
			{"user_name": primitive.Regex{Pattern: req.Keyword, Options: "i"}},
		},
	})
	defer cursor.Close(context.Background())

	var borrows []bson.M
	cursor.All(context.Background(), &borrows)
	results["borrows"] = borrows

	utils.Success(c, results)
}
