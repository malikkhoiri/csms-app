package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/malikkhoiri/csms/internal/application/service"
	"github.com/malikkhoiri/csms/internal/config"
	"github.com/malikkhoiri/csms/internal/domain"
	"github.com/malikkhoiri/csms/internal/infrastructure/database"
	"github.com/malikkhoiri/csms/internal/infrastructure/repository"
)

func main() {
	// Command line flags
	var (
		forceReset = flag.Bool("force", false, "Force reset existing data (delete and recreate)")
		skipUser   = flag.Bool("skip-user", false, "Skip seeding admin user")
		skipCP     = flag.Bool("skip-cp", false, "Skip seeding charge points")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	postgresDB, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgresDB.Close()

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(postgresDB.DB)
	userService := service.NewUserService(userRepo)

	chargePointRepo := repository.NewChargePointRepository(postgresDB.DB)
	connectorRepo := repository.NewConnectorRepository(postgresDB.DB)
	chargePointService := service.NewChargePointService(chargePointRepo, connectorRepo)

	ctx := context.Background()

	log.Println("🚀 Starting CSMS database seeding...")
	if *forceReset {
		log.Println("⚠️  Force reset mode enabled - existing data will be deleted!")
	}

	// Seed admin user
	if !*skipUser {
		seedAdminUser(ctx, userService, *forceReset)
	} else {
		log.Println("⏭️  Skipping admin user seeding...")
	}

	// Seed charge points
	if !*skipCP {
		seedChargePoints(ctx, chargePointService, *forceReset)
	} else {
		log.Println("⏭️  Skipping charge points seeding...")
	}

	log.Println("✨ Seeding completed!")
}

func seedAdminUser(ctx context.Context, userService domain.UserService, forceReset bool) {
	// Check if admin user already exists
	existingUser, err := userService.GetUserByEmail(ctx, "admin@csms.com")
	if err == nil && existingUser != nil {
		if forceReset {
			err = userService.DeleteUser(ctx, existingUser.ID)
			if err != nil {
				log.Printf("❌ Failed to delete existing admin user: %v", err)
				return
			}
			log.Println("⚠️  Existing admin user deleted!")
		} else {
			log.Println("ℹ️  Admin user already exists, skipping...")
			log.Printf("   Email: %s", existingUser.Email)
			return
		}
	}

	// Create default admin user
	adminUser := &domain.User{
		Name:     "Admin",
		Email:    "admin@csms.com",
		Password: "admin123",
		Phone:    "+1234567890",
		Role:     "admin",
		Status:   "active",
	}

	err = userService.CreateUser(ctx, adminUser)
	if err != nil {
		log.Printf("❌ Failed to create admin user: %v", err)
	} else {
		log.Println("✅ Admin user created successfully!")
		log.Printf("   Email: %s", adminUser.Email)
		log.Printf("   Password: admin123")
	}
}

func seedChargePoints(ctx context.Context, chargePointService domain.ChargePointService, forceReset bool) {
	// Sample charge points data
	chargePoints := []*domain.ChargePoint{
		{
			ChargePointCode:         "CP001",
			ChargePointModel:        "Tesla Supercharger V3",
			ChargePointVendor:       "Tesla",
			ChargePointSerialNumber: "TS-V3-001",
			FirmwareVersion:         "1.2.3",
			Iccid:                   "89014103211118510720",
			Imsi:                    "310260123456789",
			MeterType:               "AC",
			MeterSerialNumber:       "MTR-001",
			Status:                  "Available",
			LastHeartbeat:           time.Now(),
			LastBootNotification:    time.Now(),
		},
		{
			ChargePointCode:         "CP002",
			ChargePointModel:        "ABB Terra 184",
			ChargePointVendor:       "ABB",
			ChargePointSerialNumber: "ABB-184-002",
			FirmwareVersion:         "2.1.0",
			Iccid:                   "89014103211118510721",
			Imsi:                    "310260123456790",
			MeterType:               "DC",
			MeterSerialNumber:       "MTR-002",
			Status:                  "Available",
			LastHeartbeat:           time.Now(),
			LastBootNotification:    time.Now(),
		},
		{
			ChargePointCode:         "CP003",
			ChargePointModel:        "Siemens VersiCharge",
			ChargePointVendor:       "Siemens",
			ChargePointSerialNumber: "SI-VC-003",
			FirmwareVersion:         "1.8.5",
			Iccid:                   "89014103211118510722",
			Imsi:                    "310260123456791",
			MeterType:               "AC",
			MeterSerialNumber:       "MTR-003",
			Status:                  "Available",
			LastHeartbeat:           time.Now(),
			LastBootNotification:    time.Now(),
		},
		{
			ChargePointCode:         "CP004",
			ChargePointModel:        "ChargePoint CT4000",
			ChargePointVendor:       "ChargePoint",
			ChargePointSerialNumber: "CP-CT4-004",
			FirmwareVersion:         "3.2.1",
			Iccid:                   "89014103211118510723",
			Imsi:                    "310260123456792",
			MeterType:               "DC",
			MeterSerialNumber:       "MTR-004",
			Status:                  "Available",
			LastHeartbeat:           time.Now(),
			LastBootNotification:    time.Now(),
		},
		{
			ChargePointCode:         "CP005",
			ChargePointModel:        "EVBox BusinessLine",
			ChargePointVendor:       "EVBox",
			ChargePointSerialNumber: "EV-BL-005",
			FirmwareVersion:         "2.0.4",
			Iccid:                   "89014103211118510724",
			Imsi:                    "310260123456793",
			MeterType:               "AC",
			MeterSerialNumber:       "MTR-005",
			Status:                  "Available",
			LastHeartbeat:           time.Now(),
			LastBootNotification:    time.Now(),
		},
	}

	log.Println("🌱 Seeding charge points...")

	createdCount := 0
	skippedCount := 0

	for _, cp := range chargePoints {
		// Check if charge point already exists
		existingCP, err := chargePointService.GetChargePointByCode(ctx, cp.ChargePointCode)
		if err == nil && existingCP != nil {
			if forceReset {
				err = chargePointService.DeleteChargePoint(ctx, existingCP.ID)
				if err != nil {
					log.Printf("❌ Failed to delete existing charge point %s: %v", cp.ChargePointCode, err)
					continue
				}
				log.Printf("⚠️  Existing charge point %s deleted!", cp.ChargePointCode)
			} else {
				log.Printf("ℹ️  Charge point %s already exists, skipping...", cp.ChargePointCode)
				skippedCount++
				continue
			}
		}

		// Create charge point
		_, err = chargePointService.RegisterChargePoint(ctx, &domain.BootNotificationRequest{
			ChargePointVendor:       cp.ChargePointVendor,
			ChargePointModel:        cp.ChargePointModel,
			ChargePointSerialNumber: cp.ChargePointSerialNumber,
			ChargeBoxSerialNumber:   cp.ChargeBoxSerialNumber,
			FirmwareVersion:         cp.FirmwareVersion,
			Iccid:                   cp.Iccid,
			Imsi:                    cp.Imsi,
			MeterType:               cp.MeterType,
			MeterSerialNumber:       cp.MeterSerialNumber,
		}, cp.ChargePointCode)

		if err != nil {
			log.Printf("❌ Failed to create charge point %s: %v", cp.ChargePointCode, err)
		} else {
			log.Printf("✅ Charge point %s created successfully!", cp.ChargePointCode)
			createdCount++
		}
	}

	log.Printf("🎉 Charge points seeding completed! Created: %d, Skipped: %d", createdCount, skippedCount)
}
