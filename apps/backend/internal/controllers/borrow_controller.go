package controllers

import (
        "context"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/library-system/backend/internal/models"
        "github.com/library-system/backend/internal/utils"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "gorm.io/gorm"
)

type BorrowController struct {
        mysqlDB    *gorm.DB
        mongoDB    *mongo.Database
        collection *mongo.Collection
}

func NewBorrowController(mysqlDB *gorm.DB, mongoDB *mongo.Database) *BorrowController {
        return &BorrowController{
                mysqlDB:    mysqlDB,
                mongoDB:    mongoDB,
                collection: mongoDB.Collection("borrow_records"),
        }
}

// GetBorrowRecords 获取借阅记录列表
func (bc *BorrowController) GetBorrowRecords(c *gin.Context) {
        var req struct {
                Page     int    `json:"page"`
                PageSize int    `json:"page_size"`
                Status   string `json:"status"`
                UserID   int    `json:"user_id"`
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

        filter := bson.M{}
        if req.Status != "" {
                filter["status"] = req.Status
        }
        if req.UserID > 0 {
                filter["user_id"] = req.UserID
        }

        skip := int64((req.Page - 1) * req.PageSize)
        limit := int64(req.PageSize)

        total, _ := bc.collection.CountDocuments(context.Background(), filter)

        opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
        cursor, err := bc.collection.Find(context.Background(), filter, opts)
        if err != nil {
                utils.ServerError(c, "查询失败")
                return
        }
        defer cursor.Close(context.Background())

        var records []models.BorrowRecord
        if err := cursor.All(context.Background(), &records); err != nil {
                utils.ServerError(c, "解析数据失败")
                return
        }

        utils.PageSuccess(c, records, total, req.Page, req.PageSize)
}

// GetBorrowRecord 获取借阅记录详情
func (bc *BorrowController) GetBorrowRecord(c *gin.Context) {
        var req struct {
                ID string `json:"id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        objID, err := primitive.ObjectIDFromHex(req.ID)
        if err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        var record models.BorrowRecord
        err = bc.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&record)
        if err != nil {
                utils.NotFound(c, "借阅记录不存在")
                return
        }

        utils.Success(c, record)
}

// BorrowBook 借书
func (bc *BorrowController) BorrowBook(c *gin.Context) {
        var req struct {
                BookID int `json:"book_id" binding:"required"`
                Days   int `json:"days"` // 借阅天数，默认30天
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        userID := c.GetUint("user_id")
        var user models.User
        if err := bc.mysqlDB.First(&user, userID).Error; err != nil {
                utils.NotFound(c, "用户不存在")
                return
        }

        // 检查图书库存
        var book models.Book
        if err := bc.mysqlDB.Preload("Stock").First(&book, req.BookID).Error; err != nil {
                utils.NotFound(c, "图书不存在")
                return
        }

        if book.Stock == nil || book.Stock.AvailableQty <= 0 {
                utils.Error(c, 400, "库存不足")
                return
        }

        // 检查用户是否已借阅此书
        count, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "user_id": userID,
                "book_id": req.BookID,
                "status":  bson.M{"$in": []string{"borrowed", "overdue"}},
        })
        if count > 0 {
                utils.Error(c, 400, "您已借阅此书，请先归还")
                return
        }

        // 默认借阅30天
        if req.Days <= 0 {
                req.Days = 30
        }

        now := time.Now()
        dueDate := now.AddDate(0, 0, req.Days)

        // 获取用户显示名称
        displayName := user.Username
        if user.RealName != "" {
                displayName = user.RealName
        }

        record := models.BorrowRecord{
                BookID:     uint(req.BookID),
                BookTitle:  book.Title,
                BookISBN:   book.ISBN,
                UserID:     userID,
                UserName:   displayName,
                BorrowDate: now,
                DueDate:    dueDate,
                Status:     "borrowed",
                RenewCount: 0,
                Fine:       0,
                FinePaid:   false,
                CreatedAt:  now,
                UpdatedAt:  now,
        }

        // 插入借阅记录
        _, err = bc.collection.InsertOne(context.Background(), record)
        if err != nil {
                utils.ServerError(c, "借阅失败")
                return
        }

        // 更新库存
        bc.mysqlDB.Model(&book.Stock).Updates(map[string]interface{}{
                "available_quantity": book.Stock.AvailableQty - 1,
                "borrowed_quantity":  book.Stock.BorrowedQty + 1,
        })

        utils.SuccessWithMsg(c, "借阅成功", record)
}

// ReturnBook 还书
func (bc *BorrowController) ReturnBook(c *gin.Context) {
        var req struct {
                ID string `json:"id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        objID, err := primitive.ObjectIDFromHex(req.ID)
        if err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        var record models.BorrowRecord
        err = bc.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&record)
        if err != nil {
                utils.NotFound(c, "借阅记录不存在")
                return
        }

        if record.Status == "returned" {
                utils.Error(c, 400, "图书已归还")
                return
        }

        now := time.Now()
        update := bson.M{
                "$set": bson.M{
                        "status":      "returned",
                        "return_date": now,
                        "updated_at":  now,
                },
        }

        // 计算罚款
        if now.After(record.DueDate) {
                days := int(now.Sub(record.DueDate).Hours() / 24)
                fine := float64(days) * 0.5 // 每天0.5元
                update["$set"].(bson.M)["fine"] = fine
        }

        _, err = bc.collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
        if err != nil {
                utils.ServerError(c, "归还失败")
                return
        }

        // 更新库存
        var stock models.Stock
        bc.mysqlDB.Where("book_id = ?", record.BookID).First(&stock)
        bc.mysqlDB.Model(&stock).Updates(map[string]interface{}{
                "available_quantity": stock.AvailableQty + 1,
                "borrowed_quantity":  stock.BorrowedQty - 1,
        })

        utils.SuccessWithMsg(c, "归还成功", nil)
}

// RenewBook 续借
func (bc *BorrowController) RenewBook(c *gin.Context) {
        var req struct {
                ID string `json:"id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        objID, err := primitive.ObjectIDFromHex(req.ID)
        if err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        var record models.BorrowRecord
        err = bc.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&record)
        if err != nil {
                utils.NotFound(c, "借阅记录不存在")
                return
        }

        if record.Status != "borrowed" {
                utils.Error(c, 400, "当前状态不可续借")
                return
        }

        if record.RenewCount >= 2 {
                utils.Error(c, 400, "已达最大续借次数")
                return
        }

        // 续借14天
        newDueDate := record.DueDate.AddDate(0, 0, 14)
        now := time.Now()

        update := bson.M{
                "$set": bson.M{
                        "due_date":    newDueDate,
                        "renew_count": record.RenewCount + 1,
                        "updated_at":  now,
                },
        }

        _, err = bc.collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
        if err != nil {
                utils.ServerError(c, "续借失败")
                return
        }

        utils.SuccessWithMsg(c, "续借成功", gin.H{
                "new_due_date": newDueDate,
        })
}

// GetOverdueRecords 获取逾期记录
func (bc *BorrowController) GetOverdueRecords(c *gin.Context) {
        var req struct {
                Page     int `json:"page"`
                PageSize int `json:"page_size"`
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

        filter := bson.M{
                "status":    "borrowed",
                "due_date": bson.M{"$lt": time.Now()},
        }

        skip := int64((req.Page - 1) * req.PageSize)
        limit := int64(req.PageSize)

        total, _ := bc.collection.CountDocuments(context.Background(), filter)

        opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "due_date", Value: 1}})
        cursor, err := bc.collection.Find(context.Background(), filter, opts)
        if err != nil {
                utils.ServerError(c, "查询失败")
                return
        }
        defer cursor.Close(context.Background())

        var records []models.BorrowRecord
        cursor.All(context.Background(), &records)

        utils.PageSuccess(c, records, total, req.Page, req.PageSize)
}

// GetUserBorrowStats 获取用户借阅统计
func (bc *BorrowController) GetUserBorrowStats(c *gin.Context) {
        userID := c.GetUint("user_id")

        // 当前借阅数
        borrowedCount, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "user_id": userID,
                "status":  "borrowed",
        })

        // 历史借阅数
        totalCount, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "user_id": userID,
        })

        // 逾期数
        overdueCount, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "user_id": userID,
                "status":  "overdue",
        })

        utils.Success(c, gin.H{
                "current_borrowed": borrowedCount,
                "total_borrowed":   totalCount,
                "overdue_count":    overdueCount,
        })
}

// GetBorrowRules 获取借阅规则
func (bc *BorrowController) GetBorrowRules(c *gin.Context) {
        collection := bc.mongoDB.Collection("borrow_rules")

        cursor, err := collection.Find(context.Background(), bson.M{})
        if err != nil {
                utils.ServerError(c, "查询失败")
                return
        }
        defer cursor.Close(context.Background())

        var rules []models.BorrowRule
        if err := cursor.All(context.Background(), &rules); err != nil {
                utils.ServerError(c, "解析数据失败")
                return
        }

        utils.Success(c, rules)
}

// GetBorrowStats 获取借阅统计（管理员）
func (bc *BorrowController) GetBorrowStats(c *gin.Context) {
        var req struct {
                StartDate string `json:"start_date"`
                EndDate   string `json:"end_date"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                req.StartDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
                req.EndDate = time.Now().Format("2006-01-02")
        }

        // 总借阅数
        totalBorrowed, _ := bc.collection.CountDocuments(context.Background(), bson.M{})

        // 当前借阅数
        currentBorrowed, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "status": "borrowed",
        })

        // 逾期数
        overdueCount, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "status": "overdue",
        })

        // 已归还数
        returnedCount, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "status": "returned",
        })

        // 今日借阅数
        todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
        todayBorrowed, _ := bc.collection.CountDocuments(context.Background(), bson.M{
                "borrow_date": bson.M{"$gte": todayStart},
        })

        utils.Success(c, gin.H{
                "total_borrowed":   totalBorrowed,
                "current_borrowed": currentBorrowed,
                "overdue_count":    overdueCount,
                "returned_count":   returnedCount,
                "today_borrowed":   todayBorrowed,
        })
}

// PayFine 支付罚款
func (bc *BorrowController) PayFine(c *gin.Context) {
        var req struct {
                ID string `json:"id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        objID, err := primitive.ObjectIDFromHex(req.ID)
        if err != nil {
                utils.BadRequest(c, "参数错误")
                return
        }

        var record models.BorrowRecord
        err = bc.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&record)
        if err != nil {
                utils.NotFound(c, "借阅记录不存在")
                return
        }

        if record.Fine <= 0 {
                utils.Error(c, 400, "无罚款需要支付")
                return
        }

        if record.FinePaid {
                utils.Error(c, 400, "罚款已支付")
                return
        }

        now := time.Now()
        update := bson.M{
                "$set": bson.M{
                        "fine_paid":  true,
                        "updated_at": now,
                },
        }

        _, err = bc.collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
        if err != nil {
                utils.ServerError(c, "支付失败")
                return
        }

        utils.SuccessWithMsg(c, "支付成功", nil)
}
