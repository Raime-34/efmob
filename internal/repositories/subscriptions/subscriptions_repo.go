package subscriptions

import (
	"context"
	"efmob/internal/dto"
	dbinterface "efmob/internal/repositories/db_interface"
	"fmt"
	"time"
)

type SubscriptionRepo struct {
	db dbinterface.DbIface
}

func NewSubscriptionRepo(pool dbinterface.DbIface) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: pool,
	}
}

func (r *SubscriptionRepo) CreateSubscriptionInfo(ctx context.Context, subscriptionInfo *dto.SubscriptionInfo) error {
	_, err := r.db.Exec(
		ctx, insertSubscriptionData(),
		subscriptionInfo.UserID,
		subscriptionInfo.ServiceID,
		subscriptionInfo.Price,
		subscriptionInfo.StartData,
		subscriptionInfo.EndData,
	)

	if err != nil {
		return fmt.Errorf("CreateSubscriptionInfo - Failed to insert subscription info: %w", err)
	}

	return nil
}

func (r *SubscriptionRepo) GetSubscriptionInfo(ctx context.Context, subscriptionInfo *dto.SubscriptionInfo) error {
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
		endDate   time.Time
	)

	if err := row.Scan(
		&userID,
		&serviceID,
		&price,
		&startDate,
		&endDate,
	); err != nil {
		return fmt.Errorf("GetSubscriptionInfo - Failed to get subscription info: %w", err)
	}

	subscriptionInfo.Price = price

	return nil
}
