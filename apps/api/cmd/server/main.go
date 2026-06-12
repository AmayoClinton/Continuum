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
	// 1. Load .env (best-effort) and extract environment variables
	_ = godotenv.Load()

	// Extract Environment variables or fallback to local Docker infrastructure defaults
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/continuum?sslmode=disable"
	}

	log.Println("🔌 Connecting to Continuum Relational Storage Layer...")
	db, err := repository.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("❌ Core Database connection crash: %v", err)
	}

	// 2. Initialize Layered Domain Entities & Business Services
	vaultService := service.NewVaultService(db)
	schedulerService := service.NewScheduler(db)
	lightningService := service.NewLightningService(nil)
	multisigService := service.NewMultisigServiceFromEnv()

	// 3. Launch Non-Blocking Background Scheduler Daemon
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	schedulerService.StartCheckLoop(ctx, 5*time.Second)
	log.Println("⏰ Autonomous dead-man switch tracking ticker online.")

	// 4. Provision HTTP Transport Adapters (Fiber Framework Instance)
	app := fiber.New(fiber.Config{
		AppName: "Continuum Anonymized Inheritance Core v1.0",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST", "PATCH", "OPTIONS"},
	}))

	vaultHandler := handler.NewVaultHandler(db, vaultService, lightningService, multisigService)
	proofHandler := handler.NewProofHandler(db)
	recoveryHandler := handler.NewRecoveryHandler(db)

	// 5. Structure API Routing Table Tree Layouts
	api := app.Group("/api")
	{
		api.Get("/health", func(c fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "ok", "service": "continuum-api"})
		})
		api.Get("/vaults", vaultHandler.ListVaults)
		api.Post("/vaults", vaultHandler.CreateVault)
		api.Get("/vaults/:id", recoveryHandler.GetVaultStatus)
		api.Post("/vaults/:id/check-in", vaultHandler.CheckIn)
		api.Post("/vaults/:id/invoice", vaultHandler.CreateProofInvoice)
		api.Patch("/vaults/:id/timer", vaultHandler.UpdateTimer)
		api.Post("/vaults/:id/beneficiaries", vaultHandler.AddBeneficiary)
		api.Post("/vaults/:id/warp", proofHandler.SimulateTimeWarp)
	}

	// 6. Establish Graceful Shutdown Routines
	go func() {
		log.Println("🚀 Continuum Protocol Core HTTP node listening on port :8080")
		if err := app.Listen(":8080"); err != nil {
			log.Printf("⚠️ Server shut down warning: %v", err)
		}
	}()

	// Listen for system kill interrupts
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Continuum Node gracefully...")
	cancel()
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Forced cluster closure error: %v", err)
	}
	log.Println("✨ Server instance closed safely.")
}
