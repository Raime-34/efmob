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

func (r *Repo) InsertSubscriptionInfo(ctx context.Context, data dto.CreateSubscriptionDTO) error {
	// TODO чет сделать со стрктурами которые я в слой репозитория отпраляю
	serviceInfo := dto.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		return err
	}

	subscriptionInfo := dto.SubscriptionInfo{
		ServiceID: *serviceInfo.Id,
		Price:     data.Price,
		UserID:    data.UserID,
		StartDate: data.StartDate,
		EndDate:   data.EndDate,
	}
	if err := r.subscriptionRepo.CreateSubscriptionInfo(ctx, &subscriptionInfo); err != nil {
		return err
	}

	return nil
}
