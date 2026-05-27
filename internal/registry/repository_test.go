package registry

import (
	"context"
	"testing"
	"time"

	"github.com/arche/sentinelmesh/pkg/db"
)

// TestRepositoryIntegration uses an in-memory SQLite database to test the repository layer
func TestRepositoryIntegration(t *testing.T) {
	database, err := db.NewSqliteDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	// Migrate schemas for testing
	err = database.AutoMigrate(&Service{}, &SLAContract{}, &Incident{})
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	repo := NewPostgresRepository(database)
	ctx := context.Background()

	t.Run("Create and Retrieve Service", func(t *testing.T) {
		svc := &Service{
			ID:          "test-svc-1",
			Name:        "OrderService",
			Description: "Handles customer orders",
			Owner:       "CoreTeam",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Test Creation
		err := repo.CreateService(ctx, svc)
		if err != nil {
			t.Fatalf("expected no error creating service, got %v", err)
		}

		// Test Retrieval
		fetched, err := repo.GetService(ctx, "test-svc-1")
		if err != nil {
			t.Fatalf("expected no error fetching service, got %v", err)
		}

		if fetched.Name != "OrderService" {
			t.Errorf("expected name OrderService, got %s", fetched.Name)
		}
	})

	t.Run("Record and List Incidents", func(t *testing.T) {
		incident := &Incident{
			ID:               "inc-1",
			ServiceID:        "test-svc-1",
			ContractID:       "contract-1",
			ErrorRate:        0.08,
			P99Latency:       150.0,
			Status:           "OPEN",
			BlockchainTxHash: "0x123abc",
			CreatedAt:        time.Now(),
		}

		err := repo.RecordIncident(ctx, incident)
		if err != nil {
			t.Fatalf("expected no error recording incident, got %v", err)
		}

		incidents, err := repo.ListIncidents(ctx)
		if err != nil {
			t.Fatalf("expected no error listing incidents, got %v", err)
		}

		if len(incidents) != 1 {
			t.Errorf("expected 1 incident, got %d", len(incidents))
		}
	})
}
