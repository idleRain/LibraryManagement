package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type PurchaseController struct {
	db *gorm.DB
}

func NewPurchaseController(db *gorm.DB) *PurchaseController {
	return &PurchaseController{db: db}
}

// GetSuppliers 获取供应商列表
func (pc *PurchaseController) GetSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	pc.db.Find(&suppliers)
	utils.Success(c, suppliers)
}

// CreateSupplier 创建供应商
func (pc *PurchaseController) CreateSupplier(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Contact     string `json:"contact"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Address     string `json:"address"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	supplier := models.Supplier{
		Name:        req.Name,
		Contact:     req.Contact,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		Description: req.Description,
		Status:      1,
	}

	if err := pc.db.Create(&supplier).Error; err != nil {
		utils.ServerError(c, "创建供应商失败")
		return
	}

	utils.SuccessWithMsg(c, "创建成功", supplier)
}

// UpdateSupplier 更新供应商
func (pc *PurchaseController) UpdateSupplier(c *gin.Context) {
	var req struct {
		ID          int    `json:"id" binding:"required"`
		Name        string `json:"name"`
		Contact     string `json:"contact"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Address     string `json:"address"`
		Description string `json:"description"`
		Status      *int8  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var supplier models.Supplier
	if err := pc.db.First(&supplier, req.ID).Error; err != nil {
		utils.NotFound(c, "供应商不存在")
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"contact":     req.Contact,
		"phone":       req.Phone,
		"email":       req.Email,
		"address":     req.Address,
		"description": req.Description,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	pc.db.Model(&supplier).Updates(updates)
	utils.SuccessWithMsg(c, "更新成功", supplier)
}

// DeleteSupplier 删除供应商
func (pc *PurchaseController) DeleteSupplier(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := pc.db.Delete(&models.Supplier{}, req.ID).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// GetPurchaseOrders 获取采购订单列表
func (pc *PurchaseController) GetPurchaseOrders(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Status   string `json:"status"`
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

	var orders []models.PurchaseOrder
	var total int64

	query := pc.db.Model(&models.PurchaseOrder{}).Preload("Supplier").Preload("Operator").Preload("Items.Book")
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&orders)

	utils.PageSuccess(c, orders, total, req.Page, req.PageSize)
}

// GetPurchaseOrder 获取采购订单详情
func (pc *PurchaseController) GetPurchaseOrder(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.PurchaseOrder
	if err := pc.db.Preload("Supplier").Preload("Operator").Preload("Items.Book").First(&order, req.ID).Error; err != nil {
		utils.NotFound(c, "采购订单不存在")
		return
	}

	utils.Success(c, order)
}

// CreatePurchaseOrder 创建采购订单
func (pc *PurchaseController) CreatePurchaseOrder(c *gin.Context) {
	var req struct {
		SupplierID int    `json:"supplier_id" binding:"required"`
		Remark     string `json:"remark"`
		Items      []struct {
			BookID    int     `json:"book_id" binding:"required"`
			Quantity  int     `json:"quantity" binding:"required,min=1"`
			UnitPrice float64 `json:"unit_price" binding:"required"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("user_id")

	// 生成订单号
	orderNo := fmt.Sprintf("PO%s%d", time.Now().Format("20060102150405"), userID)

	tx := pc.db.Begin()

	// 计算总金额
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.UnitPrice * float64(item.Quantity)
	}

	order := models.PurchaseOrder{
		OrderNo:     orderNo,
		SupplierID:  uint(req.SupplierID),
		TotalAmount: totalAmount,
		Status:      "pending",
		OperatorID:  userID,
		Remark:      req.Remark,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "创建采购订单失败")
		return
	}

	// 创建订单明细
	for _, item := range req.Items {
		orderItem := models.PurchaseOrderItem{
			OrderID:    order.ID,
			BookID:     uint(item.BookID),
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.UnitPrice * float64(item.Quantity),
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			utils.ServerError(c, "创建订单明细失败")
			return
		}
	}

	tx.Commit()
	utils.SuccessWithMsg(c, "创建成功", order)
}

// UpdatePurchaseOrderStatus 更新采购订单状态
func (pc *PurchaseController) UpdatePurchaseOrderStatus(c *gin.Context) {
	var req struct {
		ID     int    `json:"id" binding:"required"`
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.PurchaseOrder
	if err := pc.db.First(&order, req.ID).Error; err != nil {
		utils.NotFound(c, "采购订单不存在")
		return
	}

	// 如果是完成状态，需要入库
	if req.Status == "completed" && order.Status != "completed" {
		tx := pc.db.Begin()

		var items []models.PurchaseOrderItem
		tx.Where("order_id = ?", order.ID).Find(&items)

		for _, item := range items {
			var stock models.Stock
			if err := tx.Where("book_id = ?", item.BookID).First(&stock).Error; err != nil {
				tx.Rollback()
				utils.ServerError(c, "库存记录不存在")
				return
			}

			stock.TotalQuantity += item.Quantity - item.ReceivedQuantity
			stock.AvailableQty += item.Quantity - item.ReceivedQuantity
			tx.Save(&stock)

			item.ReceivedQuantity = item.Quantity
			tx.Save(&item)
		}

		order.Status = req.Status
		tx.Save(&order)
		tx.Commit()
	} else {
		pc.db.Model(&order).Update("status", req.Status)
	}

	utils.SuccessWithMsg(c, "更新成功", order)
}

// DeletePurchaseOrder 删除采购订单
func (pc *PurchaseController) DeletePurchaseOrder(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.PurchaseOrder
	if err := pc.db.First(&order, req.ID).Error; err != nil {
		utils.NotFound(c, "采购订单不存在")
		return
	}

	if order.Status != "pending" && order.Status != "cancelled" {
		utils.Error(c, 400, "只能删除待审核或已取消的订单")
		return
	}

	tx := pc.db.Begin()
	tx.Where("order_id = ?", order.ID).Delete(&models.PurchaseOrderItem{})
	tx.Delete(&order)
	tx.Commit()

	utils.SuccessWithMsg(c, "删除成功", nil)
}
