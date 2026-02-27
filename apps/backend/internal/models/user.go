package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Email       string         `gorm:"size:100;uniqueIndex" json:"email"`
	Phone       string         `gorm:"size:20" json:"phone"`
	RealName    string         `gorm:"size:50" json:"real_name"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Status      int8           `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	LastLoginAt *time.Time     `json:"last_login_at"`
	LastLoginIP string         `gorm:"size:45" json:"last_login_ip"`
	Roles       []Role         `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int8           `gorm:"default:1" json:"status"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Resource    string         `gorm:"size:50;not null" json:"resource"`
	Action      string         `gorm:"size:20;not null" json:"action"`
	ParentID    uint           `gorm:"default:0" json:"parent_id"`
	Description string         `gorm:"size:255" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
