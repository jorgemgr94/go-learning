package main

import (
	"context"
	"fmt"
	"go-learning/internal/config"
	"go-learning/internal/routers"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	config := config.LoadConfig()

	fmt.Println("Config:", config)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":    Version,
			"build_time": BuildTime,
			"git_commit": GitCommit,
		})
	})
	router.GET("/graceful-test", func(c *gin.Context) {
		slog.Info("slow request started")
		time.Sleep(20 * time.Second) // Simulate slow processing
		slog.Info("graceful request completed")
		c.JSON(200, gin.H{"message": "Graceful response completed"})
	})

	apiV1 := router.Group("/api/v1")
	{
		routers.UserRouter(apiV1)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	// Start server in a goroutine so it doesn't block
	go func() {
		slog.Info("server started", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	// Register the channel to receive specific signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	<-quit
	slog.Info("shutting down server...")

	// Create a context with timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("server exited gracefully")
}
