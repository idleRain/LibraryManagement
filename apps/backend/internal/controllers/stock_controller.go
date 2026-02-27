package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type StockController struct {
	db *gorm.DB
}

func NewStockController(db *gorm.DB) *StockController {
	return &StockController{db: db}
}

// GetStocks 获取库存列表
func (sc *StockController) GetStocks(c *gin.Context) {
	var req struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"page_size"`
		BookTitle string `json:"book_title"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Page = 1
		req.PageSize = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var stocks []models.Stock
	var total int64

	query := sc.db.Model(&models.Stock{}).Preload("Book")
	if req.BookTitle != "" {
		query = query.Joins("JOIN books ON books.id = stocks.book_id").Where("books.title LIKE ?", "%"+req.BookTitle+"%")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("stocks.created_at DESC").Find(&stocks)

	utils.PageSuccess(c, stocks, total, req.Page, req.PageSize)
}

// GetStockDetail 获取库存详情
func (sc *StockController) GetStockDetail(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var stock models.Stock
	if err := sc.db.Preload("Book").First(&stock, req.ID).Error; err != nil {
		utils.NotFound(c, "库存不存在")
		return
	}

	utils.Success(c, stock)
}

// StockIn 入库
func (sc *StockController) StockIn(c *gin.Context) {
	var req struct {
		BookID   uint   `json:"book_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
		Reason   string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("user_id")

	tx := sc.db.Begin()

	var stock models.Stock
	if err := tx.Where("book_id = ?", req.BookID).First(&stock).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "库存记录不存在")
		return
	}

	// 更新库存
	stock.TotalQuantity += req.Quantity
	stock.AvailableQty += req.Quantity
	now := time.Now()
	stock.LastStockInAt = &now
	if err := tx.Save(&stock).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "更新库存失败")
		return
	}

	// 记录入库日志
	record := models.StockRecord{
		BookID:         req.BookID,
		Type:           "in",
		Quantity:       req.Quantity,
		BeforeQuantity: stock.AvailableQty - req.Quantity,
		AfterQuantity:  stock.AvailableQty,
		Reason:         req.Reason,
		OperatorID:     userID,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "记录入库日志失败")
		return
	}

	tx.Commit()
	utils.SuccessWithMsg(c, "入库成功", stock)
}

// StockOut 出库
func (sc *StockController) StockOut(c *gin.Context) {
	var req struct {
		BookID   uint   `json:"book_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
		Reason   string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("user_id")

	tx := sc.db.Begin()

	var stock models.Stock
	if err := tx.Where("book_id = ?", req.BookID).First(&stock).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "库存记录不存在")
		return
	}

	if stock.AvailableQty < req.Quantity {
		tx.Rollback()
		utils.Error(c, 400, "库存不足")
		return
	}

	// 更新库存
	stock.TotalQuantity -= req.Quantity
	stock.AvailableQty -= req.Quantity
	now := time.Now()
	stock.LastStockOutAt = &now
	if err := tx.Save(&stock).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "更新库存失败")
		return
	}

	// 记录出库日志
	record := models.StockRecord{
		BookID:         req.BookID,
		Type:           "out",
		Quantity:       req.Quantity,
		BeforeQuantity: stock.AvailableQty + req.Quantity,
		AfterQuantity:  stock.AvailableQty,
		Reason:         req.Reason,
		OperatorID:     userID,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "记录出库日志失败")
		return
	}

	tx.Commit()
	utils.SuccessWithMsg(c, "出库成功", stock)
}

// GetStockRecords 获取库存变动记录
func (sc *StockController) GetStockRecords(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		BookID   int    `json:"book_id"`
		Type     string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.Page = 1
		req.PageSize = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var records []models.StockRecord
	var total int64

	query := sc.db.Model(&models.StockRecord{}).Preload("Book").Preload("Operator")
	if req.BookID > 0 {
		query = query.Where("book_id = ?", req.BookID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&records)

	utils.PageSuccess(c, records, total, req.Page, req.PageSize)
}

// GetLowStock 获取低库存预警
func (sc *StockController) GetLowStock(c *gin.Context) {
	var req struct {
		Threshold int `json:"threshold"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Threshold <= 0 {
		req.Threshold = 5
	}

	var stocks []models.Stock
	sc.db.Preload("Book").Where("available_quantity <= ?", req.Threshold).Find(&stocks)

	utils.Success(c, stocks)
}
