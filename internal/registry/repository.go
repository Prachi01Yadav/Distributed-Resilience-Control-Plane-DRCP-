package registry

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the interface for registry data access
type Repository interface {
	CreateService(ctx context.Context, service *Service) error
	GetService(ctx context.Context, id string) (*Service, error)
	ListServices(ctx context.Context) ([]Service, error)
	
	CreateSLAContract(ctx context.Context, contract *SLAContract) error
	GetSLAContract(ctx context.Context, id string) (*SLAContract, error)
	GetContractsByService(ctx context.Context, serviceID string) ([]SLAContract, error)
	ListSLAContracts(ctx context.Context) ([]SLAContract, error)

	RecordIncident(ctx context.Context, incident *Incident) error
	ListIncidents(ctx context.Context) ([]Incident, error)
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new postgres repository
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateService(ctx context.Context, service *Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *postgresRepository) GetService(ctx context.Context, id string) (*Service, error) {
	var service Service
	err := r.db.WithContext(ctx).Preload("SLAContracts").First(&service, "id = ?", id).Error
	return &service, err
}

func (r *postgresRepository) ListServices(ctx context.Context) ([]Service, error) {
	var services []Service
	err := r.db.WithContext(ctx).Find(&services).Error
	return services, err
}

func (r *postgresRepository) CreateSLAContract(ctx context.Context, contract *SLAContract) error {
	return r.db.WithContext(ctx).Create(contract).Error
}

func (r *postgresRepository) GetSLAContract(ctx context.Context, id string) (*SLAContract, error) {
	var contract SLAContract
	err := r.db.WithContext(ctx).First(&contract, "id = ?", id).Error
	return &contract, err
}

func (r *postgresRepository) GetContractsByService(ctx context.Context, serviceID string) ([]SLAContract, error) {
	var contracts []SLAContract
	err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).Find(&contracts).Error
	return contracts, err
}

func (r *postgresRepository) ListSLAContracts(ctx context.Context) ([]SLAContract, error) {
	var contracts []SLAContract
	err := r.db.WithContext(ctx).Find(&contracts).Error
	return contracts, err
}

func (r *postgresRepository) RecordIncident(ctx context.Context, incident *Incident) error {
	return r.db.WithContext(ctx).Create(incident).Error
}

func (r *postgresRepository) ListIncidents(ctx context.Context) ([]Incident, error) {
	var incidents []Incident
	err := r.db.WithContext(ctx).Find(&incidents).Error
	return incidents, err
}
