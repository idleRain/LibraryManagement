package models

import (
	"time"

	"gorm.io/gorm"
)

// Book 图书模型
type Book struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ISBN        string         `gorm:"size:20;uniqueIndex;not null" json:"isbn"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Subtitle    string         `gorm:"size:200" json:"subtitle"`
	Author      string         `gorm:"size:100" json:"author"`
	Translator  string         `gorm:"size:100" json:"translator"`
	Publisher   string         `gorm:"size:100" json:"publisher"`
	PublishDate *time.Time     `json:"publish_date"`
	Edition     string         `gorm:"size:50" json:"edition"`
	Category    string         `gorm:"size:50" json:"category"`
	Language    string         `gorm:"size:20;default:中文" json:"language"`
	Pages       int            `json:"pages"`
	Price       float64        `json:"price"`
	CoverImage  string         `gorm:"size:255" json:"cover_image"`
	Description string         `json:"description"`
	Keywords    string         `gorm:"size:255" json:"keywords"`
	Status      int8           `gorm:"default:1" json:"status"` // 1:上架 0:下架
	Stock       *Stock         `gorm:"foreignKey:BookID" json:"stock"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Book) TableName() string {
	return "books"
}

// Stock 库存模型
type Stock struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	BookID           uint           `gorm:"uniqueIndex;not null" json:"book_id"`
	Book             *Book          `gorm:"foreignKey:BookID" json:"book"`
	TotalQuantity    int            `gorm:"default:0" json:"total_quantity"`
	AvailableQty     int            `gorm:"default:0" json:"available_quantity"`
	BorrowedQty      int            `gorm:"default:0" json:"borrowed_quantity"`
	SoldQty          int            `gorm:"default:0" json:"sold_quantity"`
	DamagedQty       int            `gorm:"default:0" json:"damaged_quantity"`
	Location         string         `gorm:"size:100" json:"location"`
	WarningThreshold int            `gorm:"default:5" json:"warning_threshold"`
	LastStockInAt    *time.Time     `json:"last_stock_in_at"`
	LastStockOutAt   *time.Time     `json:"last_stock_out_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Stock) TableName() string {
	return "stocks"
}

// StockRecord 库存变动记录
type StockRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	BookID          uint           `gorm:"not null;index" json:"book_id"`
	Book            *Book          `gorm:"foreignKey:BookID" json:"book"`
	Type            string         `gorm:"size:10;not null" json:"type"` // in, out, adjust
	Quantity        int            `gorm:"not null" json:"quantity"`
	BeforeQuantity  int            `json:"before_quantity"`
	AfterQuantity   int            `json:"after_quantity"`
	Reason          string         `gorm:"size:255" json:"reason"`
	RelatedOrderType string        `gorm:"size:50" json:"related_order_type"`
	RelatedOrderID  uint           `json:"related_order_id"`
	OperatorID      uint           `json:"operator_id"`
	Operator        *User          `gorm:"foreignKey:OperatorID" json:"operator"`
	CreatedAt       time.Time      `json:"created_at"`
}

// TableName 指定表名
func (StockRecord) TableName() string {
	return "stock_records"
}
