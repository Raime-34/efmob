package repositories

import (
	"context"
	"strings"

	"efmob/internal/constants"
	"efmob/internal/dto"
	"efmob/internal/repositories/services"
	"efmob/internal/repositories/subscriptions"
	"efmob/internal/util"
	"efmob/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	logger.Log().Info("Repo.InsertSubscriptionInfo call",
		zap.String("service_name", data.ServiceName),
		zap.String("user_id", data.UserID),
		zap.Int("price", data.Price),
		zap.String("start_date", data.StartDate),
		zap.String("end_date", data.EndDate),
	)

	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		logger.Log().Error("Repo.InsertSubscriptionInfo GetOrCreate service failed", zap.Error(err))
		return err
	}

	subscriptionInfo, err := subscriptions.NewSubscriptionInfoFromCreate(data, *serviceInfo.Id)
	if err != nil {
		logger.Log().Error("Repo.InsertSubscriptionInfo NewSubscriptionInfoFromCreate failed", zap.Error(err))
		return err
	}

	if err := r.subscriptionRepo.CreateSubscriptionInfo(ctx, subscriptionInfo); err != nil {
		return err
	}

	logger.Log().Info("Repo.InsertSubscriptionInfo ok", zap.String("user_id", data.UserID), zap.Int("service_id", *serviceInfo.Id))
	return nil
}

func (r *Repo) DeleteSubscriptionInfo(ctx context.Context, data dto.DeleteSubscriptionRequest) error {
	logger.Log().Info("Repo.DeleteSubscriptionInfo call",
		zap.String("service_name", data.ServiceName),
		zap.String("user_id", data.UserID),
	)

	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		logger.Log().Error("Repo.DeleteSubscriptionInfo GetOrCreate service failed", zap.Error(err))
		return err
	}

	subscriptionInfo := &subscriptions.SubscriptionInfo{
		ServiceID: *serviceInfo.Id,
		UserID:    data.UserID,
	}
	if err := r.subscriptionRepo.DeleteSubscriptionInfo(ctx, subscriptionInfo); err != nil {
		return err
	}

	logger.Log().Info("Repo.DeleteSubscriptionInfo ok", zap.String("user_id", data.UserID), zap.Int("service_id", *serviceInfo.Id))
	return nil
}

func (r *Repo) UpdateSubscriptionInfo(ctx context.Context, data dto.CreateOrUpdateSubscriptionRequest) error {
	logger.Log().Info("Repo.UpdateSubscriptionInfo call",
		zap.String("service_name", data.ServiceName),
		zap.String("user_id", data.UserID),
		zap.Int("price", data.Price),
		zap.String("start_date", data.StartDate),
		zap.String("end_date", data.EndDate),
	)

	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		logger.Log().Error("Repo.UpdateSubscriptionInfo GetOrCreate service failed", zap.Error(err))
		return err
	}

	subscriptionInfo, err := subscriptions.NewSubscriptionInfoFromCreate(data, *serviceInfo.Id)
	if err != nil {
		logger.Log().Error("Repo.UpdateSubscriptionInfo NewSubscriptionInfoFromCreate failed", zap.Error(err))
		return err
	}

	if err := r.subscriptionRepo.UpdateSubscriptionInfo(ctx, subscriptionInfo); err != nil {
		return err
	}

	logger.Log().Info("Repo.UpdateSubscriptionInfo ok", zap.String("user_id", data.UserID), zap.Int("service_id", *serviceInfo.Id))
	return nil
}

func (r *Repo) PatchSubscriptionInfo(ctx context.Context, data dto.PatchSubscriptionRequest) error {
	logger.Log().Info("Repo.PatchSubscriptionInfo call",
		zap.String("service_name", data.ServiceName),
		zap.String("user_id", data.UserID),
		zap.Bool("patch_price", data.Price != nil),
		zap.Bool("patch_start_date", data.StartDate != nil),
		zap.Bool("patch_end_date", data.EndDate != nil),
	)

	if data.Price == nil && data.StartDate == nil && data.EndDate == nil {
		return constants.ErrPatchNothingToUpdate
	}

	serviceInfo := services.ServiceInfo{Name: data.ServiceName}
	if err := r.servicesRepo.GetOrCreate(ctx, &serviceInfo); err != nil {
		logger.Log().Error("Repo.PatchSubscriptionInfo GetOrCreate service failed", zap.Error(err))
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
			logger.Log().Error("Repo.PatchSubscriptionInfo parse start_date failed", zap.Error(err))
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
				logger.Log().Error("Repo.PatchSubscriptionInfo parse end_date failed", zap.Error(err))
				return err
			}
			info.EndDate = &t
		}
	}

	if err := r.subscriptionRepo.UpdateSubscriptionInfo(ctx, info); err != nil {
		return err
	}
	logger.Log().Info("Repo.PatchSubscriptionInfo ok", zap.String("user_id", data.UserID), zap.Int("service_id", *serviceInfo.Id))
	return nil
}

func (r *Repo) ListSubscriptionsByUserID(ctx context.Context, userID string, limit, offset int) ([]dto.SubscriptionListItem, int, error) {
	logger.Log().Info("Repo.ListSubscriptionsByUserID call",
		zap.String("user_id", userID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	rows, total, err := r.subscriptionRepo.ListSubscriptionsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
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
	logger.Log().Info("Repo.ListSubscriptionsByUserID ok",
		zap.String("user_id", userID),
		zap.Int("items_count", len(out)),
		zap.Int("total", total),
	)
	return out, total, nil
}

// SumSubscriptionPriceByFilter считает SUM(price) по подпискам. Все параметры опциональны (пустая строка — без фильтра по полю).
func (r *Repo) SumSubscriptionPriceByFilter(ctx context.Context, userID, serviceName, startDate, endDate string) (int64, error) {
	logger.Log().Info("Repo.SumSubscriptionPriceByFilter call",
		zap.String("user_id", userID),
		zap.String("service_name", serviceName),
		zap.String("start_date", startDate),
		zap.String("end_date", endDate),
	)

	filters, err := subscriptions.ComposePriceSumFilters(userID, serviceName, startDate, endDate)
	if err != nil {
		logger.Log().Error("Repo.SumSubscriptionPriceByFilter ComposePriceSumFilters failed", zap.Error(err))
		return 0, err
	}
	sum, err := r.subscriptionRepo.SumPrice(ctx, filters)
	if err != nil {
		return 0, err
	}
	logger.Log().Info("Repo.SumSubscriptionPriceByFilter ok", zap.Int64("sum_price", sum))
	return sum, nil
}
