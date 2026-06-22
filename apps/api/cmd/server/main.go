package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"continuum/api/internal/handler"
	"continuum/api/internal/repository"
	"continuum/api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the environment configuration BEFORE evaluating variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("ℹ️ No root .env file found, relying purely on host environment parameters.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/continuum?sslmode=disable"
	}

	log.Println("🔌 Connecting to Continuum Relational Storage Layer...")
	db, err := repository.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("❌ Core Database connection crash: %v", err)
	}

	// 2. Initialize Layered Domain Entities & Real Business Services
	lndConfig := &service.LNDConfig{
		Host:         os.Getenv("LND_HOST"),
		MacaroonPath: os.Getenv("LND_MACAROON_PATH"),
		TLSCertPath:  os.Getenv("LND_TLS_CERT_PATH"),
	}
	
	lightningService, err := service.NewLightningService(lndConfig)
	if err != nil {
		log.Printf("⚠️ LND connection initialization failed: %v. Node falling back onto simulation mode.", err)
	}

	// Connect the dependencies safely across services
	vaultService := service.NewVaultService(db, lightningService)
	schedulerService := service.NewScheduler(db)

	// 3. Launch Non-Blocking Background Scheduler Daemon
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Evaluate the expiration matrix every 10 seconds for clean production polling efficiency
	schedulerService.StartCheckLoop(ctx, 10*time.Second)
	log.Println("⏰ Autonomous dead-man switch tracking ticker online.")

	// 4. Provision HTTP Transport Adapters (Fiber Framework Instance)
	app := fiber.New(fiber.Config{
		AppName: "Continuum Anonymized Inheritance Core v1.0",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, 
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	// Inject the proper architectural layers into handlers
	vaultHandler := handler.NewVaultHandler(db, vaultService)
	proofHandler := handler.NewProofHandler(db, lightningService)
	recoveryHandler := handler.NewRecoveryHandler(db)

	// 5. Structure API Routing Table Tree Layouts
	api := app.Group("/api")
	{
		api.Post("/vaults", vaultHandler.CreateVault)
		api.Get("/vaults/:id", recoveryHandler.GetVaultStatus)
		
		// ⚡ Lightning Network Interactivity Core Flow Endpoints
		api.Post("/vaults/:id/invoice", vaultHandler.RequestCheckInToken)
		api.Post("/vaults/:id/checkin", vaultHandler.ConfirmCheckIn)
		
		// 🔒 Production Security Guard: Protect the time-warp testing backdoor
		if os.Getenv("ALLOW_DEV_TIME_WARP") == "true" {
			log.Println("⚠️ WARNING: Time Warp simulation backdoor route activated.")
			api.Post("/vaults/:id/warp", proofHandler.SimulateTimeWarp)
		} else {
			log.Println("🔒 Production lockdown active: Time warp endpoints are explicitly disabled.")
		}
	}

	// 6. Establish Graceful Shutdown Routines
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Printf("🚀 Continuum Protocol Core HTTP node listening on port :%s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("⚠️ Server shut down warning: %v", err)
		}
	}()

	// Listen for system kill interrupts
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Continuum Node gracefully...")
	cancel() // Shuts down background tracking checks immediately
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Forced cluster closure error: %v", err)
	}
	log.Println("✨ Server instance closed safely.")
}