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

	"github.com/JJIiSSH/jewelry-store/internal/handler/httphandler"

	"github.com/JJIiSSH/jewelry-store/internal/repository/postgres"
	"github.com/JJIiSSH/jewelry-store/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	dsn := "postgres://IvanDev:1111@127.0.0.1:5432/mydb?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	productRepo := postgres.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := httphandler.NewProductHandler(productService)

	router := gin.Default()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	v1 := router.Group("/api/v1")

	productHandler.RegisterRoutes(v1)

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
