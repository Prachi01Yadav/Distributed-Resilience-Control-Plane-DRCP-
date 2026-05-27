package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/arche/sentinelmesh/internal/xds"
	"github.com/arche/sentinelmesh/pkg/cache"
	"github.com/arche/sentinelmesh/pkg/logger"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logg, _ := logger.InitLogger("development")
	defer logg.Sync()

	ctx := context.Background()

	// 1. Setup Envoy Cache & Server
	snapshotCache := envoycache.NewSnapshotCache(false, envoycache.IDHash{}, nil)
	srv := serverv3.NewServer(ctx, snapshotCache, nil)

	// 2. Setup Redis to listen for Breach events
	rdb := cache.NewRedisClient("localhost:6379", "", 0)
	if err := rdb.Ping(ctx); err != nil {
		logg.Warn("Redis not connected", zap.Error(err))
	}

	builder := xds.NewSnapshotBuilder(snapshotCache, logg)

	// NodeID should match what Envoy proxy sets in its bootstrap config (e.g., "envoy-node-1")
	err := builder.UpdateServiceBreaker(ctx, "envoy-node-1", "backend_service", false)
	if err != nil {
		logg.Error("Failed to build initial xDS snapshot", zap.Error(err))
	}

	// 3. Listen to Redis Pub/Sub for Breach Events
	pubsub := rdb.Client.Subscribe(ctx, "sla_breach_events")
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			serviceID := msg.Payload // The breached service ID
			logg.Warn("Received SLA Breach Event from Pub/Sub! Triggering Circuit Breaker...", zap.String("service", serviceID))
			
			// Trip the breaker
			err := builder.UpdateServiceBreaker(ctx, "envoy-node-1", serviceID, true)
			if err != nil {
				logg.Error("Failed to update circuit breaker via xDS", zap.Error(err))
			}
		}
	}()

	// 4. Start gRPC Server for xDS
	grpcServer := grpc.NewServer()
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, srv)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, srv)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 18000))
	if err != nil {
		logg.Fatal("Failed to listen on port 18000", zap.Error(err))
	}

	go func() {
		logg.Info("xDS Control Plane listening on :18000")
		if err := grpcServer.Serve(lis); err != nil {
			logg.Fatal("gRPC serve error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pubsub.Close()
	grpcServer.GracefulStop()
	logg.Info("xDS Server shut down gracefully")
}
