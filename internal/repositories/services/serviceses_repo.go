package services

import (
	"context"
	dbinterface "efmob/internal/repositories/db_interface"
	"fmt"

	"efmob/logger"

	"go.uber.org/zap"
)

type ServiceRepo struct {
	db dbinterface.DbIface
}

func NewServiceRepo(pool dbinterface.DbIface) *ServiceRepo {
	return &ServiceRepo{
		db: pool,
	}
}

func (r *ServiceRepo) CreateService(ctx context.Context, serviceInfo *ServiceInfo) error {
	logger.Log().Info("ServiceRepo.CreateService call", zap.String("name", serviceInfo.Name))

	row := r.db.QueryRow(ctx, insertServiceQuery(), serviceInfo.Name)
	var (
		id int
	)
	if err := row.Scan(&id); err != nil {
		logger.Log().Error("ServiceRepo.CreateService failed", zap.Error(err), zap.String("name", serviceInfo.Name))
		return fmt.Errorf("CreateService - failed to insert service: %w", err)
	}
	serviceInfo.Id = &id

	logger.Log().Info("ServiceRepo.CreateService ok", zap.String("name", serviceInfo.Name), zap.Int("service_id", id))
	return nil
}

func (r *ServiceRepo) GetService(ctx context.Context, serviceInfo *ServiceInfo) error {
	logger.Log().Info("ServiceRepo.GetService call", zap.String("name", serviceInfo.Name))

	row := r.db.QueryRow(ctx, getServiceQuery(), serviceInfo.Name)
	var id int
	if err := row.Scan(&id); err != nil {
		logger.Log().Error("ServiceRepo.GetService failed", zap.Error(err), zap.String("name", serviceInfo.Name))
		return fmt.Errorf("GetService - failed to get service by name: %w", err)
	}
	serviceInfo.Id = &id

	logger.Log().Info("ServiceRepo.GetService ok", zap.String("name", serviceInfo.Name), zap.Int("service_id", id))
	return nil
}

func (r *ServiceRepo) GetOrCreate(ctx context.Context, serviceInfo *ServiceInfo) error {
	logger.Log().Info("ServiceRepo.GetOrCreate call", zap.String("name", serviceInfo.Name))

	row := r.db.QueryRow(ctx, getOrCreateServiceQuery(), serviceInfo.Name)
	var id int
	if err := row.Scan(&id); err != nil {
		logger.Log().Error("ServiceRepo.GetOrCreate failed", zap.Error(err), zap.String("name", serviceInfo.Name))
		return fmt.Errorf("GetOrCreate - failed to get or create service: %w", err)
	}
	serviceInfo.Id = &id

	logger.Log().Info("ServiceRepo.GetOrCreate ok", zap.String("name", serviceInfo.Name), zap.Int("service_id", id))
	return nil
}
