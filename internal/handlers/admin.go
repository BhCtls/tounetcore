package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tounetcore/internal/auth"
	"tounetcore/internal/config"
	"tounetcore/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAdminHandler(db *gorm.DB, cfg *config.Config) *AdminHandler {
	return &AdminHandler{db: db, cfg: cfg}
}

// CreateUserRequest represents admin user creation request
type CreateUserRequest struct {
	Username      string            `json:"username" binding:"required"`
	Password      string            `json:"password" binding:"required"`
	Status        models.UserStatus `json:"status"`
	Phone         string            `json:"phone"`
	PushDeerToken string            `json:"pushdeer_token"`
}

// CreateAppRequest represents app creation request
type CreateAppRequest struct {
	AppID                   string            `json:"app_id" binding:"required"`
	Name                    string            `json:"name" binding:"required"`
	Description             string            `json:"description"`
	URL                     string            `json:"url"`
	RequiredPermissionLevel models.UserStatus `json:"required_permission_level"`
	IsActive                bool              `json:"is_active"`
}

// UpdateAppRequest represents app update request
type UpdateAppRequest struct {
	Name                    string            `json:"name"`
	Description             string            `json:"description"`
	URL                     string            `json:"url"`
	SecretKey               string            `json:"secret_key"`
	RequiredPermissionLevel models.UserStatus `json:"required_permission_level"`
	IsActive                *bool             `json:"is_active"`
}

// AdminUpdateUserRequest represents admin user update request
type AdminUpdateUserRequest struct {
	Username      string            `json:"username"`
	Password      string            `json:"password"`
	Status        models.UserStatus `json:"status"`
	Phone         string            `json:"phone"`
	PushDeerToken string            `json:"pushdeer_token"`
}

// CreateUser creates a new user (admin only)
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "username already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to hash password",
		})
		return
	}

	// Force default status to user for security - admin/trusted status must be set via update
	req.Status = models.StatusUser

	// Create user
	user := models.User{
		Username:      req.Username,
		PasswordHash:  hashedPassword,
		Phone:         req.Phone,
		PushDeerToken: req.PushDeerToken,
		Status:        req.Status,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"user_id": user.ID,
		},
	})
}

// ListUsers returns all users with pagination
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size

	var users []models.User
	var total int64

	// Get total count
	h.db.Model(&models.User{}).Count(&total)

	// Get users with pagination
	if err := h.db.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to fetch users",
		})
		return
	}

	// Build response
	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"status":     user.Status,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
			"last_login": user.LastLogin,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total": total,
			"users": userList,
		},
	})
}

// ListApps returns all applications
func (h *AdminHandler) ListApps(c *gin.Context) {
	var apps []models.App
	if err := h.db.Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to fetch apps",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    apps,
	})
}

// CreateApp creates a new application
func (h *AdminHandler) CreateApp(c *gin.Context) {
	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	// Check if app_id already exists
	var existingApp models.App
	if err := h.db.Where("app_id = ?", req.AppID).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "app_id already exists",
		})
		return
	}

	// Generate secret key
	secretKey, err := auth.GenerateSecretKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to generate secret key",
		})
		return
	}

	// Set defaults
	if req.RequiredPermissionLevel == "" {
		req.RequiredPermissionLevel = models.StatusUser
	}

	// Create app
	app := models.App{
		AppID:                   req.AppID,
		SecretKey:               secretKey,
		Name:                    req.Name,
		Description:             req.Description,
		URL:                     req.URL,
		RequiredPermissionLevel: req.RequiredPermissionLevel,
		IsActive:                req.IsActive,
	}

	if err := h.db.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to create app",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"app_id":                    app.AppID,
			"name":                      app.Name,
			"url":                       app.URL,
			"secret_key":                app.SecretKey,
			"required_permission_level": app.RequiredPermissionLevel,
			"is_active":                 app.IsActive,
		},
	})
}

// UpdateApp updates an existing application
func (h *AdminHandler) UpdateApp(c *gin.Context) {
	appID := c.Param("app_id")

	var req UpdateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	var app models.App
	if err := h.db.Where("app_id = ?", appID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "app not found",
		})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	if req.URL != "" {
		app.URL = req.URL
	}
	if req.SecretKey != "" {
		app.SecretKey = req.SecretKey
	}
	if req.RequiredPermissionLevel != "" {
		app.RequiredPermissionLevel = req.RequiredPermissionLevel
	}
	if req.IsActive != nil {
		app.IsActive = *req.IsActive
	}

	if err := h.db.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to update app",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// DeleteApp deletes an application (admin only)
func (h *AdminHandler) DeleteApp(c *gin.Context) {
	appID := c.Param("app_id")
	operatorID, _ := c.Get("user_id")

	var app models.App
	if err := h.db.Where("app_id = ?", appID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "app not found",
		})
		return
	}

	// Delete related records first
	h.db.Where("app_id = ?", appID).Delete(&models.UserAllowedApp{})

	// Delete the app
	if err := h.db.Delete(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to delete app",
		})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		ActionType: "DELETE_APP",
		TargetType: "APP",
		TargetID:   appID,
		OperatorID: operatorID.(uint),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Details:    fmt.Sprintf("Deleted app: %s", app.Name),
	}
	h.db.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// ToggleAppStatus toggles application active status (admin only)
