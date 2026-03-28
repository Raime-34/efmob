package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"efmob/internal/dto"
	dbinterface "efmob/internal/repositories/db_interface"
)

// ErrSubscriptionNotFound — DELETE/UPDATE не затронули ни одной строки (записи нет).
var ErrSubscriptionNotFound = errors.New("subscriptions: no matching subscription row")

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
		subscriptionInfo.StartDate,
		subscriptionInfo.EndDate,
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

func (r *SubscriptionRepo) DeleteSubscriptionInfo(ctx context.Context, subscriptionInfo *dto.SubscriptionInfo) error {
	tag, err := r.db.Exec(
		ctx,
		deleteSubscriptionData(),
		subscriptionInfo.UserID,
		subscriptionInfo.ServiceID,
	)
	if err != nil {
		return fmt.Errorf("DeleteSubscriptionInfo - Failed to delete subscription info: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteSubscriptionInfo - %w", ErrSubscriptionNotFound)
	}

	return nil
}

func (r *SubscriptionRepo) UpdateSubscriptionInfo(ctx context.Context, subscriptionInfo *dto.SubscriptionInfo) error {
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
		return fmt.Errorf("UpdateSubscriptionInfo - Failed to update subscription info: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateSubscriptionInfo - %w", ErrSubscriptionNotFound)
	}

	return nil
}
