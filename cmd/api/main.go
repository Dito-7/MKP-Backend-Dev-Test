package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mkp-cinema-ticketing/internal/config"
	delivery "mkp-cinema-ticketing/internal/delivery/http"
	"mkp-cinema-ticketing/internal/delivery/http/handler"
	"mkp-cinema-ticketing/internal/repository/postgres"
	"mkp-cinema-ticketing/internal/usecase"
	"mkp-cinema-ticketing/pkg/database"
)

func main() {
	log.Println("Starting Cinema Ticketing API Server...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Printf("[WARNING] PostgreSQL connection failed (%v). Continuing to start server for health check.", err)
	} else {
		defer db.Close()
	}

	userRepo := postgres.NewUserRepositoryPG(db)
	movieRepo := postgres.NewMovieRepositoryPG(db)
	studioRepo := postgres.NewStudioRepositoryPG(db)
	scheduleRepo := postgres.NewScheduleRepositoryPG(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
	scheduleUsecase := usecase.NewScheduleUsecase(scheduleRepo, movieRepo, studioRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	scheduleHandler := handler.NewScheduleHandler(scheduleUsecase)

	router := delivery.SetupRouter(cfg, authHandler, scheduleHandler)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on port :%s (Environment: %s)", cfg.AppPort, cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly.")
}
