package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"tounetcore/internal/auth"
	"tounetcore/internal/config"
	"tounetcore/internal/database"
	"tounetcore/internal/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/seed/main.go <command>")
		fmt.Println("Commands:")
		fmt.Println("  apps      - Seed default applications")
		fmt.Println("  admin     - Create admin user")
		fmt.Println("  invite    - Generate invite codes")
		os.Exit(1)
	}

	command := os.Args[1]

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	switch command {
	case "apps":
		if err := database.SeedData(db); err != nil {
			log.Fatal("Failed to seed apps:", err)
		}
		fmt.Println("✅ Default applications seeded successfully")

	case "admin":
		// Create admin user
		hashedPassword, err := auth.HashPassword("admin123")
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		admin := models.User{
			Username:      "admin",
			PasswordHash:  hashedPassword,
			Status:        models.StatusAdmin,
			PushDeerToken: "PUSHDEER_TOKEN_HERE", // Update this
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		fmt.Println("✅ Admin user created successfully")
		fmt.Println("   Username: admin")
		fmt.Println("   Password: admin123")
		fmt.Println("   Please change the password after login!")

	case "invite":
		// Generate 10 invite codes
		for i := 0; i < 10; i++ {
			code, err := auth.GenerateInviteCode()
			if err != nil {
				log.Printf("Failed to generate invite code %d: %v", i+1, err)
				continue
			}
			inviteCode := models.InviteCode{
				Code: code,
				Time: time.Now(),
			}

			if err := db.Create(&inviteCode).Error; err != nil {
				log.Printf("Failed to create invite code %d: %v", i+1, err)
				continue
			}

			fmt.Printf("📨 Invite code %d: %s\n", i+1, code)
		}
		fmt.Println("✅ Invite codes generated successfully")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
