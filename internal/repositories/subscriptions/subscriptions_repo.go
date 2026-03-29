package subscriptions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	dbinterface "efmob/internal/repositories/db_interface"
	"efmob/internal/constants"

	"github.com/jackc/pgx/v5"
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

func (r *SubscriptionRepo) GetSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
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
			return fmt.Errorf("GetSubscriptionInfo - %w", constants.ErrSubscriptionNotFound)
		}
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

	return nil
}

func (r *SubscriptionRepo) ListSubscriptionsByUserID(ctx context.Context, userID string) ([]ListRow, error) {
	rows, err := r.db.Query(ctx, listSubscriptionsByUserID(), userID)
	if err != nil {
		return nil, fmt.Errorf("ListSubscriptionsByUserID - query: %w", err)
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
			return nil, fmt.Errorf("ListSubscriptionsByUserID - scan: %w", err)
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
		return nil, fmt.Errorf("ListSubscriptionsByUserID - rows: %w", err)
	}

	return out, nil
}

func (r *SubscriptionRepo) DeleteSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
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
		return fmt.Errorf("DeleteSubscriptionInfo - %w", constants.ErrSubscriptionNotFound)
	}

	return nil
}

func (r *SubscriptionRepo) UpdateSubscriptionInfo(ctx context.Context, subscriptionInfo *SubscriptionInfo) error {
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
		return fmt.Errorf("UpdateSubscriptionInfo - %w", constants.ErrSubscriptionNotFound)
	}

	return nil
}
