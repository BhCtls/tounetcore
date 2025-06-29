package api

import (
	"tounetcore/internal/config"
	"tounetcore/internal/handlers"
	"tounetcore/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Add middleware
	router.Use(middleware.CORSMiddleware())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db, cfg)
	nkeyHandler := handlers.NewNKeyHandler(db, cfg)
	adminHandler := handlers.NewAdminHandler(db, cfg)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/nkey/validate", nkeyHandler.ValidateNKey)

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/me", userHandler.GetUserInfo)
				user.PUT("/me", userHandler.UpdateUser)
				user.GET("/apps", userHandler.ListAllowedApps)
			}

			// NKey routes
			nkey := protected.Group("/nkey")
			{
				nkey.POST("/generate", nkeyHandler.ApplyNKey)
			}

			// Admin routes (require admin role)
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// User management
				admin.POST("/users", adminHandler.CreateUser)
				admin.GET("/users", adminHandler.ListUsers)
				admin.POST("/users/:user_id/update", adminHandler.UpdateUser)
				admin.POST("/users/:user_id/delete", adminHandler.DeleteUser)

				// Invite code management
				admin.POST("/invite-codes", adminHandler.GenerateInviteCode)
				admin.GET("/invite-codes", adminHandler.ListInviteCodes)
				admin.POST("/invite-codes/:invite_code/delete", adminHandler.DeleteInviteCode)

				// App management
				admin.GET("/apps", adminHandler.ListApps)
				admin.POST("/apps", adminHandler.CreateApp)
				admin.POST("/apps/:app_id/update", adminHandler.UpdateApp)
				admin.POST("/apps/:app_id/delete", adminHandler.DeleteApp)
				admin.POST("/apps/:app_id/toggle", adminHandler.ToggleAppStatus)

				// Audit logs
				admin.GET("/logs", adminHandler.ViewAuditLogs)
			}
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "tounetcore",
		})
	})
}
