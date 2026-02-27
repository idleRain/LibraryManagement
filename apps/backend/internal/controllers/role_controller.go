package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type RoleController struct {
	db *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{db: db}
}

// GetRoles 获取角色列表
func (rc *RoleController) GetRoles(c *gin.Context) {
	var roles []models.Role
	rc.db.Preload("Permissions").Find(&roles)
	utils.Success(c, roles)
}

// CreateRole 创建角色
func (rc *RoleController) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	role := models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := rc.db.Create(&role).Error; err != nil {
		utils.ServerError(c, "创建角色失败")
		return
	}

	utils.SuccessWithMsg(c, "创建成功", role)
}

// UpdateRole 更新角色
func (rc *RoleController) UpdateRole(c *gin.Context) {
	var req struct {
		ID          int    `json:"id" binding:"required"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := rc.db.First(&role, req.ID).Error; err != nil {
		utils.NotFound(c, "角色不存在")
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}

	rc.db.Model(&role).Updates(updates)
	utils.SuccessWithMsg(c, "更新成功", role)
}

// DeleteRole 删除角色
func (rc *RoleController) DeleteRole(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := rc.db.Delete(&models.Role{}, req.ID).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// AssignPermissions 分配权限
func (rc *RoleController) AssignPermissions(c *gin.Context) {
	var req struct {
		ID            int   `json:"id" binding:"required"`
		PermissionIDs []int `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := rc.db.First(&role, req.ID).Error; err != nil {
		utils.NotFound(c, "角色不存在")
		return
	}

	var permissions []models.Permission
	rc.db.Where("id IN ?", req.PermissionIDs).Find(&permissions)

	rc.db.Model(&role).Association("Permissions").Replace(permissions)

	utils.SuccessWithMsg(c, "分配权限成功", nil)
}
