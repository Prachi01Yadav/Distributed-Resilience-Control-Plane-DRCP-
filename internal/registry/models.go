package registry

import (
	"time"
)

// Service represents a microservice registered in the mesh
type Service struct {
	ID          string    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	SLAContracts []SLAContract `gorm:"foreignKey:ServiceID" json:"sla_contracts,omitempty"`
}

// SLAContract defines the latency and error rate budgets for a service
type SLAContract struct {
	ID        string    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	ServiceID string    `gorm:"index;not null" json:"service_id"`
	Policy    string    `gorm:"type:text;not null" json:"policy"` // YAML or JSON representing the OPA policy
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Incident logs an SLA breach event
type Incident struct {
	ID               string     `gorm:"primaryKey;type:varchar(100)" json:"id"`
	ServiceID        string     `gorm:"index;not null" json:"service_id"`
	ContractID       string     `gorm:"index;not null" json:"contract_id"`
	ErrorRate        float64    `json:"error_rate"`
	P99Latency       float64    `json:"p99_latency"`
	Status           string     `gorm:"type:varchar(50);default:'OPEN'" json:"status"` // OPEN, RESOLVED, ANCHORED
	BlockchainTxHash string     `gorm:"type:varchar(100)" json:"blockchain_tx_hash,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
}
