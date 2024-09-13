package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ycd/leaderboard/internal/service/api/handlers"
)

type Server struct {
	Addr        string
	MustClose   bool
	server      *fiber.App
	maintenance bool

	userHandler  *handlers.UserHandler
	eventHandler *handlers.EventHandler
	adminHandler *handlers.AdminHandler
}

func NewServer(userHandler *handlers.UserHandler, eventHandler *handlers.EventHandler, adminHandler *handlers.AdminHandler) *Server {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Println("failed to retrieve port from env")
		port = "8080"
	}

	server := &Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		userHandler:  userHandler,
		eventHandler: eventHandler,
		adminHandler: adminHandler,
	}

	return server
}

func (s *Server) Serve() error {
	app := fiber.New()

	s.server = app
	s.RegisterHandlers()

	if err := s.server.Listen(s.Addr); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf(err.Error())
		}
	}

	return nil
}

func (s *Server) RegisterHandlers() {
	userGroup := s.server.Group("/user")
	eventGroup := s.server.Group("/event")
	adminGroup := s.server.Group("/admin")

	// User Operations
	userGroup.Post("/profile", s.userHandler.SetProfile)
	userGroup.Post("/level", s.userHandler.SetLevel)
	userGroup.Post("/coin", s.userHandler.SetCoin)

	// Event Operations
	eventGroup.Post("/join", s.eventHandler.JoinEvent)
	eventGroup.Post("/leaderboard/progress", s.eventHandler.SetLeaderboardProgress)
	eventGroup.Get("/leaderboard", s.eventHandler.GetLeaderboard)
	eventGroup.Get("/country-leaderboard", s.eventHandler.GetCountryLeaderboard)
	eventGroup.Post("/leaderboard/claim", s.eventHandler.ClaimReward)

	// Admin Operations
	adminGroup.Post("/event", s.adminHandler.CreateEvent)
	adminGroup.Get("/events", s.adminHandler.ListEvents)
	adminGroup.Delete("/event", s.adminHandler.DeleteEvent)
	adminGroup.Get("/user", s.adminHandler.GetUserDetails)
	adminGroup.Get("/users", s.adminHandler.ListUsers)
	adminGroup.Post("/user/ban", s.adminHandler.BanUser)
	adminGroup.Post("/user/unban", s.adminHandler.UnbanUser)
	adminGroup.Post("/leaderboard/reset", s.adminHandler.ResetLeaderboard)
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.ShutdownWithContext(ctx)
}

func (s *Server) ShutdownWithContext(ctx context.Context) error {
	if err := s.server.Server().ShutdownWithContext(ctx); err != nil {
		if s.MustClose {
			s.server.Server()
		}
		return err
	}

	return nil
}
