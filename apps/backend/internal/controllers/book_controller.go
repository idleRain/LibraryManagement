package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type BookController struct {
	db *gorm.DB
}

func NewBookController(db *gorm.DB) *BookController {
	return &BookController{db: db}
}

// GetBooks 获取图书列表
func (bc *BookController) GetBooks(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Title    string `json:"title"`
		Author   string `json:"author"`
		Category string `json:"category"`
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

	var books []models.Book
	var total int64

	query := bc.db.Model(&models.Book{}).Preload("Stock")
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Author != "" {
		query = query.Where("author LIKE ?", "%"+req.Author+"%")
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&books)

	utils.PageSuccess(c, books, total, req.Page, req.PageSize)
}

// GetBook 获取图书详情
func (bc *BookController) GetBook(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var book models.Book
	if err := bc.db.Preload("Stock").First(&book, req.ID).Error; err != nil {
		utils.NotFound(c, "图书不存在")
		return
	}

	utils.Success(c, book)
}

// CreateBook 创建图书
func (bc *BookController) CreateBook(c *gin.Context) {
	var req struct {
		ISBN        string     `json:"isbn" binding:"required"`
		Title       string     `json:"title" binding:"required"`
		Author      string     `json:"author"`
		Publisher   string     `json:"publisher"`
		PublishDate *time.Time `json:"publish_date"`
		Category    string     `json:"category"`
		Price       float64    `json:"price"`
		Description string     `json:"description"`
		CoverImage  string     `json:"cover_image"`
		Pages       int        `json:"pages"`
		Language    string     `json:"language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	book := models.Book{
		ISBN:        req.ISBN,
		Title:       req.Title,
		Author:      req.Author,
		Publisher:   req.Publisher,
		PublishDate: req.PublishDate,
		Category:    req.Category,
		Price:       req.Price,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Pages:       req.Pages,
		Language:    req.Language,
	}

	tx := bc.db.Begin()
	if err := tx.Create(&book).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "创建图书失败")
		return
	}

	// 创建库存记录
	stock := models.Stock{
		BookID:        book.ID,
		TotalQuantity: 0,
		AvailableQty:  0,
		BorrowedQty:   0,
		SoldQty:       0,
	}
	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "创建库存记录失败")
		return
	}

	tx.Commit()
	utils.SuccessWithMsg(c, "创建成功", book)
}

// UpdateBook 更新图书
func (bc *BookController) UpdateBook(c *gin.Context) {
	var req struct {
		ID          int        `json:"id" binding:"required"`
		Title       string     `json:"title"`
		Author      string     `json:"author"`
		Publisher   string     `json:"publisher"`
		PublishDate *time.Time `json:"publish_date"`
		Category    string     `json:"category"`
		Price       float64    `json:"price"`
		Description string     `json:"description"`
		CoverImage  string     `json:"cover_image"`
		Pages       int        `json:"pages"`
		Language    string     `json:"language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var book models.Book
	if err := bc.db.First(&book, req.ID).Error; err != nil {
		utils.NotFound(c, "图书不存在")
		return
	}

	updates := map[string]interface{}{
		"title":        req.Title,
		"author":       req.Author,
		"publisher":    req.Publisher,
		"publish_date": req.PublishDate,
		"category":     req.Category,
		"price":        req.Price,
		"description":  req.Description,
		"cover_image":  req.CoverImage,
		"pages":        req.Pages,
		"language":     req.Language,
	}

	bc.db.Model(&book).Updates(updates)
	utils.SuccessWithMsg(c, "更新成功", book)
}

// DeleteBook 删除图书
func (bc *BookController) DeleteBook(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tx := bc.db.Begin()

	// 删除库存
	tx.Delete(&models.Stock{}, "book_id = ?", req.ID)
	// 删除图书
	if err := tx.Delete(&models.Book{}, req.ID).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "删除失败")
		return
	}

	tx.Commit()
	utils.SuccessWithMsg(c, "删除成功", nil)
}

// GetCategories 获取图书分类
func (bc *BookController) GetCategories(c *gin.Context) {
	var categories []string
	bc.db.Model(&models.Book{}).Distinct().Pluck("category", &categories)
	utils.Success(c, categories)
}
