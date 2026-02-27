package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type SaleController struct {
	db *gorm.DB
}

func NewSaleController(db *gorm.DB) *SaleController {
	return &SaleController{db: db}
}

// GetSaleOrders 获取销售订单列表
func (sc *SaleController) GetSaleOrders(c *gin.Context) {
	var req struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"page_size"`
		Status    string `json:"status"`
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
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var orders []models.SaleOrder
	var total int64

	query := sc.db.Model(&models.SaleOrder{}).Preload("Operator").Preload("Items.Book")
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&orders)

	utils.PageSuccess(c, orders, total, req.Page, req.PageSize)
}

// GetSaleOrder 获取销售订单详情
func (sc *SaleController) GetSaleOrder(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.SaleOrder
	if err := sc.db.Preload("Operator").Preload("Items.Book").First(&order, req.ID).Error; err != nil {
		utils.NotFound(c, "销售订单不存在")
		return
	}

	utils.Success(c, order)
}

// CreateSaleOrder 创建销售订单
func (sc *SaleController) CreateSaleOrder(c *gin.Context) {
	var req struct {
		CustomerName  string  `json:"customer_name"`
		CustomerPhone string  `json:"customer_phone"`
		Discount      float64 `json:"discount"`
		PaymentMethod string  `json:"payment_method"`
		Remark        string  `json:"remark"`
		Items         []struct {
			BookID    int     `json:"book_id" binding:"required"`
			Quantity  int     `json:"quantity" binding:"required,min=1"`
			UnitPrice float64 `json:"unit_price" binding:"required"`
			Discount  float64 `json:"discount"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("user_id")

	// 生成订单号
	orderNo := fmt.Sprintf("SO%s%d", time.Now().Format("20060102150405"), userID)

	tx := sc.db.Begin()

	// 计算总金额
	var totalAmount float64
	for _, item := range req.Items {
		discount := item.Discount
		if discount == 0 {
			discount = 100
		}
		totalAmount += item.UnitPrice * float64(item.Quantity) * discount / 100
	}

	order := models.SaleOrder{
		OrderNo:       orderNo,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TotalAmount:   totalAmount,
		PayAmount:     totalAmount,
		Discount:      req.Discount,
		Status:        "paid",
		PaymentMethod: req.PaymentMethod,
		OperatorID:    userID,
		Remark:        req.Remark,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "创建销售订单失败")
		return
	}

	// 创建订单明细并扣减库存
	for _, item := range req.Items {
		discount := item.Discount
		if discount == 0 {
			discount = 100
		}

		orderItem := models.SaleOrderItem{
			OrderID:    order.ID,
			BookID:     uint(item.BookID),
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.UnitPrice * float64(item.Quantity) * discount / 100,
			Discount:   discount,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			utils.ServerError(c, "创建订单明细失败")
			return
		}

		// 扣减库存
		var stock models.Stock
		if err := tx.Where("book_id = ?", item.BookID).First(&stock).Error; err != nil {
			tx.Rollback()
			utils.Error(c, 400, "图书库存不存在")
			return
		}

		if stock.AvailableQty < item.Quantity {
			tx.Rollback()
			utils.Error(c, 400, "库存不足")
			return
		}

		stock.AvailableQty -= item.Quantity
		stock.TotalQuantity -= item.Quantity
		stock.SoldQty += item.Quantity
		if err := tx.Save(&stock).Error; err != nil {
			tx.Rollback()
			utils.ServerError(c, "更新库存失败")
			return
		}
	}

	tx.Commit()
	utils.SuccessWithMsg(c, "创建成功", order)
}

