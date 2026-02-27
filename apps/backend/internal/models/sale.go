package models

import (
	"time"

	"gorm.io/gorm"
)

// SaleOrder 销售订单
type SaleOrder struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	OrderNo      string          `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	CustomerName string          `json:"customer_name" gorm:"size:50"`
	CustomerPhone string         `json:"customer_phone" gorm:"size:20"`
	TotalAmount  float64         `json:"total_amount" gorm:"type:decimal(12,2)"`
	PayAmount    float64         `json:"pay_amount" gorm:"type:decimal(12,2)"`
	Discount     float64         `json:"discount" gorm:"type:decimal(5,2);default:100"` // 折扣百分比
	Status       string          `json:"status" gorm:"size:20;not null"` // pending, paid, completed, cancelled
	PaymentMethod string          `json:"payment_method" gorm:"size:20"` // cash, wechat, alipay, card
	OperatorID   uint            `json:"operator_id" gorm:"index"`
	Operator     *User           `json:"operator" gorm:"foreignKey:OperatorID"`
	Remark       string          `json:"remark" gorm:"type:text"`
	Items        []SaleOrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `json:"-" gorm:"index"`
}

// SaleOrderItem 销售订单明细
type SaleOrderItem struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	OrderID    uint        `json:"order_id" gorm:"index;not null"`
	Order      *SaleOrder  `json:"order" gorm:"foreignKey:OrderID"`
	BookID     uint        `json:"book_id" gorm:"index;not null"`
	Book       *Book       `json:"book" gorm:"foreignKey:BookID"`
	Quantity   int         `json:"quantity" gorm:"not null"`
	UnitPrice  float64     `json:"unit_price" gorm:"type:decimal(10,2)"`
	TotalPrice float64     `json:"total_price" gorm:"type:decimal(10,2)"`
	Discount   float64     `json:"discount" gorm:"type:decimal(5,2);default:100"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// Cart 购物车
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SessionID string    `json:"session_id" gorm:"index;size:100"`
	UserID    uint      `json:"user_id" gorm:"index"`
	User      *User     `json:"user" gorm:"foreignKey:UserID"`
	BookID    uint      `json:"book_id" gorm:"index;not null"`
	Book      *Book     `json:"book" gorm:"foreignKey:BookID"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SaleOrder) TableName() string {
	return "sale_orders"
}

func (SaleOrderItem) TableName() string {
	return "sale_order_items"
}

func (Cart) TableName() string {
	return "carts"
}