func (h *AdminHandler) ToggleAppStatus(c *gin.Context) {
	appID := c.Param("app_id")
	operatorID, _ := c.Get("user_id")

	var app models.App
	if err := h.db.Where("app_id = ?", appID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "app not found",
		})
		return
	}

	// Toggle the status
	app.IsActive = !app.IsActive

	if err := h.db.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to update app status",
		})
		return
	}

	// Create audit log
	status := "DISABLED"
	if app.IsActive {
		status = "ENABLED"
	}
	auditLog := models.AuditLog{
		ActionType: "TOGGLE_APP_STATUS",
		TargetType: "APP",
		TargetID:   appID,
		OperatorID: operatorID.(uint),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Details:    fmt.Sprintf("App %s status changed to: %s", app.Name, status),
	}
	h.db.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"app_id":    app.AppID,
			"is_active": app.IsActive,
		},
	})
}

// GenerateInviteCode generates a new invite code
func (h *AdminHandler) GenerateInviteCode(c *gin.Context) {
	code, err := auth.GenerateInviteCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to generate invite code",
		})
		return
	}

	inviteCode := models.InviteCode{
		Code: code,
		Time: time.Now(),
	}

	if err := h.db.Create(&inviteCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to create invite code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"invite_code": code,
		},
	})
}

// ViewAuditLogs returns audit logs with pagination
func (h *AdminHandler) ViewAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	offset := (page - 1) * size

	var logs []models.AuditLog
	var total int64

	// Get total count
	h.db.Model(&models.AuditLog{}).Count(&total)

	// Get logs with pagination, ordered by creation time desc
	if err := h.db.Preload("Operator").Offset(offset).Limit(size).Order("created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to fetch audit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total": total,
			"logs":  logs,
		},
	})
}

// ListInviteCodes returns all invite codes with pagination
func (h *AdminHandler) ListInviteCodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	offset := (page - 1) * size

	var inviteCodes []models.InviteCode
	var total int64

	// Get total count
	h.db.Model(&models.InviteCode{}).Count(&total)

	// Get invite codes with pagination, ordered by creation time desc
	if err := h.db.Preload("User").Offset(offset).Limit(size).Order("time DESC").Find(&inviteCodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to fetch invite codes",
		})
		return
	}

	// Build response
	var codeList []gin.H
	for _, inviteCode := range inviteCodes {
		codeData := gin.H{
			"code":         inviteCode.Code,
			"time":         inviteCode.Time,
			"code_user_id": inviteCode.CodeUserID,
			"used_at":      inviteCode.UsedAt,
			"used_by":      nil,
		}

		// Add user information if the code has been used
		if inviteCode.User != nil {
			codeData["used_by"] = gin.H{
				"id":       inviteCode.User.ID,
				"username": inviteCode.User.Username,
			}
		}

		codeList = append(codeList, codeData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"total":        total,
			"invite_codes": codeList,
		},
	})
}

// UpdateUser updates user information (admin only)
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("user_id")

	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}

	// Update fields if provided
	if req.Username != "" {
		// Check if new username already exists
		var existingUser models.User
		if err := h.db.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "username already exists",
			})
			return
		}
		user.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "failed to hash password",
			})
			return
		}
		user.PasswordHash = hashedPassword
	}
	if req.Status != "" {
		operatorID, _ := c.Get("user_id")

		// Prevent admin from demoting themselves
		if userID == fmt.Sprintf("%v", operatorID) && user.Status == models.StatusAdmin && req.Status != models.StatusAdmin {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "cannot demote your own admin privileges",
			})
			return
		}

		user.Status = req.Status
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.PushDeerToken != "" {
		user.PushDeerToken = req.PushDeerToken
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// DeleteUser deletes a user (admin only)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	operatorID, _ := c.Get("user_id")

	// Prevent admin from deleting themselves
	if userID == fmt.Sprintf("%v", operatorID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "cannot delete your own account",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}

	// Use soft delete
	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to delete user",
		})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		ActionType: "DELETE_USER",
		TargetType: "USER",
		TargetID:   userID,
		OperatorID: operatorID.(uint),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Details:    fmt.Sprintf("Deleted user: %s", user.Username),
	}
	h.db.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

// DeleteInviteCode deletes an invite code (admin only)
func (h *AdminHandler) DeleteInviteCode(c *gin.Context) {
	inviteCode := c.Param("invite_code")
	operatorID, _ := c.Get("user_id")

	var code models.InviteCode
	if err := h.db.Where("code = ?", inviteCode).First(&code).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "invite code not found",
		})
		return
	}

	// Check if code has been used
	if code.CodeUserID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "cannot delete used invite code",
		})
		return
	}

	if err := h.db.Delete(&code).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to delete invite code",
		})
		return
	}

	// Create audit log
	auditLog := models.AuditLog{
		ActionType: "DELETE_INVITE_CODE",
		TargetType: "INVITE_CODE",
		TargetID:   inviteCode,
		OperatorID: operatorID.(uint),
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Details:    fmt.Sprintf("Deleted invite code: %s", inviteCode),
	}
	h.db.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}
