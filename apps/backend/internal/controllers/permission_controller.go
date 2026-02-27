package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type PermissionController struct {
	db *gorm.DB
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{db: db}
}

// GetPermissions 获取权限列表
func (pc *PermissionController) GetPermissions(c *gin.Context) {
	var permissions []models.Permission
	pc.db.Order("sort_order ASC").Find(&permissions)
	utils.Success(c, permissions)
}

// CreatePermission 创建权限
func (pc *PermissionController) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Resource    string `json:"resource" binding:"required"`
		Action      string `json:"action" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	permission := models.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := pc.db.Create(&permission).Error; err != nil {
		utils.ServerError(c, "创建权限失败")
		return
	}

	utils.SuccessWithMsg(c, "创建成功", permission)
}

// UpdatePermission 更新权限
func (pc *PermissionController) UpdatePermission(c *gin.Context) {
	var req struct {
		ID          int    `json:"id" binding:"required"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var permission models.Permission
	if err := pc.db.First(&permission, req.ID).Error; err != nil {
		utils.NotFound(c, "权限不存在")
		return
	}

	pc.db.Model(&permission).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	})
	utils.SuccessWithMsg(c, "更新成功", permission)
}

// DeletePermission 删除权限
func (pc *PermissionController) DeletePermission(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := pc.db.Delete(&models.Permission{}, req.ID).Error; err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}
