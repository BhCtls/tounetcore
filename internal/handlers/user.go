package handlers

import (
	"net/http"
	"time"
	"tounetcore/internal/auth"
	"tounetcore/internal/config"
	"tounetcore/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{db: db, cfg: cfg}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Phone         string `json:"phone"`
	PushDeerToken string `json:"pushdeer_token"`
	InviteCode    string `json:"invite_code" binding:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	Phone         string `json:"phone"`
	PushDeerToken string `json:"pushdeer_token"`
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validate invite code
	var inviteCode models.InviteCode
	if err := h.db.Where("code = ? AND code_user_id IS NULL", req.InviteCode).First(&inviteCode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid or used invite code",
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

	// Create user
	user := models.User{
		Username:      req.Username,
		PasswordHash:  hashedPassword,
		Phone:         req.Phone,
		PushDeerToken: req.PushDeerToken,
		Status:        models.StatusUser,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to create user",
		})
		return
	}

	// Mark invite code as used
	now := time.Now()
	inviteCode.CodeUserID = &user.ID
	inviteCode.UsedAt = &now
	h.db.Save(&inviteCode)

	// Generate JWT token
	token, err := auth.GenerateJWT(&user, h.cfg.JWTSecret, h.cfg.JWTExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"user_id": user.ID,
			"token":   token,
		},
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "invalid credentials",
		})
		return
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "invalid credentials",
		})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	h.db.Save(&user)

	// Generate JWT token
	token, err := auth.GenerateJWT(&user, h.cfg.JWTSecret, h.cfg.JWTExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"token": token,
		},
	})
}

// GetUserInfo returns current user information
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"status":     user.Status,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
			"last_login": user.LastLogin,
		},
	})
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UpdateUserRequest
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

// ListAllowedApps returns user's allowed applications
func (h *UserHandler) ListAllowedApps(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userStatus, _ := c.Get("user_status")

	// Get all apps that match user's permission level
	var apps []models.App
	query := h.db.Where("is_active = ?", true)

	// Filter based on user status
	switch userStatus {
	case models.StatusAdmin:
		// Admin can see all apps
	case models.StatusTrusted:
		query = query.Where("required_permission_level IN (?)", []models.UserStatus{models.StatusDisabledUser, models.StatusUser, models.StatusTrusted})
	case models.StatusUser:
		query = query.Where("required_permission_level IN (?)", []models.UserStatus{models.StatusDisabledUser, models.StatusUser})
	case models.StatusDisabledUser:
		query = query.Where("required_permission_level = ?", models.StatusDisabledUser)
	}

	if err := query.Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to fetch apps",
		})
		return
	}

	// Get user-specific app permissions
	var userAllowedApps []models.UserAllowedApp
	h.db.Where("user_id = ?", userID).Find(&userAllowedApps)

	// Create a map for quick lookup
	userAppMap := make(map[string]models.UserAllowedApp)
	for _, ua := range userAllowedApps {
		userAppMap[ua.AppID] = ua
	}

	// Build response
	var result []gin.H
	for _, app := range apps {
		appData := gin.H{
			"app_id":      app.AppID,
			"name":        app.Name,
			"url":         app.URL,
			"enabled":     true,
			"valid_until": nil,
		}

		// Check user-specific permissions
		if userApp, exists := userAppMap[app.AppID]; exists {
			appData["enabled"] = userApp.Enabled
			appData["valid_until"] = userApp.ValidUntil
		}

		result = append(result, appData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}
