package models

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商模型
type Supplier struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex" json:"code"`
	Contact     string         `gorm:"size:50" json:"contact"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Email       string         `gorm:"size:100" json:"email"`
	Address     string         `gorm:"size:255" json:"address"`
	BankName    string         `gorm:"size:100" json:"bank_name"`
	BankAccount string         `gorm:"size:50" json:"bank_account"`
	TaxNumber   string         `gorm:"size:50" json:"tax_number"`
	Description string         `json:"description"`
	Status      int8           `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Supplier) TableName() string {
	return "suppliers"
}

// PurchaseOrder 采购订单
type PurchaseOrder struct {
	ID           uint                  `gorm:"primaryKey" json:"id"`
	OrderNo      string                `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	SupplierID   uint                  `gorm:"not null;index" json:"supplier_id"`
	Supplier     *Supplier             `gorm:"foreignKey:SupplierID" json:"supplier"`
	TotalAmount  float64               `gorm:"default:0" json:"total_amount"`
	PaidAmount   float64               `gorm:"default:0" json:"paid_amount"`
	Status       string                `gorm:"size:20;default:pending" json:"status"` // pending, approved, receiving, completed, cancelled
	OperatorID   uint                  `json:"operator_id"`
	Operator     *User                 `gorm:"foreignKey:OperatorID" json:"operator"`
	ApprovedBy   uint                  `json:"approved_by"`
	ApprovedAt   *time.Time            `json:"approved_at"`
	Remark       string                `json:"remark"`
	Items        []PurchaseOrderItem   `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	DeletedAt    gorm.DeletedAt        `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// PurchaseOrderItem 采购订单明细
type PurchaseOrderItem struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	OrderID          uint           `gorm:"not null;index" json:"order_id"`
	BookID           uint           `gorm:"not null" json:"book_id"`
	Book             *Book          `gorm:"foreignKey:BookID" json:"book"`
	Quantity         int            `gorm:"not null" json:"quantity"`
	UnitPrice        float64        `gorm:"not null" json:"unit_price"`
	TotalPrice       float64        `gorm:"not null" json:"total_price"`
	ReceivedQuantity int            `gorm:"default:0" json:"received_quantity"`
	Status           string         `gorm:"size:20;default:pending" json:"status"` // pending, partial, completed
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

// SaleOrder 销售订单
type SaleOrder struct {
	ID             uint               `gorm:"primaryKey" json:"id"`
	OrderNo        string             `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	CustomerName   string             `gorm:"size:50" json:"customer_name"`
	CustomerPhone  string             `gorm:"size:20" json:"customer_phone"`
	CustomerEmail  string             `gorm:"size:100" json:"customer_email"`
	TotalAmount    float64            `gorm:"default:0" json:"total_amount"`
	DiscountAmount float64            `gorm:"default:0" json:"discount_amount"`
	PayAmount      float64            `gorm:"default:0" json:"pay_amount"`
	PaymentMethod  string             `gorm:"size:20" json:"payment_method"` // cash, wechat, alipay, card, other
	PaymentStatus  string             `gorm:"size:20;default:unpaid" json:"payment_status"` // unpaid, paid, refunded
	Status         string             `gorm:"size:20;default:pending" json:"status"` // pending, paid, completed, cancelled
	OperatorID     uint               `json:"operator_id"`
	Operator       *User              `gorm:"foreignKey:OperatorID" json:"operator"`
	Remark         string             `json:"remark"`
	PaidAt         *time.Time         `json:"paid_at"`
	Items          []SaleOrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SaleOrder) TableName() string {
	return "sale_orders"
}

// SaleOrderItem 销售订单明细
type SaleOrderItem struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	OrderID     uint       `gorm:"not null;index" json:"order_id"`
	BookID      uint       `gorm:"not null" json:"book_id"`
	Book        *Book      `gorm:"foreignKey:BookID" json:"book"`
	Quantity    int        `gorm:"not null" json:"quantity"`
	UnitPrice   float64    `gorm:"not null" json:"unit_price"`
	Discount    float64    `gorm:"default:100" json:"discount"` // 折扣百分比
	TotalPrice  float64    `gorm:"not null" json:"total_price"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TableName 指定表名
func (SaleOrderItem) TableName() string {
	return "sale_order_items"
}

// Cart 购物车
type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	SessionID string     `gorm:"size:100;index" json:"session_id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	BookID    uint       `gorm:"not null" json:"book_id"`
	Book      *Book      `gorm:"foreignKey:BookID" json:"book"`
	Quantity  int        `gorm:"default:1" json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Cart) TableName() string {
	return "carts"
}
