package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/ycd/leaderboard/internal/repository"
	"github.com/ycd/leaderboard/internal/service/admin"
	"github.com/ycd/leaderboard/internal/service/api/handlers"
	"github.com/ycd/leaderboard/internal/service/event"
	"github.com/ycd/leaderboard/internal/service/user"
	_ "modernc.org/sqlite"

	"github.com/ycd/leaderboard/internal/service/api"
	"github.com/ycd/leaderboard/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables.")
	}
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initialize ClickHouse connection

	chDB, err := storage.NewClickHouseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer chDB.Close()
	log.Println("Connected to ClickHouse")

	// Initialize SQLite connection
	sqliteDB, err := storage.NewSQLiteConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDB.Close()

	// Initialize Redis connection
	redisDB, err := storage.NewRedisConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer redisDB.Close()
	log.Println("Connected to Redis")

	// Create AdminRepository
	adminRepo := repository.NewAdminRepository(chDB, sqliteDB)

	// Use adminRepo in your AdminService
	userRepo := repository.NewUserRepository(sqliteDB)
	eventRepo := repository.NewEventRepository(chDB, sqliteDB, redisDB)

	userService := user.NewUserService(userRepo)
	eventService := event.NewEventService(eventRepo, userRepo)
	adminService := admin.NewAdminService(userRepo, eventRepo, adminRepo)

	adminHandler := handlers.NewAdminHandler(adminService)
	eventHandler := handlers.NewEventHandler(eventService)
	userHandler := handlers.NewUserHandler(userService)

	server := api.NewServer(userHandler, eventHandler, adminHandler)

	errCh := make(chan error)

	go func() {
		if err := server.Serve(); err != nil {
			errCh <- err
		}
		close(errCh)
	}()
	log.Printf("Starting server at %s\n", server.Addr)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case err, ok := <-errCh:
		if ok {
			log.Println("server error:", err)
		}

	case sig := <-sigCh:
		log.Printf("Signal %s received\n", sig)
		if err := server.Shutdown(); err != nil {
			log.Println("Failed to shutdown server:", err)
		}
		log.Println("server shutdown")
	}
}
