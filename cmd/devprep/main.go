package main

import (
	"context"
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

	app.Use(func(c *fiber.Ctx) error {
		if isShuttingDown.Load() {
			return c.Status(503).SendString("Service Unavailable")
		}
		return c.Next()
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		if isShuttingDown.Load() {
			return c.Status(503).JSON(fiber.Map{
				"status": "shutting_down",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
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

	go func() {
		log.Println("Server starting on :3000")
		if err := app.Listen(":3000"); err != nil {
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
