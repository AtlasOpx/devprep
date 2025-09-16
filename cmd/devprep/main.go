package main

import (
	"context"
	"fmt"
	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/routes"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	defer stopOngoingGracefully()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "DevPrep",
		AppName:       "Dev Prep app v1.0.1",
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   120 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Use(func(c *fiber.Ctx) error {
		if isShuttingDown.Load() {
			return c.Status(503).SendString("Service Unavailable")
		}
		return c.Next()
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		if isShuttingDown.Load() {
			return c.Status(503).JSON(fiber.Map{
				"status": "shutting_down",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		if isShuttingDown.Load() {
			return c.Status(503).JSON(fiber.Map{
				"ready":  false,
				"reason": "shutting_down",
			})
		}
		return c.JSON(fiber.Map{
			"ready": true,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		select {
		case <-time.After(100 * time.Millisecond):
			return c.SendString("Hello, World!")
		case <-ongoingCtx.Done():
			return c.Status(503).SendString("Request cancelled due to shutdown")
		}
	})

	routes.SetupRoutes(app, db, cfg)

	go func() {
		log.Println("Server starting on :3000")
		if err := app.Listen(fmt.Sprintf(":%v", cfg.ServerPort)); err != nil {
			log.Printf("Server failed to start: %v", err)
			stop()
		}
	}()

	<-rootCtx.Done()
	log.Println("Received shutdown signal, initiating graceful shutdown...")

	isShuttingDown.Store(true)

	log.Printf("Waiting %v for readiness checks to propagate...", _readinessDrainDelay)
	time.Sleep(_readinessDrainDelay)

	log.Println("Stopping acceptance of new requests and waiting for ongoing requests to finish...")
	stopOngoingGracefully()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown gracefully within %v: %v", _shutdownPeriod, err)

		log.Printf("Forcing shutdown in %v...", _shutdownHardPeriod)
		time.Sleep(_shutdownHardPeriod)

		os.Exit(1)
	}

	log.Println("Server shut down gracefully")
}
