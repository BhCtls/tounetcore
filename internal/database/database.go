package database

import (
	"strings"
	"tounetcore/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if strings.HasPrefix(databaseURL, "sqlite://") {
		dbPath := strings.TrimPrefix(databaseURL, "sqlite://")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		// Default to SQLite if no scheme provided
		db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.InviteCode{},
		&models.NKey{},
		&models.App{},
		&models.UserAllowedApp{},
		&models.AuditLog{},
	)
}

// SeedData seeds initial data into the database
func SeedData(db *gorm.DB) error {
	// Create default apps
	apps := []models.App{
		{
			AppID:                   "searchall",
			Name:                    "Search All",
			Description:             "Global search functionality",
			RequiredPermissionLevel: models.StatusUser,
			IsActive:                true,
		},
		{
			AppID:                   "segaasstes",
			Name:                    "Sega Assets",
			Description:             "Sega assets management",
			RequiredPermissionLevel: models.StatusUser,
			IsActive:                true,
		},
		{
			AppID:                   "dxprender",
			Name:                    "DXP Render",
			Description:             "DXP rendering service",
			RequiredPermissionLevel: models.StatusUser,
			IsActive:                true,
		},
		{
			AppID:                   "CardPreview",
			Name:                    "Card Preview",
			Description:             "Card preview functionality",
			RequiredPermissionLevel: models.StatusUser,
			IsActive:                true,
		},
		{
			AppID:                   "livecontent_basic",
			Name:                    "Live Content Basic",
			Description:             "Basic live content access",
			RequiredPermissionLevel: models.StatusUser,
			IsActive:                true,
		},
		{
			AppID:                   "advanced_analytics",
			Name:                    "Advanced Analytics",
			Description:             "Advanced analytics and reporting tools",
			RequiredPermissionLevel: models.StatusTrusted,
			IsActive:                true,
		},
		{
			AppID:                   "livecontent_admin",
			Name:                    "Live Content Admin",
			Description:             "Administrative live content access",
			RequiredPermissionLevel: models.StatusAdmin,
			IsActive:                true,
		},
	}

	for _, app := range apps {
		// Check if app exists
		var existingApp models.App
		if err := db.Where("app_id = ?", app.AppID).First(&existingApp).Error; err == nil {
			continue // App already exists
		}

		// Generate a default secret key (should be properly generated in production)
		app.SecretKey = generateSecretKey(app.AppID)

		if err := db.Create(&app).Error; err != nil {
			return err
		}
	}

	return nil
}

// generateSecretKey generates a basic secret key for an app
// In production, this should use proper cryptographic methods
func generateSecretKey(appID string) string {
	// This is a placeholder - in production, use proper key generation
	return "secret_" + appID + "_key_change_in_production"
}
