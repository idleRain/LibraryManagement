package models

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商模型
type Supplier struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Contact     string         `json:"contact" gorm:"size:50"`
	Phone       string         `json:"phone" gorm:"size:20"`
	Email       string         `json:"email" gorm:"size:100"`
	Address     string         `json:"address" gorm:"size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Status      int            `json:"status" gorm:"default:1"` // 1: 启用, 0: 禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// PurchaseOrder 采购订单
type PurchaseOrder struct {
	ID            uint                `json:"id" gorm:"primaryKey"`
	OrderNo       string              `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	SupplierID    uint                `json:"supplier_id" gorm:"index;not null"`
	Supplier      *Supplier           `json:"supplier" gorm:"foreignKey:SupplierID"`
	TotalAmount   float64             `json:"total_amount" gorm:"type:decimal(12,2)"`
	Status        string              `json:"status" gorm:"size:20;not null"` // pending, approved, completed, cancelled
	OperatorID    uint                `json:"operator_id" gorm:"index"`
	Operator      *User               `json:"operator" gorm:"foreignKey:OperatorID"`
	Remark        string              `json:"remark" gorm:"type:text"`
	Items         []PurchaseOrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `json:"-" gorm:"index"`
}

// PurchaseOrderItem 采购订单明细
type PurchaseOrderItem struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	OrderID      uint           `json:"order_id" gorm:"index;not null"`
	Order        *PurchaseOrder `json:"order" gorm:"foreignKey:OrderID"`
	BookID       uint           `json:"book_id" gorm:"index;not null"`
	Book         *Book          `json:"book" gorm:"foreignKey:BookID"`
	Quantity     int            `json:"quantity" gorm:"not null"`
	UnitPrice    float64        `json:"unit_price" gorm:"type:decimal(10,2)"`
	TotalPrice   float64        `json:"total_price" gorm:"type:decimal(10,2)"`
	ReceivedQty  int            `json:"received_quantity" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Supplier) TableName() string {
	return "suppliers"
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}
