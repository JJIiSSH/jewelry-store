package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JJIiSSH/jewelry-store/internal/config"
	"github.com/JJIiSSH/jewelry-store/internal/handler/httphandler"

	"github.com/JJIiSSH/jewelry-store/internal/repository/postgres"
	"github.com/JJIiSSH/jewelry-store/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.DB.DSN())

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	productRepo := postgres.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := httphandler.NewProductHandler(productService)

	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userHandler := httphandler.NewAuthHandler(authService)

	router := gin.Default()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	v1 := router.Group("/api/v1")

	productHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Printf("server error: %v", err)
	case sig := <-quit:
		log.Printf("received signal %v, shutting down...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
