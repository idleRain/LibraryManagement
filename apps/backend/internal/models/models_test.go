package models

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	user := User{
		Username:  "testuser",
		Password:  "hashedpassword",
		Email:     "test@example.com",
		Status:    1,
		CreatedAt: time.Now(),
	}

	if user.Username != "testuser" {
		t.Errorf("Expected Username = testuser, got %s", user.Username)
	}
}

func TestBook(t *testing.T) {
	book := Book{
		ISBN:      "978-7-111-54784-2",
		Title:     "深入理解计算机系统",
		Author:    "Randal E. Bryant",
		Publisher: "机械工业出版社",
		Price:     139.00,
		Status:    1,
	}

	if book.ISBN != "978-7-111-54784-2" {
		t.Errorf("Expected ISBN = 978-7-111-54784-2, got %s", book.ISBN)
	}
}

func TestStock(t *testing.T) {
	stock := Stock{
		BookID:           1,
		TotalQuantity:    100,
		AvailableQty:     80,
		BorrowedQty:      10,
		SoldQty:          10,
		WarningThreshold: 10,
	}

	if stock.AvailableQty != 80 {
		t.Errorf("Expected AvailableQty = 80, got %d", stock.AvailableQty)
	}
}
