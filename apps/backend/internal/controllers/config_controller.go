package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/library-system/backend/internal/models"
	"github.com/library-system/backend/internal/utils"
	"gorm.io/gorm"
)

type ConfigController struct {
	db *gorm.DB
}

func NewConfigController(db *gorm.DB) *ConfigController {
	return &ConfigController{db: db}
}

// GetConfigs 获取系统配置列表
func (cc *ConfigController) GetConfigs(c *gin.Context) {
	var configs []models.SystemConfig
	cc.db.Order("id ASC").Find(&configs)

	// 转换为 map
	configMap := make(map[string]interface{})
	for _, cfg := range configs {
		configMap[cfg.ConfigKey] = cfg.ConfigValue
	}

	utils.Success(c, configMap)
}

// GetConfig 获取单个配置
func (cc *ConfigController) GetConfig(c *gin.Context) {
	var req struct {
		Key string `json:"key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var config models.SystemConfig
	if err := cc.db.Where("config_key = ?", req.Key).First(&config).Error; err != nil {
		utils.NotFound(c, "配置不存在")
		return
	}

	utils.Success(c, config)
}

// UpdateConfig 更新配置
func (cc *ConfigController) UpdateConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var config models.SystemConfig
	if err := cc.db.Where("config_key = ?", req.Key).First(&config).Error; err != nil {
		// 配置不存在，创建新配置
		config = models.SystemConfig{
			ConfigKey:   req.Key,
			ConfigValue: req.Value,
		}
		if err := cc.db.Create(&config).Error; err != nil {
			utils.ServerError(c, "创建配置失败")
			return
		}
	} else {
		// 更新配置
		config.ConfigValue = req.Value
		if err := cc.db.Save(&config).Error; err != nil {
			utils.ServerError(c, "更新配置失败")
			return
		}
	}

	utils.SuccessWithMsg(c, "更新成功", config)
}

// BatchUpdateConfigs 批量更新配置
func (cc *ConfigController) BatchUpdateConfigs(c *gin.Context) {
	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tx := cc.db.Begin()

	for key, value := range req {
		var config models.SystemConfig
		if err := tx.Where("config_key = ?", key).First(&config).Error; err != nil {
			// 创建新配置
			config = models.SystemConfig{
				ConfigKey:   key,
				ConfigValue: value,
			}
			tx.Create(&config)
		} else {
			// 更新配置
			config.ConfigValue = value
			tx.Save(&config)
		}
	}

	tx.Commit()

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// GetBorrowRules 获取借阅规则配置
func (cc *ConfigController) GetBorrowRules(c *gin.Context) {
	rules := map[string]interface{}{
		"borrow_days":      30,
		"max_borrow_books": 5,
		"max_renew_times":  2,
		"fine_per_day":     0.5,
	}

	// 从数据库获取实际配置
	var configs []models.SystemConfig
	cc.db.Where("config_key IN ?", []string{
		"borrow_days", "max_borrow_books", "max_renew_times", "fine_per_day",
	}).Find(&configs)

	for _, cfg := range configs {
		rules[cfg.ConfigKey] = cfg.ConfigValue
	}

	utils.Success(c, rules)
}

// UpdateBorrowRules 更新借阅规则配置
func (cc *ConfigController) UpdateBorrowRules(c *gin.Context) {
	var req struct {
		BorrowDays     *int     `json:"borrow_days"`
		MaxBorrowBooks *int     `json:"max_borrow_books"`
		MaxRenewTimes  *int     `json:"max_renew_times"`
		FinePerDay     *float64 `json:"fine_per_day"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tx := cc.db.Begin()

	if req.BorrowDays != nil {
		cc.updateOrCreateConfig(tx, "borrow_days", string(rune(*req.BorrowDays)))
	}
	if req.MaxBorrowBooks != nil {
		cc.updateOrCreateConfig(tx, "max_borrow_books", string(rune(*req.MaxBorrowBooks)))
	}
	if req.MaxRenewTimes != nil {
		cc.updateOrCreateConfig(tx, "max_renew_times", string(rune(*req.MaxRenewTimes)))
	}
	if req.FinePerDay != nil {
		cc.updateOrCreateConfig(tx, "fine_per_day", string(rune(int(*req.FinePerDay*100))))
	}

	tx.Commit()

	utils.SuccessWithMsg(c, "更新成功", nil)
}

func (cc *ConfigController) updateOrCreateConfig(tx *gorm.DB, key, value string) {
	var config models.SystemConfig
	if err := tx.Where("config_key = ?", key).First(&config).Error; err != nil {
		config = models.SystemConfig{
			ConfigKey:   key,
			ConfigValue: value,
		}
		tx.Create(&config)
	} else {
		config.ConfigValue = value
		tx.Save(&config)
	}
}

// GetSiteInfo 获取站点信息
func (cc *ConfigController) GetSiteInfo(c *gin.Context) {
	info := map[string]interface{}{
		"site_name":    "图书管理系统",
		"site_logo":    "",
		"site_footer":  "",
		"contact_email": "",
		"contact_phone": "",
	}

	var configs []models.SystemConfig
	cc.db.Where("config_key IN ?", []string{
		"site_name", "site_logo", "site_footer", "contact_email", "contact_phone",
	}).Find(&configs)

	for _, cfg := range configs {
		info[cfg.ConfigKey] = cfg.ConfigValue
	}

	utils.Success(c, info)
}

// UpdateSiteInfo 更新站点信息
func (cc *ConfigController) UpdateSiteInfo(c *gin.Context) {
	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	allowedFields := map[string]bool{
		"site_name":     true,
		"site_logo":     true,
		"site_footer":   true,
		"contact_email": true,
		"contact_phone": true,
	}

	tx := cc.db.Begin()

	for key, value := range req {
		if !allowedFields[key] {
			continue
		}
		var config models.SystemConfig
		if err := tx.Where("config_key = ?", key).First(&config).Error; err != nil {
			config = models.SystemConfig{
				ConfigKey:   key,
				ConfigValue: value,
			}
			tx.Create(&config)
		} else {
			config.ConfigValue = value
			tx.Save(&config)
		}
	}

	tx.Commit()

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// ResetConfigs 重置配置为默认值
func (cc *ConfigController) ResetConfigs(c *gin.Context) {
	defaultConfigs := map[string]string{
		"site_name":        "图书管理系统",
		"borrow_days":      "30",
		"max_borrow_books": "5",
		"max_renew_times":  "2",
		"fine_per_day":     "0.5",
		"low_stock_threshold": "5",
	}

	tx := cc.db.Begin()

	for key, value := range defaultConfigs {
		var config models.SystemConfig
		if err := tx.Where("config_key = ?", key).First(&config).Error; err != nil {
			config = models.SystemConfig{
				ConfigKey:   key,
				ConfigValue: value,
			}
			tx.Create(&config)
		} else {
			config.ConfigValue = value
			tx.Save(&config)
		}
	}

	tx.Commit()

	utils.SuccessWithMsg(c, "重置成功", defaultConfigs)
}
