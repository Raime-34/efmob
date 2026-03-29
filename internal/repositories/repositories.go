package repositories

import (
	"context"
	"strings"

	"efmob/internal/constants"
	"efmob/internal/dto"
	"efmob/internal/repositories/services"
	"efmob/internal/repositories/subscriptions"
	"efmob/internal/util"

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

func (r *Repo) PatchSubscriptionInfo(ctx context.Context, data dto.PatchSubscriptionRequest) error {
	if data.Price == nil && data.StartDate == nil && data.EndDate == nil {
		return constants.ErrPatchNothingToUpdate
	}

	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		return err
	}

	info := &subscriptions.SubscriptionInfo{
		UserID:    data.UserID,
		ServiceID: *serviceInfo.Id,
	}
	if err := r.subscriptionRepo.GetSubscriptionInfo(ctx, info); err != nil {
		return err
	}

	if data.Price != nil {
		info.Price = *data.Price
	}
	if data.StartDate != nil {
		t, err := util.MonthYearToTime(*data.StartDate)
		if err != nil {
			return err
		}
		info.StartDate = t
	}
	if data.EndDate != nil {
		s := strings.TrimSpace(*data.EndDate)
		if s == "" {
			info.EndDate = nil
		} else {
			t, err := util.MonthYearToTime(s)
			if err != nil {
				return err
			}
			info.EndDate = &t
		}
	}

	return r.subscriptionRepo.UpdateSubscriptionInfo(ctx, info)
}

func (r *Repo) ListSubscriptionsByUserID(ctx context.Context, userID string) ([]dto.SubscriptionListItem, error) {
	rows, err := r.subscriptionRepo.ListSubscriptionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.SubscriptionListItem, 0, len(rows))
	for _, row := range rows {
		item := dto.SubscriptionListItem{
			ServiceName: row.ServiceName,
			Price:       row.Price,
			StartDate:   util.TimeToMonthYear(row.StartDate),
		}
		if row.EndDate != nil {
			item.EndDate = util.TimeToMonthYear(*row.EndDate)
		}
		out = append(out, item)
	}
	return out, nil
}
