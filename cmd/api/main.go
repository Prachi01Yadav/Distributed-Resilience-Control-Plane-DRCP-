package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arche/sentinelmesh/internal/registry"
	"github.com/arche/sentinelmesh/internal/telemetry"
	"github.com/arche/sentinelmesh/pkg/db"
	"github.com/arche/sentinelmesh/pkg/kafka"
	"github.com/arche/sentinelmesh/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Initialize Logger
	logg, err := logger.InitLogger("development")
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logg.Sync()

	// Initialize DB
	var database *gorm.DB
	
	if os.Getenv("DATABASE_TYPE") == "sqlite" {
		logg.Info("Using SQLite database")
		database, err = db.NewSqliteDB("sentinelmesh.db")
		if err != nil {
			logg.Fatal("Failed to connect to SQLite", zap.Error(err))
		}
	} else {
		logg.Info("Using PostgreSQL database")
		dbCfg := db.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "sentinelmesh",
			SSLMode:  "disable",
		}
		database, err = db.NewPostgresDB(dbCfg)
		if err != nil {
			logg.Fatal("Failed to connect to database", zap.Error(err))
		}
	}

	// Auto-Migrate schema
	logg.Info("Running database migrations...")
	err = database.AutoMigrate(
		&registry.Service{},
		&registry.SLAContract{},
		&registry.Incident{},
	)
	if err != nil {
		logg.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize repository and handler
	repo := registry.NewPostgresRepository(database)
	handler := registry.NewHandler(repo, logg)

	// Setup Kafka Producer
	producer := kafka.NewProducer([]string{"localhost:9092"}, "telemetry-events", logg)
	defer producer.Close()

	telemetryHandler := telemetry.NewHandler(producer, logg)

	// Setup Gin Router
	router := gin.Default()
	
	// Serve Static Files for Dashboard
	router.Static("/static", "./web")
	router.StaticFile("/", "./web/index.html")

	// Register Routes
	handler.RegisterRoutes(router)
	router.POST("/api/v1/telemetry", telemetryHandler.ReceiveTelemetry)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Server Configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful Shutdown
	go func() {
		logg.Info("Starting API server on :" + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("listen: ", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logg.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logg.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logg.Info("Server exiting")
}
