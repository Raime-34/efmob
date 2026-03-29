package subscriptions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	dbinterface "efmob/internal/repositories/db_interface"
	"efmob/internal/constants"
	"efmob/logger"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type SubscriptionRepo struct {
	db dbinterface.DbIface
}

func NewSubscriptionRepo(pool dbinterface.DbIface) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: pool,
	}
}

func (r *SubscriptionRepo) CreateSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
	logger.Log().Info("SubscriptionRepo.CreateSubscriptionInfo call",
		zap.String("user_id", subscriptionInfo.UserID),
		zap.Int("service_id", subscriptionInfo.ServiceID),
		zap.Int("price", subscriptionInfo.Price),
		zap.Time("start_date", subscriptionInfo.StartDate),
		zapOptionalTime("end_date", subscriptionInfo.EndDate),
	)

	_, err := r.db.Exec(
		ctx, insertSubscriptionData(),
		subscriptionInfo.UserID,
		subscriptionInfo.ServiceID,
		subscriptionInfo.Price,
		subscriptionInfo.StartDate,
		subscriptionInfo.EndDate,
	)

	if err != nil {
		logger.Log().Error("SubscriptionRepo.CreateSubscriptionInfo failed", zap.Error(err))
		return fmt.Errorf("CreateSubscriptionInfo - Failed to insert subscription info: %w", err)
	}

	logger.Log().Info("SubscriptionRepo.CreateSubscriptionInfo ok", zap.String("user_id", subscriptionInfo.UserID), zap.Int("service_id", subscriptionInfo.ServiceID))
	return nil
}

func (r *SubscriptionRepo) GetSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
	logger.Log().Info("SubscriptionRepo.GetSubscriptionInfo call",
		zap.String("user_id", subscriptionInfo.UserID),
		zap.Int("service_id", subscriptionInfo.ServiceID),
	)

	row := r.db.QueryRow(
		ctx,
		getSubscriptionData(),
		subscriptionInfo.UserID,
		subscriptionInfo.ServiceID,
	)

	var (
		userID    string
		serviceID int
		price     int
		startDate time.Time
		endDate   sql.NullTime
	)

	if err := row.Scan(
		&userID,
		&serviceID,
		&price,
		&startDate,
		&endDate,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log().Info("SubscriptionRepo.GetSubscriptionInfo not found",
				zap.String("user_id", subscriptionInfo.UserID),
				zap.Int("service_id", subscriptionInfo.ServiceID),
			)
			return fmt.Errorf("GetSubscriptionInfo - %w", constants.ErrSubscriptionNotFound)
		}
		logger.Log().Error("SubscriptionRepo.GetSubscriptionInfo failed", zap.Error(err))
		return fmt.Errorf("GetSubscriptionInfo - Failed to get subscription info: %w", err)
	}

	subscriptionInfo.UserID = userID
	subscriptionInfo.ServiceID = serviceID
	subscriptionInfo.Price = price
	subscriptionInfo.StartDate = startDate
	if endDate.Valid {
		t := endDate.Time
		subscriptionInfo.EndDate = &t
	} else {
		subscriptionInfo.EndDate = nil
	}

	logger.Log().Info("SubscriptionRepo.GetSubscriptionInfo ok",
		zap.String("user_id", userID),
		zap.Int("service_id", serviceID),
		zap.Int("price", price),
		zap.Time("start_date", startDate),
		zapOptionalTime("end_date", subscriptionInfo.EndDate),
	)
	return nil
}

