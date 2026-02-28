package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type BatchController struct {
	db *gorm.DB
}

func NewBatchController(db *gorm.DB) *BatchController {
	return &BatchController{db: db}
}

// BatchImportBooks 批量导入图书
func (bc *BatchController) BatchImportBooks(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要导入的文件")
		return
	}
	defer file.Close()

	// 解析 CSV 文件
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // 允许变长字段

	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		utils.Error(c, 400, "读取文件失败")
		return
	}

	// 映射表头索引
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	// 必需字段检查
	requiredFields := []string{"isbn", "title"}
	for _, f := range requiredFields {
		if _, ok := headerMap[f]; !ok {
			utils.Error(c, 400, fmt.Sprintf("缺少必需字段: %s", f))
			return
		}
	}

	// 读取数据行
	var books []models.Book
	var errors []string
	rowNum := 1

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("第 %d 行解析错误", rowNum))
			continue
		}
		rowNum++

		getValue := func(key string) string {
			if i, ok := headerMap[key]; ok && i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}

		book := models.Book{
			ISBN:      getValue("isbn"),
			Title:     getValue("title"),
			Author:    getValue("author"),
			Publisher: getValue("publisher"),
			Category:  getValue("category"),
			Language:  getValue("language"),
			Status:    1,
		}

		// 解析价格
		if price := getValue("price"); price != "" {
			fmt.Sscanf(price, "%f", &book.Price)
		}

		// 解析页数
		if pages := getValue("pages"); pages != "" {
			fmt.Sscanf(pages, "%d", &book.Pages)
		}

		// 解析出版日期
		if date := getValue("publish_date"); date != "" {
			if t, err := time.Parse("2006-01-02", date); err == nil {
				book.PublishDate = &t
			}
		}

		// 检查 ISBN 是否已存在
		var count int64
		bc.db.Model(&models.Book{}).Where("isbn = ?", book.ISBN).Count(&count)
		if count > 0 {
			errors = append(errors, fmt.Sprintf("第 %d 行: ISBN %s 已存在", rowNum, book.ISBN))
			continue
		}

		books = append(books, book)
	}

	// 批量插入
	if len(books) > 0 {
		if err := bc.db.CreateInBatches(&books, 100).Error; err != nil {
			utils.ServerError(c, "导入失败")
			return
		}
	}

	utils.Success(c, gin.H{
		"imported": len(books),
		"errors":   errors,
		"total":    rowNum - 1,
	})
}

