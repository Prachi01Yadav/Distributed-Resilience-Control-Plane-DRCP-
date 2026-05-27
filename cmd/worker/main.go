package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/arche/sentinelmesh/internal/budget"
	"github.com/arche/sentinelmesh/internal/policy"
	"github.com/arche/sentinelmesh/internal/telemetry"
	"github.com/arche/sentinelmesh/pkg/cache"
	"github.com/arche/sentinelmesh/pkg/kafka"
	"github.com/arche/sentinelmesh/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logg, _ := logger.InitLogger("development")
	defer logg.Sync()

	rdb := cache.NewRedisClient("localhost:6379", "", 0)
	err := rdb.Ping(context.Background())
	if err != nil {
		logg.Warn("Redis not connected. Ensure it's running", zap.Error(err))
	} else {
		logg.Info("Connected to Redis")
	}

	calculator := budget.NewWindowCalculator(rdb)
	opaEngine := policy.NewOPAEngine()

	consumer := kafka.NewConsumer([]string{"localhost:9092"}, "telemetry-events", "budget-worker", logg)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logg.Info("Starting budget worker...")

	go consumer.ReadMessages(ctx, func(msg []byte) error {
		var data telemetry.TelemetryData
		if err := json.Unmarshal(msg, &data); err != nil {
			return err
		}

		isError := data.StatusCode >= 500
		err := calculator.AddMetric(context.Background(), data.ServiceID, isError, data.LatencyMs)
		if err != nil {
			logg.Error("Failed to add metric", zap.Error(err))
			return err
		}

		// Calculate current rate
		errorRate, _ := calculator.GetErrorRate(context.Background(), data.ServiceID)
		
		// In a real app, fetch policy from Registry/DB here.
		// For now, we mock the evaluation.
		mockPolicy := `
		package sla
		default breach = false
		breach { input.error_rate > 0.05 }
		`
		breach, err := opaEngine.EvaluateSLA(context.Background(), mockPolicy, errorRate, data.LatencyMs)
		if err != nil {
			logg.Error("Failed to evaluate SLA", zap.Error(err))
			return err
		}

		if breach {
			logg.Warn("SLA BREACH DETECTED!", zap.String("service", data.ServiceID), zap.Float64("error_rate", errorRate))
			// Push to xDS or trigger incident
			rdb.Client.Publish(context.Background(), "sla_breach_events", data.ServiceID)
		} else {
			logg.Debug("SLA Budget OK", zap.String("service", data.ServiceID), zap.Float64("error_rate", errorRate))
		}

		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logg.Info("Shutting down worker...")
}
