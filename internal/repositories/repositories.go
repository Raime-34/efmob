package repositories

import (
	"context"
	"efmob/internal/dto"
	"efmob/internal/repositories/services"
	"efmob/internal/repositories/subscriptions"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	servicesRepo     *services.ServiceRepo
	subscriptionRepo *subscriptions.SubscriptionRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	return &Repo{
		servicesRepo:     services.NewServiceRepo(conn),
		subscriptionRepo: subscriptions.NewSubscriptionRepo(conn),
	}
}

func (r *Repo) InsertSubscriptionInfo(ctx context.Context, data dto.CreateOrUpdateSubscriptionRequest) error {
	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		return err
	}

	subscriptionInfo, err := subscriptions.NewSubscriptionInfoFromCreate(data, *serviceInfo.Id)
	if err != nil {
		return err
	}

	if err := r.subscriptionRepo.CreateSubscriptionInfo(ctx, subscriptionInfo); err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteSubscriptionInfo(ctx context.Context, data dto.DeleteSubscriptionRequest) error {
	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		return err
	}

	subscriptionInfo := &subscriptions.SubscriptionInfo{
		ServiceID: *serviceInfo.Id,
		UserID:    data.UserID,
	}
	if err := r.subscriptionRepo.DeleteSubscriptionInfo(ctx, subscriptionInfo); err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateSubscriptionInfo(ctx context.Context, data dto.CreateOrUpdateSubscriptionRequest) error {
	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		return err
	}

	subscriptionInfo, err := subscriptions.NewSubscriptionInfoFromCreate(data, *serviceInfo.Id)
	if err != nil {
		return err
	}

	if err := r.subscriptionRepo.UpdateSubscriptionInfo(ctx, subscriptionInfo); err != nil {
		return err
	}

	return nil
}