// CancelSaleOrder 取消销售订单
func (sc *SaleController) CancelSaleOrder(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.SaleOrder
	if err := sc.db.Preload("Items").First(&order, req.ID).Error; err != nil {
		utils.NotFound(c, "销售订单不存在")
		return
	}

	if order.Status == "cancelled" {
		utils.Error(c, 400, "订单已取消")
		return
	}

	tx := sc.db.Begin()

	// 恢复库存
	for _, item := range order.Items {
		var stock models.Stock
		if err := tx.Where("book_id = ?", item.BookID).First(&stock).Error; err != nil {
			tx.Rollback()
			utils.ServerError(c, "库存记录不存在")
			return
		}

		stock.AvailableQty += item.Quantity
		stock.TotalQuantity += item.Quantity
		stock.SoldQty -= item.Quantity
		tx.Save(&stock)
	}

	order.Status = "cancelled"
	tx.Save(&order)
	tx.Commit()

	utils.SuccessWithMsg(c, "取消成功", order)
}

// GetSalesStats 获取销售统计
func (sc *SaleController) GetSalesStats(c *gin.Context) {
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.StartDate = time.Now().Format("2006-01-02")
		req.EndDate = time.Now().Format("2006-01-02")
	}

	if req.StartDate == "" {
		req.StartDate = time.Now().Format("2006-01-02")
	}
	if req.EndDate == "" {
		req.EndDate = time.Now().Format("2006-01-02")
	}

	var stats struct {
		TotalOrders   int64   `json:"total_orders"`
		TotalAmount   float64 `json:"total_amount"`
		TotalQuantity int64   `json:"total_quantity"`
	}

	sc.db.Model(&models.SaleOrder{}).
		Where("status = 'paid' AND created_at >= ? AND created_at <= ?", req.StartDate, req.EndDate+" 23:59:59").
		Count(&stats.TotalOrders)

	sc.db.Model(&models.SaleOrder{}).
		Where("status = 'paid' AND created_at >= ? AND created_at <= ?", req.StartDate, req.EndDate+" 23:59:59").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalAmount)

	sc.db.Table("sale_order_items").
		Joins("JOIN sale_orders ON sale_orders.id = sale_order_items.order_id").
		Where("sale_orders.status = 'paid' AND sale_orders.created_at >= ? AND sale_orders.created_at <= ?", req.StartDate, req.EndDate+" 23:59:59").
		Select("COALESCE(SUM(sale_order_items.quantity), 0)").Scan(&stats.TotalQuantity)

	utils.Success(c, stats)
}

// GetCart 获取购物车
func (sc *SaleController) GetCart(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	userID, _ := c.Get("user_id")

	var carts []models.Cart
	query := sc.db.Preload("Book")
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}
	query.Find(&carts)

	utils.Success(c, carts)
}

// AddToCart 添加到购物车
func (sc *SaleController) AddToCart(c *gin.Context) {
	var req struct {
		BookID   int `json:"book_id" binding:"required"`
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	sessionID := c.GetHeader("X-Session-ID")
	userID, _ := c.Get("user_id")

	// 检查库存
	var stock models.Stock
	if err := sc.db.Where("book_id = ?", req.BookID).First(&stock).Error; err != nil {
		utils.Error(c, 400, "图书不存在")
		return
	}

	if stock.AvailableQty < req.Quantity {
		utils.Error(c, 400, "库存不足")
		return
	}

	cart := models.Cart{
		SessionID: sessionID,
		BookID:    uint(req.BookID),
		Quantity:  req.Quantity,
	}

	if userID != nil {
		cart.UserID = userID.(uint)
	}

	// 检查是否已存在
	var existingCart models.Cart
	query := sc.db.Where("book_id = ?", req.BookID)
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	if err := query.First(&existingCart).Error; err == nil {
		existingCart.Quantity += req.Quantity
		sc.db.Save(&existingCart)
		utils.SuccessWithMsg(c, "更新成功", existingCart)
		return
	}

	if err := sc.db.Create(&cart).Error; err != nil {
		utils.ServerError(c, "添加失败")
		return
	}

	utils.SuccessWithMsg(c, "添加成功", cart)
}

// RemoveFromCart 从购物车移除
func (sc *SaleController) RemoveFromCart(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := sc.db.Delete(&models.Cart{}, req.ID).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}
