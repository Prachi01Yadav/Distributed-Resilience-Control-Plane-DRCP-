package main

import (
	"context"
	"crypto/sha256"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arche/sentinelmesh/internal/anchor"
	"github.com/arche/sentinelmesh/internal/registry"
	"github.com/arche/sentinelmesh/pkg/cache"
	"github.com/arche/sentinelmesh/pkg/db"
	"github.com/arche/sentinelmesh/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	logg, _ := logger.InitLogger("development")
	defer logg.Sync()

	ctx := context.Background()

	// 1. Setup Postgres
	dbCfg := db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "sentinelmesh",
		SSLMode:  "disable",
	}
	database, err := db.NewPostgresDB(dbCfg)
	if err != nil {
		logg.Fatal("Failed to connect to DB", zap.Error(err))
	}
	repo := registry.NewPostgresRepository(database)

	// 2. Setup Ethereum Anchor
	// Using a local RPC (Ganache/Hardhat) and a dummy key for testing
	ethRpc := "http://127.0.0.1:8545"
	dummyKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19" 
	contractAddr := "0x0000000000000000000000000000000000000000"

	ethAnchor, err := anchor.NewEthereumAnchor(ethRpc, dummyKey, contractAddr, logg)
	if err != nil {
		logg.Warn("Could not connect to Ethereum node. Will mock transaction hashes.", zap.Error(err))
	}

	// 3. Listen to Redis Pub/Sub for Breach Events
	rdb := cache.NewRedisClient("localhost:6379", "", 0)
	pubsub := rdb.Client.Subscribe(ctx, "sla_breach_events")
	
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			serviceID := msg.Payload
			logg.Info("Received SLA Breach Event! Anchoring to Blockchain...", zap.String("service", serviceID))

			// Generate a simple hash (acting as merkle root of incident data)
			hash := sha256.Sum256([]byte(serviceID + time.Now().String()))

			// Submit to Blockchain
			var txHash string
			if ethAnchor != nil {
				tx, err := ethAnchor.RecordBreach(ctx, serviceID, hash)
				if err != nil {
					logg.Error("Failed to anchor breach", zap.Error(err))
					txHash = "failed"
				} else {
					txHash = tx
				}
			} else {
				// Mock mode if node is down
				txHash = "0xmock_tx_" + uuid.New().String()
			}

			// Save Incident to Database
			incident := &registry.Incident{
				ID:               uuid.New().String(),
				ServiceID:        serviceID,
				ContractID:       "unknown-for-now", // Fetched from cache in a real app
				ErrorRate:        1.0,               // Mock data
				P99Latency:       500.0,             // Mock data
				Status:           "ANCHORED",
				BlockchainTxHash: txHash,
				CreatedAt:        time.Now(),
			}

			if err := repo.RecordIncident(ctx, incident); err != nil {
				logg.Error("Failed to record incident to Postgres", zap.Error(err))
			} else {
				logg.Info("Incident anchored to Ethereum and saved to DB", zap.String("txHash", txHash))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pubsub.Close()
	logg.Info("Anchor Service shut down gracefully")
}
