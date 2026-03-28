package services

import (
	"context"
	"efmob/internal/dto"
	dbinterface "efmob/internal/repositories/db_interface"
	"fmt"
)

type ServiceRepo struct {
	db dbinterface.DbIface
}

func NewServiceRepo(pool dbinterface.DbIface) *ServiceRepo {
	return &ServiceRepo{
		db: pool,
	}
}

func (r *ServiceRepo) CreateService(ctx context.Context, serviceInfo *dto.ServiceInfo) error {
	row := r.db.QueryRow(ctx, insertServiceQuery(), serviceInfo.Name)
	var (
		id int
	)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("CreateService - failed to insert service: %w", err)
	}
	serviceInfo.Id = &id

	return nil
}

func (r *ServiceRepo) GetService(ctx context.Context, serviceInfo *dto.ServiceInfo) error {
	row := r.db.QueryRow(ctx, getServiceQuery(), serviceInfo.Name)
	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("GetService - failed to get service by name: %w", err)
	}
	serviceInfo.Id = &id

	return nil
}

func (r *ServiceRepo) GetOrCreate(ctx context.Context, serviceInfo *dto.ServiceInfo) error {
	row := r.db.QueryRow(ctx, getOrCreateServiceQuery(), serviceInfo.Name)
	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("GetOrCreate - failed to get or create service: %w", err)
	}
	serviceInfo.Id = &id

	return nil
}
