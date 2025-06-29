package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"tounetcore/internal/auth"
	"tounetcore/internal/config"
	"tounetcore/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NKeyHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewNKeyHandler(db *gorm.DB, cfg *config.Config) *NKeyHandler {
	return &NKeyHandler{db: db, cfg: cfg}
}

// ApplyNKeyRequest represents the request to generate an NKey
type ApplyNKeyRequest struct {
	Username []string `json:"username"`
	AppIDs   []string `json:"app_ids" binding:"required"`
}

// ValidateNKeyRequest represents the request to validate an NKey
type ValidateNKeyRequest struct {
	NKey  string `json:"nkey" binding:"required"`
	AppID string `json:"app_id" binding:"required"`
}

// ApplyNKey generates a new NKey for the user
func (h *NKeyHandler) ApplyNKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userStatus, _ := c.Get("user_status")

	var req ApplyNKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	// Get user information
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
		return
	}

	// Validate requested apps
	var validAppIDs []string
	for _, appID := range req.AppIDs {
		var app models.App
		if err := h.db.Where("app_id = ? AND is_active = ?", appID, true).First(&app).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "invalid app_id: " + appID,
			})
			return
		}

		// Check if user has permission for this app
		if !h.userHasAppPermission(userID.(uint), appID, userStatus.(models.UserStatus), app.RequiredPermissionLevel) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "no permission for app: " + appID,
			})
			return
		}

		validAppIDs = append(validAppIDs, appID)
	}

	// Generate NKey
	nkey, err := auth.GenerateNKey(userID.(uint), validAppIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to generate nkey",
		})
		return
	}

	// Store NKey in database
	appIDsJSON, _ := json.Marshal(validAppIDs)
	nkeyRecord := models.NKey{
		KeyValue:  nkey,
		UserID:    userID.(uint),
		AppIDs:    string(appIDsJSON),
		ExpiresAt: time.Now().Add(h.cfg.NKeyExpiration),
	}

	if err := h.db.Create(&nkeyRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to store nkey",
		})
		return
	}

	// Send push notification if PushDeer token is available
	if user.PushDeerToken != "" {
		go h.sendPushNotification(user.PushDeerToken, nkey)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"nkey":       nkey,
			"expires_in": int(h.cfg.NKeyExpiration.Seconds()),
			"apps":       validAppIDs,
		},
	})
}

// ValidateNKey validates an NKey for a specific app
func (h *NKeyHandler) ValidateNKey(c *gin.Context) {
	var req ValidateNKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
		})
		return
	}

	// Find NKey in database
	var nkey models.NKey
	if err := h.db.Preload("User").Where("key_value = ?", req.NKey).First(&nkey).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "invalid nkey",
		})
		return
	}

	// Check if NKey is expired
	if time.Now().After(nkey.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "expired nkey",
		})
		return
	}

	// Check if user has permission for the requested app
	var appIDs []string
	if err := json.Unmarshal([]byte(nkey.AppIDs), &appIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "invalid nkey data",
		})
		return
	}

	// Check if the requested app is in the allowed apps
	hasPermission := false
	for _, appID := range appIDs {
		if appID == req.AppID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "no permission for app",
		})
		return
	}

	// Update first use information if not already used
	if !nkey.IsUsed {
		now := time.Now()
		nkey.FirstUsedAt = &now
		nkey.FirstUsedApp = req.AppID
		nkey.IsUsed = true
		h.db.Save(&nkey)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"valid":     true,
			"username":  nkey.User.Username,
			"user_role": nkey.User.Status,
		},
	})
}

// userHasAppPermission checks if a user has permission for a specific app
func (h *NKeyHandler) userHasAppPermission(userID uint, appID string, userStatus models.UserStatus, requiredLevel models.UserStatus) bool {
	// Check basic permission level
	switch requiredLevel {
	case models.StatusAdmin:
		if userStatus != models.StatusAdmin {
			return false
		}
	case models.StatusUser:
		if userStatus == models.StatusDisabledUser {
			return false
		}
	case models.StatusDisabledUser:
		// All users can access disabled user level apps
	}

	// Check user-specific app permissions
	var userApp models.UserAllowedApp
	if err := h.db.Where("user_id = ? AND app_id = ?", userID, appID).First(&userApp).Error; err == nil {
		// User has specific permission settings
		if !userApp.Enabled {
			return false
		}
		if userApp.ValidUntil != nil && time.Now().After(*userApp.ValidUntil) {
			return false
		}
	}

	return true
}

// sendPushNotification sends a push notification via PushDeer
func (h *NKeyHandler) sendPushNotification(token, nkey string) {
	// Implementation for PushDeer notification
	// This is a placeholder - implement actual HTTP request to PushDeer API
	// Example: POST to https://api2.pushdeer.com/message/push
	// with pushkey=token and text=nkey
}
