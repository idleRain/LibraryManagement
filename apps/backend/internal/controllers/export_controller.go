package controllers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExportController struct {
	db *gorm.DB
}

func NewExportController(db *gorm.DB) *ExportController {
	return &ExportController{db: db}
}

// ExportBooks 导出图书列表
func (ec *ExportController) ExportBooks(c *gin.Context) {
	var req struct {
		Title    string `json:"title"`
		Author   string `json:"author"`
		Category string `json:"category"`
	}

	c.ShouldBindJSON(&req)

	// 查询数据
	var books []models.Book
	query := ec.db.Preload("Stock")
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Author != "" {
		query = query.Where("author LIKE ?", "%"+req.Author+"%")
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	query.Find(&books)

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	sheet := "图书列表"
	f.SetSheetName("Sheet1", sheet)

	// 设置表头
	headers := []string{"ID", "ISBN", "书名", "作者", "出版社", "分类", "价格", "库存", "可借", "状态", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	// 填充数据
	for i, book := range books {
		row := i + 2
		stock := 0
		available := 0
		if book.Stock != nil {
			stock = book.Stock.TotalQuantity
			available = book.Stock.AvailableQty
		}
		status := "上架"
		if book.Status != 1 {
			status = "下架"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), book.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), book.ISBN)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), book.Title)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), book.Author)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), book.Publisher)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), book.Category)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), book.Price)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), stock)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), available)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), status)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), book.CreatedAt.Format("2006-01-02"))
	}

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 30)
	f.SetColWidth(sheet, "D", "E", 15)
	f.SetColWidth(sheet, "F", "F", 12)
	f.SetColWidth(sheet, "G", "G", 10)
	f.SetColWidth(sheet, "H", "I", 8)
	f.SetColWidth(sheet, "J", "J", 8)
	f.SetColWidth(sheet, "K", "K", 12)

	// 写入缓冲区
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		utils.ServerError(c, "生成文件失败")
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("books_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportBorrows 导出借阅记录
func (ec *ExportController) ExportBorrows(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}

	c.ShouldBindJSON(&req)

	// 从 MongoDB 查询数据（这里简化为从 MySQL 查询）
	// 实际项目中应该从 MongoDB 查询

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	sheet := "借阅记录"
	f.SetSheetName("Sheet1", sheet)

	// 设置表头
	headers := []string{"ID", "图书名称", "ISBN", "借阅人", "借阅日期", "应还日期", "归还日期", "状态", "续借次数", "罚款"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 25)
	f.SetColWidth(sheet, "B", "B", 30)
	f.SetColWidth(sheet, "C", "C", 18)
	f.SetColWidth(sheet, "D", "D", 12)
	f.SetColWidth(sheet, "E", "G", 12)
	f.SetColWidth(sheet, "H", "H", 10)
	f.SetColWidth(sheet, "I", "I", 10)
	f.SetColWidth(sheet, "J", "J", 10)

	// 写入缓冲区
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		utils.ServerError(c, "生成文件失败")
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("borrows_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportSales 导出销售记录
func (ec *ExportController) ExportSales(c *gin.Context) {
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	c.ShouldBindJSON(&req)

	// 查询数据
	var orders []models.SaleOrder
	query := ec.db.Preload("Items.Book")
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}
	query.Order("created_at DESC").Find(&orders)

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	sheet := "销售记录"
	f.SetSheetName("Sheet1", sheet)

	// 设置表头
	headers := []string{"订单号", "客户姓名", "客户电话", "图书名称", "数量", "单价", "总价", "订单金额", "支付方式", "状态", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	// 填充数据
	row := 2
	for _, order := range orders {
		for _, item := range order.Items {
			bookName := ""
			if item.Book != nil {
				bookName = item.Book.Title
			}
			status := "待支付"
			switch order.Status {
			case "paid":
				status = "已支付"
			case "completed":
				status = "已完成"
			case "cancelled":
				status = "已取消"
			}
			paymentMethod := order.PaymentMethod
			if paymentMethod == "" {
				paymentMethod = "-"
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), order.OrderNo)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), order.CustomerName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), order.CustomerPhone)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), bookName)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.UnitPrice)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), item.TotalPrice)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), order.TotalAmount)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), paymentMethod)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), status)
			f.SetCellValue(sheet, fmt.Sprintf("K%d", row), order.CreatedAt.Format("2006-01-02 15:04:05"))
			row++
		}
	}

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 20)
	f.SetColWidth(sheet, "B", "B", 12)
	f.SetColWidth(sheet, "C", "C", 15)
	f.SetColWidth(sheet, "D", "D", 30)
	f.SetColWidth(sheet, "E", "E", 8)
	f.SetColWidth(sheet, "F", "H", 12)
	f.SetColWidth(sheet, "I", "I", 10)
	f.SetColWidth(sheet, "J", "J", 10)
	f.SetColWidth(sheet, "K", "K", 18)

	// 写入缓冲区
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		utils.ServerError(c, "生成文件失败")
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("sales_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ExportUsers 导出用户列表
func (ec *ExportController) ExportUsers(c *gin.Context) {
	// 查询数据
	var users []models.User
	ec.db.Preload("Roles").Find(&users)

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer f.Close()

	sheet := "用户列表"
	f.SetSheetName("Sheet1", sheet)

	// 设置表头
	headers := []string{"ID", "用户名", "邮箱", "手机号", "真实姓名", "角色", "状态", "最后登录时间", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 设置表头样式
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	// 填充数据
	for i, user := range users {
		row := i + 2
		roles := ""
		for j, r := range user.Roles {
			if j > 0 {
				roles += ", "
			}
			roles += r.Name
		}
		status := "启用"
		if user.Status != 1 {
			status = "禁用"
		}
		lastLogin := ""
		if user.LastLoginAt != nil {
			lastLogin = user.LastLoginAt.Format("2006-01-02 15:04:05")
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), user.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), user.Username)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), user.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), user.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), user.RealName)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), roles)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), status)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), lastLogin)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), user.CreatedAt.Format("2006-01-02"))
	}

	// 设置列宽
	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 25)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 12)
	f.SetColWidth(sheet, "F", "F", 20)
	f.SetColWidth(sheet, "G", "G", 8)
	f.SetColWidth(sheet, "H", "H", 18)
	f.SetColWidth(sheet, "I", "I", 12)

	// 写入缓冲区
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		utils.ServerError(c, "生成文件失败")
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("users_%s.xlsx", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
