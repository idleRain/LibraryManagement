package models

import (
	"time"
)

// BorrowRecord 借阅记录 (MongoDB)
type BorrowRecord struct {
	ID         string     `bson:"_id,omitempty" json:"id"`
	BookID     uint       `bson:"book_id" json:"book_id"`
	BookTitle  string     `bson:"book_title" json:"book_title"`
	BookISBN   string     `bson:"book_isbn" json:"book_isbn"`
	UserID     uint       `bson:"user_id" json:"user_id"`
	UserName   string     `bson:"user_name" json:"user_name"`
	BorrowDate time.Time  `bson:"borrow_date" json:"borrow_date"`
	DueDate    time.Time  `bson:"due_date" json:"due_date"`
	ReturnDate *time.Time `bson:"return_date,omitempty" json:"return_date"`
	Status     string     `bson:"status" json:"status"` // borrowed, returned, overdue, lost
	RenewCount int        `bson:"renew_count" json:"renew_count"`
	Fine       float64    `bson:"fine" json:"fine"`
	FinePaid   bool       `bson:"fine_paid" json:"fine_paid"`
	OperatorID uint       `bson:"operator_id,omitempty" json:"operator_id"`
	Remark     string     `bson:"remark,omitempty" json:"remark"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updated_at"`
}

// CollectionName 返回集合名称
func (BorrowRecord) CollectionName() string {
	return "borrow_records"
}

// BorrowRule 借阅规则
type BorrowRule struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	Name          string    `bson:"name" json:"name"`
	UserType      string    `bson:"user_type" json:"user_type"`
	MaxBooks      int       `bson:"max_books" json:"max_books"`
	MaxDays       int       `bson:"max_days" json:"max_days"`
	MaxRenewTimes int       `bson:"max_renew_times" json:"max_renew_times"`
	FinePerDay    float64   `bson:"fine_per_day" json:"fine_per_day"`
	MaxFine       float64   `bson:"max_fine" json:"max_fine"`
	Description   string    `bson:"description" json:"description"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

// CollectionName 返回集合名称
func (BorrowRule) CollectionName() string {
	return "borrow_rules"
}

// Reservation 预约记录
type Reservation struct {
	ID            string     `bson:"_id,omitempty" json:"id"`
	BookID        uint       `bson:"book_id" json:"book_id"`
	BookTitle     string     `bson:"book_title" json:"book_title"`
	UserID        uint       `bson:"user_id" json:"user_id"`
	UserName      string     `bson:"user_name" json:"user_name"`
	Status        string     `bson:"status" json:"status"` // waiting, notified, completed, cancelled
	QueuePosition int        `bson:"queue_position" json:"queue_position"`
	NotifyDate    *time.Time `bson:"notify_date,omitempty" json:"notify_date"`
	ExpireDate    *time.Time `bson:"expire_date,omitempty" json:"expire_date"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at" json:"updated_at"`
}

// CollectionName 返回集合名称
func (Reservation) CollectionName() string {
	return "reservations"
}

// BorrowStatistics 借阅统计
type BorrowStatistics struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	Date          time.Time `bson:"date" json:"date"`
	TotalBorrowed int       `bson:"total_borrowed" json:"total_borrowed"`
	TotalReturned int       `bson:"total_returned" json:"total_returned"`
	TotalOverdue  int       `bson:"total_overdue" json:"total_overdue"`
	TotalFine     float64   `bson:"total_fine" json:"total_fine"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

// CollectionName 返回集合名称
func (BorrowStatistics) CollectionName() string {
	return "borrow_statistics"
}
