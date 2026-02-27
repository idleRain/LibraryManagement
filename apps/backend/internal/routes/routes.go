package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/config"
	"github.com/library-system/backend/internal/controllers"
	"github.com/library-system/backend/internal/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB, mongoDB *mongo.Database, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())

	// 初始化控制器
	userCtrl := controllers.NewUserController(db, cfg)
	roleCtrl := controllers.NewRoleController(db)
	permCtrl := controllers.NewPermissionController(db)
	bookCtrl := controllers.NewBookController(db)
	stockCtrl := controllers.NewStockController(db)
	purchaseCtrl := controllers.NewPurchaseController(db)
	saleCtrl := controllers.NewSaleController(db)
	borrowCtrl := controllers.NewBorrowController(db, mongoDB)

	// API 路由组
	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/login", userCtrl.Login)
		api.POST("/register", userCtrl.Register)

		// 需要认证的接口
		auth := api.Group("")
		auth.Use(middleware.JWTAuth(&cfg.JWT))
		{
			// 用户管理
			auth.POST("/users/list", userCtrl.GetUsers)
			auth.POST("/users/detail", userCtrl.GetUser)
			auth.POST("/users/create", userCtrl.CreateUser)
			auth.POST("/users/update", userCtrl.UpdateUser)
			auth.POST("/users/delete", userCtrl.DeleteUser)
			auth.POST("/users/assign-roles", userCtrl.AssignRoles)
			auth.POST("/users/change-password", userCtrl.ChangePassword)

			// 角色管理
			auth.POST("/roles/list", roleCtrl.GetRoles)
			auth.POST("/roles/create", roleCtrl.CreateRole)
			auth.POST("/roles/update", roleCtrl.UpdateRole)
			auth.POST("/roles/delete", roleCtrl.DeleteRole)
			auth.POST("/roles/assign-permissions", roleCtrl.AssignPermissions)

			// 权限管理
			auth.POST("/permissions/list", permCtrl.GetPermissions)
			auth.POST("/permissions/create", permCtrl.CreatePermission)
			auth.POST("/permissions/update", permCtrl.UpdatePermission)
			auth.POST("/permissions/delete", permCtrl.DeletePermission)

			// 图书管理
			auth.POST("/books/list", bookCtrl.GetBooks)
			auth.POST("/books/detail", bookCtrl.GetBook)
			auth.POST("/books/create", bookCtrl.CreateBook)
			auth.POST("/books/update", bookCtrl.UpdateBook)
			auth.POST("/books/delete", bookCtrl.DeleteBook)
			auth.POST("/books/categories", bookCtrl.GetCategories)

			// 库存管理
			auth.POST("/stocks/list", stockCtrl.GetStocks)
			auth.POST("/stocks/detail", stockCtrl.GetStockDetail)
			auth.POST("/stocks/in", stockCtrl.StockIn)
			auth.POST("/stocks/out", stockCtrl.StockOut)
			auth.POST("/stocks/records", stockCtrl.GetStockRecords)
			auth.POST("/stocks/low", stockCtrl.GetLowStock)

			// 供应商管理
			auth.POST("/suppliers/list", purchaseCtrl.GetSuppliers)
			auth.POST("/suppliers/create", purchaseCtrl.CreateSupplier)
			auth.POST("/suppliers/update", purchaseCtrl.UpdateSupplier)
			auth.POST("/suppliers/delete", purchaseCtrl.DeleteSupplier)

			// 采购管理
			auth.POST("/purchases/list", purchaseCtrl.GetPurchaseOrders)
			auth.POST("/purchases/detail", purchaseCtrl.GetPurchaseOrder)
			auth.POST("/purchases/create", purchaseCtrl.CreatePurchaseOrder)
			auth.POST("/purchases/update-status", purchaseCtrl.UpdatePurchaseOrderStatus)
			auth.POST("/purchases/delete", purchaseCtrl.DeletePurchaseOrder)

			// 销售管理
			auth.POST("/sales/list", saleCtrl.GetSaleOrders)
			auth.POST("/sales/detail", saleCtrl.GetSaleOrder)
			auth.POST("/sales/create", saleCtrl.CreateSaleOrder)
			auth.POST("/sales/cancel", saleCtrl.CancelSaleOrder)
			auth.POST("/sales/stats", saleCtrl.GetSalesStats)

			// 购物车
			auth.POST("/cart/list", saleCtrl.GetCart)
			auth.POST("/cart/add", saleCtrl.AddToCart)
			auth.POST("/cart/remove", saleCtrl.RemoveFromCart)

			// 借阅管理
			auth.POST("/borrows/list", borrowCtrl.GetBorrowRecords)
			auth.POST("/borrows/detail", borrowCtrl.GetBorrowRecord)
			auth.POST("/borrows/borrow", borrowCtrl.BorrowBook)
			auth.POST("/borrows/return", borrowCtrl.ReturnBook)
			auth.POST("/borrows/renew", borrowCtrl.RenewBook)
			auth.POST("/borrows/overdue", borrowCtrl.GetOverdueRecords)
			auth.POST("/borrows/stats", borrowCtrl.GetUserBorrowStats)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Library System API is running",
		})
	})

	return r
}
