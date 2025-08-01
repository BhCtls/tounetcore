package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	Phone         string         `json:"phone"`
	PushDeerToken string         `json:"pushdeer_token"`
	Status        UserStatus     `gorm:"type:varchar(20);default:user" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	LastLogin     *time.Time     `json:"last_login"`

	// Relationships
	AllowedApps []UserAllowedApp `gorm:"foreignKey:UserID" json:"allowed_apps,omitempty"`
	NKeys       []NKey           `gorm:"foreignKey:UserID" json:"nkeys,omitempty"`
}

// UserStatus represents the status/role of a user
type UserStatus string

const (
	StatusAdmin        UserStatus = "admin"
	StatusTrusted      UserStatus = "trusted"
	StatusUser         UserStatus = "user"
	StatusDisabledUser UserStatus = "disableduser"
)

// GetPermissionLevel returns the numeric permission level for comparison
// Higher numbers indicate higher permissions
func (us UserStatus) GetPermissionLevel() int {
	switch us {
	case StatusAdmin:
		return 4
	case StatusTrusted:
		return 3
	case StatusUser:
		return 2
	case StatusDisabledUser:
		return 1
	default:
		return 0
	}
}

// HasPermission checks if the current user status has permission for the required level
func (us UserStatus) HasPermission(required UserStatus) bool {
	return us.GetPermissionLevel() >= required.GetPermissionLevel()
}

// InviteCode represents an invitation code
type InviteCode struct {
	Code       string     `gorm:"primaryKey" json:"code"`
	Time       time.Time  `gorm:"not null" json:"time"`
	CodeUserID *uint      `json:"code_user_id"`
	UsedAt     *time.Time `json:"used_at"`

	// Relationships
	User *User `gorm:"foreignKey:CodeUserID" json:"user,omitempty"`
}

// NKey represents a generated authorization key
type NKey struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	KeyValue     string     `gorm:"uniqueIndex;not null" json:"key_value"`
	UserID       uint       `gorm:"not null" json:"user_id"`
	AppIDs       string     `gorm:"type:text" json:"app_ids"` // JSON array of app IDs
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	FirstUsedAt  *time.Time `json:"first_used_at"`
	FirstUsedApp string     `json:"first_used_app"`
	IsUsed       bool       `gorm:"default:false" json:"is_used"`
	CreatedAt    time.Time  `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// App represents an application that can be authorized
type App struct {
	AppID                   string     `gorm:"primaryKey;type:text" json:"app_id"`
	SecretKey               string     `gorm:"not null" json:"-"` // Encrypted secret key
	Name                    string     `gorm:"not null" json:"name"`
	Description             string     `json:"description"`
	URL                     string     `json:"url"` // Application URL
	RequiredPermissionLevel UserStatus `gorm:"type:varchar(20);default:user" json:"required_permission_level"`
	IsActive                bool       `gorm:"default:true" json:"is_active"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// UserAllowedApp represents the relationship between users and allowed apps
type UserAllowedApp struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	AppID       string     `gorm:"not null;type:text" json:"app_id"`
	Enabled     bool       `gorm:"default:true" json:"enabled"`
	ValidUntil  *time.Time `json:"valid_until"`
	CustomLimit string     `gorm:"type:text" json:"custom_limit"` // JSON format
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AuditLog represents system audit logs
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActionType string    `gorm:"not null" json:"action_type"`
	TargetType string    `gorm:"not null" json:"target_type"`
	TargetID   string    `json:"target_id"`
	OperatorID uint      `json:"operator_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Details    string    `gorm:"type:text" json:"details"` // JSON format
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Operator *User `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}