func (r *SubscriptionRepo) ListSubscriptionsByUserID(ctx context.Context, userID string, limit, offset int) ([]ListRow, int, error) {
	logger.Log().Info("SubscriptionRepo.ListSubscriptionsByUserID call",
		zap.String("user_id", userID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	var total int
	if err := r.db.QueryRow(ctx, countSubscriptionsByUserID(), userID).Scan(&total); err != nil {
		logger.Log().Error("SubscriptionRepo.ListSubscriptionsByUserID count failed", zap.Error(err), zap.String("user_id", userID))
		return nil, 0, fmt.Errorf("ListSubscriptionsByUserID - count: %w", err)
	}

	rows, err := r.db.Query(ctx, listSubscriptionsByUserID(), userID, limit, offset)
	if err != nil {
		logger.Log().Error("SubscriptionRepo.ListSubscriptionsByUserID query failed", zap.Error(err), zap.String("user_id", userID))
		return nil, 0, fmt.Errorf("ListSubscriptionsByUserID - query: %w", err)
	}
	defer rows.Close()

	var out []ListRow
	for rows.Next() {
		var (
			name      string
			price     int
			startDate time.Time
			endDate   sql.NullTime
		)
		if err := rows.Scan(&name, &price, &startDate, &endDate); err != nil {
			return nil, 0, fmt.Errorf("ListSubscriptionsByUserID - scan: %w", err)
		}
		row := ListRow{
			ServiceName: name,
			Price:       price,
			StartDate:   startDate,
		}
		if endDate.Valid {
			t := endDate.Time
			row.EndDate = &t
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		logger.Log().Error("SubscriptionRepo.ListSubscriptionsByUserID rows failed", zap.Error(err))
		return nil, 0, fmt.Errorf("ListSubscriptionsByUserID - rows: %w", err)
	}

	logger.Log().Info("SubscriptionRepo.ListSubscriptionsByUserID ok",
		zap.String("user_id", userID),
		zap.Int("rows_count", len(out)),
		zap.Int("total", total),
	)
	return out, total, nil
}

// SumPrice возвращает сумму полей price по подпискам с учётом фильтров (все поля опциональны).
func (r *SubscriptionRepo) SumPrice(ctx context.Context, filters *PriceSumFilters) (int64, error) {
	logger.Log().Info("SubscriptionRepo.SumPrice call", zapPriceSumFilters(filters)...)

	q, args := buildSumPriceQuery(filters)
	var sum int64
	if err := r.db.QueryRow(ctx, q, args...).Scan(&sum); err != nil {
		logger.Log().Error("SubscriptionRepo.SumPrice failed", zap.Error(err), zap.Any("args", args))
		return 0, fmt.Errorf("SumPrice: %w", err)
	}
	logger.Log().Info("SubscriptionRepo.SumPrice ok", zap.Int64("sum_price", sum))
	return sum, nil
}

func (r *SubscriptionRepo) DeleteSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
	logger.Log().Info("SubscriptionRepo.DeleteSubscriptionInfo call",
		zap.String("user_id", subscriptionInfo.UserID),
		zap.Int("service_id", subscriptionInfo.ServiceID),
	)

	tag, err := r.db.Exec(
		ctx,
		deleteSubscriptionData(),
		subscriptionInfo.UserID,
		subscriptionInfo.ServiceID,
	)
	if err != nil {
		logger.Log().Error("SubscriptionRepo.DeleteSubscriptionInfo failed", zap.Error(err))
		return fmt.Errorf("DeleteSubscriptionInfo - Failed to delete subscription info: %w", err)
	}
	if tag.RowsAffected() == 0 {
		logger.Log().Info("SubscriptionRepo.DeleteSubscriptionInfo no rows", zap.String("user_id", subscriptionInfo.UserID), zap.Int("service_id", subscriptionInfo.ServiceID))
		return fmt.Errorf("DeleteSubscriptionInfo - %w", constants.ErrSubscriptionNotFound)
	}

	logger.Log().Info("SubscriptionRepo.DeleteSubscriptionInfo ok",
		zap.String("user_id", subscriptionInfo.UserID),
		zap.Int("service_id", subscriptionInfo.ServiceID),
		zap.Int64("rows_affected", tag.RowsAffected()),
	)
	return nil
}

func (r *SubscriptionRepo) UpdateSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
	logger.Log().Info("SubscriptionRepo.UpdateSubscriptionInfo call",
		zap.String("user_id", subscriptionInfo.UserID),
		zap.Int("service_id", subscriptionInfo.ServiceID),
		zap.Int("price", subscriptionInfo.Price),
		zap.Time("start_date", subscriptionInfo.StartDate),
		zapOptionalTime("end_date", subscriptionInfo.EndDate),
	)

	tag, err := r.db.Exec(
		ctx,
		updateSubscriptionData(),
		subscriptionInfo.UserID,
		subscriptionInfo.ServiceID,
		subscriptionInfo.StartDate,
		subscriptionInfo.EndDate,
		subscriptionInfo.Price,
	)
	if err != nil {
		logger.Log().Error("SubscriptionRepo.UpdateSubscriptionInfo failed", zap.Error(err))
		return fmt.Errorf("UpdateSubscriptionInfo - Failed to update subscription info: %w", err)
	}
	if tag.RowsAffected() == 0 {
		logger.Log().Info("SubscriptionRepo.UpdateSubscriptionInfo no rows",
			zap.String("user_id", subscriptionInfo.UserID),
			zap.Int("service_id", subscriptionInfo.ServiceID),
		)
		return fmt.Errorf("UpdateSubscriptionInfo - %w", constants.ErrSubscriptionNotFound)
	}

	logger.Log().Info("SubscriptionRepo.UpdateSubscriptionInfo ok",
		zap.String("user_id", subscriptionInfo.UserID),
		zap.Int("service_id", subscriptionInfo.ServiceID),
		zap.Int64("rows_affected", tag.RowsAffected()),
	)
	return nil
}

func zapOptionalTime(key string, t *time.Time) zap.Field {
	if t == nil {
		return zap.String(key, "")
	}
	return zap.Time(key, *t)
}

func zapPriceSumFilters(f *PriceSumFilters) []zap.Field {
	if f == nil {
		return []zap.Field{zap.String("filters", "nil")}
	}
	var fs []zap.Field
	if f.UserID != nil {
		fs = append(fs, zap.String("filter_user_id", *f.UserID))
	}
	if f.ServiceName != nil {
		fs = append(fs, zap.String("filter_service_name", *f.ServiceName))
	}
	if f.StartDate != nil {
		fs = append(fs, zap.Time("filter_start_date", *f.StartDate))
	}
	if f.EndDate != nil {
		fs = append(fs, zap.Time("filter_end_date", *f.EndDate))
	}
	if len(fs) == 0 {
		fs = append(fs, zap.String("filters", "empty (all rows)"))
	}
	return fs
}