// BatchUpdateBooks 批量更新图书
func (bc *BatchController) BatchUpdateBooks(c *gin.Context) {
	var req struct {
		IDs    []int                 `json:"ids" binding:"required"`
		Update map[string]interface{} `json:"update" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if len(req.IDs) == 0 {
		utils.BadRequest(c, "请选择要更新的图书")
		return
	}

	// 过滤允许更新的字段
	allowedFields := map[string]bool{
		"category": true,
		"status":   true,
		"language": true,
	}

	updateData := make(map[string]interface{})
	for k, v := range req.Update {
		if allowedFields[k] {
			updateData[k] = v
		}
	}

	if len(updateData) == 0 {
		utils.BadRequest(c, "没有可更新的字段")
		return
	}

	result := bc.db.Model(&models.Book{}).Where("id IN ?", req.IDs).Updates(updateData)

	utils.SuccessWithMsg(c, "更新成功", gin.H{
		"updated": result.RowsAffected,
	})
}

// BatchDeleteBooks 批量删除图书
func (bc *BatchController) BatchDeleteBooks(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if len(req.IDs) == 0 {
		utils.BadRequest(c, "请选择要删除的图书")
		return
	}

	// 检查是否有库存
	var stockCount int64
	bc.db.Model(&models.Stock{}).Where("book_id IN ?", req.IDs).Count(&stockCount)
	if stockCount > 0 {
		utils.Error(c, 400, "选中的图书中存在库存记录，无法删除")
		return
	}

	result := bc.db.Delete(&models.Book{}, req.IDs)

	utils.SuccessWithMsg(c, "删除成功", gin.H{
		"deleted": result.RowsAffected,
	})
}

// BatchImportUsers 批量导入用户
func (bc *BatchController) BatchImportUsers(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要导入的文件")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		utils.Error(c, 400, "读取文件失败")
		return
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	getValue := func(key string) string {
		if i, ok := headerMap[key]; ok {
			return strings.TrimSpace(headers[i])
		}
		return ""
	}

	var users []models.User
	var errors []string
	rowNum := 1

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("第 %d 行解析错误", rowNum))
			continue
		}
		rowNum++

		getRowValue := func(key string) string {
			if i, ok := headerMap[key]; ok && i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}

		username := getRowValue("username")
		if username == "" {
			errors = append(errors, fmt.Sprintf("第 %d 行: 用户名不能为空", rowNum))
			continue
		}

		// 检查用户名是否已存在
		var count int64
		bc.db.Model(&models.User{}).Where("username = ?", username).Count(&count)
		if count > 0 {
			errors = append(errors, fmt.Sprintf("第 %d 行: 用户名 %s 已存在", rowNum, username))
			continue
		}

		user := models.User{
			Username: username,
			Email:    getRowValue("email"),
			Phone:    getRowValue("phone"),
			RealName: getRowValue("real_name"),
			Status:   1,
		}

		// 默认密码
		password := getRowValue("password")
		if password == "" {
			password = "123456"
		}
		// 这里应该加密密码，简化处理

		users = append(users, user)
	}

	if len(users) > 0 {
		if err := bc.db.CreateInBatches(&users, 100).Error; err != nil {
			utils.ServerError(c, "导入失败")
			return
		}
	}

	utils.Success(c, gin.H{
		"imported": len(users),
		"errors":   errors,
		"total":    rowNum - 1,
	})
}

// BatchUpdateUsers 批量更新用户状态
func (bc *BatchController) BatchUpdateUsers(c *gin.Context) {
	var req struct {
		IDs    []int `json:"ids" binding:"required"`
		Status *int  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if len(req.IDs) == 0 {
		utils.BadRequest(c, "请选择要更新的用户")
		return
	}

	updateData := make(map[string]interface{})
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	if len(updateData) == 0 {
		utils.BadRequest(c, "没有可更新的字段")
		return
	}

	result := bc.db.Model(&models.User{}).Where("id IN ?", req.IDs).Updates(updateData)

	utils.SuccessWithMsg(c, "更新成功", gin.H{
		"updated": result.RowsAffected,
	})
}

// BatchStockIn 批量入库
func (bc *BatchController) BatchStockIn(c *gin.Context) {
	var req struct {
		Items []struct {
			BookID   int `json:"book_id" binding:"required"`
			Quantity int `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("user_id")

	tx := bc.db.Begin()
	var successCount int
	var errors []string

	for i, item := range req.Items {
		var stock models.Stock
		if err := tx.Where("book_id = ?", item.BookID).First(&stock).Error; err != nil {
			errors = append(errors, fmt.Sprintf("第 %d 项: 图书不存在", i+1))
			continue
		}

		beforeQty := stock.AvailableQty
		stock.TotalQuantity += item.Quantity
		stock.AvailableQty += item.Quantity
		now := time.Now()
		stock.LastStockInAt = &now

		if err := tx.Save(&stock).Error; err != nil {
			errors = append(errors, fmt.Sprintf("第 %d 项: 更新库存失败", i+1))
			continue
		}

		// 记录库存变动
		record := models.StockRecord{
			BookID:         uint(item.BookID),
			Type:           "in",
			Quantity:       item.Quantity,
			BeforeQuantity: beforeQty,
			AfterQuantity:  stock.AvailableQty,
			Reason:         req.Reason,
			OperatorID:     &userID,
		}
		tx.Create(&record)

		successCount++
	}

	tx.Commit()

	utils.Success(c, gin.H{
		"success": successCount,
		"errors":  errors,
		"total":   len(req.Items),
	})
}

// DownloadTemplate 下载导入模板
func (bc *BatchController) DownloadTemplate(c *gin.Context) {
	templateType := c.Query("type")

	var content string
	var filename string

	switch templateType {
	case "books":
		filename = "books_import_template.csv"
		content = "isbn,title,author,publisher,publish_date,category,price,pages,language\n"
		content += "978-7-111-54784-2,深入理解计算机系统,Randal E. Bryant,机械工业出版社,2016-07-01,计算机,139.00,737,中文\n"
		content += "978-7-115-42533-4,JavaScript高级程序设计,Matt Frisbie,人民邮电出版社,2020-05-01,计算机,129.00,912,中文\n"
	case "users":
		filename = "users_import_template.csv"
		content = "username,password,email,phone,real_name\n"
		content += "user1,123456,user1@example.com,13800138001,张三\n"
		content += "user2,123456,user2@example.com,13800138002,李四\n"
	default:
		utils.BadRequest(c, "无效的模板类型")
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.String(200, "\xEF\xBB\xBF"+content) // 添加 BOM 以支持中文
}
