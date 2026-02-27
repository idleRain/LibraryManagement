package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/library-system/backend/internal/config"
	"github.com/library-system/backend/internal/middleware"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserController(db *gorm.DB, cfg *config.Config) *UserController {
	return &UserController{db: db, cfg: cfg}
}

// Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := uc.db.Preload("Roles.Permissions").Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.Error(c, 400, "账户已被禁用")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	// 生成 Token
	token, err := middleware.GenerateToken(user.ID, user.Username, &uc.cfg.JWT)
	if err != nil {
		utils.ServerError(c, "生成Token失败")
		return
	}

	// 更新登录信息
	now := time.Now()
	uc.db.Model(&user).Updates(map[string]interface{}{
		"last_login_at": &now,
		"last_login_ip": c.ClientIP(),
	})

	// 提取权限码
	var permissions []string
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Code)
		}
	}

	utils.Success(c, gin.H{
		"token":       token,
		"user":        user,
		"permissions": permissions,
	})
}

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
		RealName string `json:"real_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 检查用户名是否存在
	var count int64
	uc.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		RealName: req.RealName,
		Status:   1,
	}

	if err := uc.db.Create(&user).Error; err != nil {
		utils.ServerError(c, "创建用户失败")
		return
	}

	utils.SuccessWithMsg(c, "注册成功", user)
}

// GetUsers 获取用户列表
func (uc *UserController) GetUsers(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Username string `json:"username"`
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

	var users []models.User
	var total int64

	query := uc.db.Model(&models.User{}).Preload("Roles")
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&users)

	utils.PageSuccess(c, users, total, req.Page, req.PageSize)
}

// GetUser 获取用户详情
func (uc *UserController) GetUser(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := uc.db.Preload("Roles.Permissions").First(&user, req.ID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		RealName string `json:"real_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 检查用户名是否存在
	var count int64
	uc.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		RealName: req.RealName,
		Status:   1,
	}

	if err := uc.db.Create(&user).Error; err != nil {
		utils.ServerError(c, "创建用户失败")
		return
	}

	utils.SuccessWithMsg(c, "创建成功", user)
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	var req struct {
		ID       int     `json:"id" binding:"required"`
		Email    *string `json:"email"`
		Phone    *string `json:"phone"`
		RealName *string `json:"real_name"`
		Status   *int8   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := uc.db.First(&user, req.ID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.RealName != nil {
		updates["real_name"] = *req.RealName
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := uc.db.Model(&user).Updates(updates).Error; err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMsg(c, "更新成功", user)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := uc.db.Delete(&models.User{}, req.ID).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// AssignRoles 分配角色
func (uc *UserController) AssignRoles(c *gin.Context) {
	var req struct {
		ID      int   `json:"id" binding:"required"`
		RoleIDs []int `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := uc.db.First(&user, req.ID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var roles []models.Role
	uc.db.Where("id IN ?", req.RoleIDs).Find(&roles)

	uc.db.Model(&user).Association("Roles").Replace(roles)

	utils.SuccessWithMsg(c, "分配角色成功", nil)
}

// ChangePassword 修改密码
func (uc *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := uc.db.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.Error(c, 400, "原密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	uc.db.Model(&user).Update("password", string(hashedPassword))

	utils.SuccessWithMsg(c, "密码修改成功", nil)
}
