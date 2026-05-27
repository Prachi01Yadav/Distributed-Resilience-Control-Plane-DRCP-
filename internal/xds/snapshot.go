package xds

import (
	"context"
	"fmt"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type SnapshotBuilder struct {
	cache   cache.SnapshotCache
	logger  *zap.Logger
	version int
}

func NewSnapshotBuilder(c cache.SnapshotCache, logger *zap.Logger) *SnapshotBuilder {
	return &SnapshotBuilder{
		cache:   c,
		logger:  logger,
		version: 0,
	}
}

// UpdateServiceBreaker generates a new xDS Snapshot that applies a strict CircuitBreaker or OutlierDetection.
func (s *SnapshotBuilder) UpdateServiceBreaker(ctx context.Context, nodeID, serviceName string, applyBreaker bool) error {
	s.version++
	versionStr := fmt.Sprintf("%d", s.version)

	maxRequests := uint32(1000)
	if applyBreaker {
		maxRequests = 1 // Trip the breaker
	}

	c := &cluster.Cluster{
		Name:                 serviceName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       &endpoint.ClusterLoadAssignment{ClusterName: serviceName},
		CircuitBreakers: &cluster.CircuitBreakers{
			Thresholds: []*cluster.CircuitBreakers_Thresholds{{
				MaxConnections:     &wrapperspb.UInt32Value{Value: maxRequests},
				MaxPendingRequests: &wrapperspb.UInt32Value{Value: maxRequests},
				MaxRequests:        &wrapperspb.UInt32Value{Value: maxRequests},
				MaxRetries:         &wrapperspb.UInt32Value{Value: 3},
			}},
		},
	}

	rt := &route.RouteConfiguration{
		Name: "local_route",
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{Prefix: "/"},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{Cluster: serviceName},
					},
				},
			}},
		}},
	}

	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http",
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: rt,
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: "envoy.filters.http.router",
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		return err
	}

	l := &listener.Listener{
		Name: "listener_0",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 10000,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: "envoy.filters.network.http_connection_manager",
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}

	snapshot, err := cache.NewSnapshot(versionStr,
		map[resource.Type][]types.Resource{
			resource.EndpointType: {},
			resource.ClusterType:  {c},
			resource.RouteType:    {rt},
			resource.ListenerType: {l},
			resource.RuntimeType:  {},
			resource.SecretType:   {},
		},
	)
	if err != nil {
		return err
	}

	if err := snapshot.Consistent(); err != nil {
		return fmt.Errorf("snapshot inconsistency: %w", err)
	}
	s.logger.Info("Pushing xDS snapshot to Envoy node", zap.String("node", nodeID), zap.Bool("breaker_active", applyBreaker), zap.String("version", versionStr))
	return s.cache.SetSnapshot(ctx, nodeID, snapshot)
}
