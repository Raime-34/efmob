package subscriptions

import (
	"efmob/internal/dto"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

// go-openapi/testify не содержит assert.AnError (в отличие от stretchr/testify).
var errMockDB = errors.New("mock db error")

func TestSubscriptionRepo_CreateSubscriptionInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		info := &dto.SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 42,
			Price:     999,
			StartDate: "2025-01-01",
			EndDate:   "2025-12-31",
		}

		mock.ExpectExec(regexp.QuoteMeta(insertSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				info.Price,
				info.StartDate,
				info.EndDate,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		repo := NewSubscriptionRepo(mock)
		err = repo.CreateSubscriptionInfo(t.Context(), info)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_GetSubscriptionInfo(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	serviceID := 7
	price := 1500
	start := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)

	t.Run("found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		info := &dto.SubscriptionInfo{
			UserID:    userID,
			ServiceID: serviceID,
		}

		mock.ExpectQuery(regexp.QuoteMeta(getSubscriptionData())).
			WithArgs(userID, serviceID).
			WillReturnRows(
				pgxmock.NewRows([]string{"user_id", "service_id", "price", "start_date", "end_date"}).
					AddRow(userID, serviceID, price, start, end),
			)

		repo := NewSubscriptionRepo(mock)
		err = repo.GetSubscriptionInfo(t.Context(), info)
		assert.NoError(t, err)
		assert.Equal(t, price, info.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		info := &dto.SubscriptionInfo{
			UserID:    userID,
			ServiceID: 999,
		}

		mock.ExpectQuery(regexp.QuoteMeta(getSubscriptionData())).
			WithArgs(userID, 999).
			WillReturnError(pgx.ErrNoRows)

		repo := NewSubscriptionRepo(mock)
		err = repo.GetSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pgx.ErrNoRows))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_DeleteSubscriptionInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		userID := "550e8400-e29b-41d4-a716-446655440000"
		serviceID := 42
		info := &dto.SubscriptionInfo{
			UserID:    userID,
			ServiceID: serviceID,
		}

		mock.ExpectExec(regexp.QuoteMeta(deleteSubscriptionData())).
			WithArgs(userID, serviceID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		repo := NewSubscriptionRepo(mock)
		err = repo.DeleteSubscriptionInfo(t.Context(), info)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		userID := "550e8400-e29b-41d4-a716-446655440000"
		info := &dto.SubscriptionInfo{UserID: userID, ServiceID: 1}

		mock.ExpectExec(regexp.QuoteMeta(deleteSubscriptionData())).
			WithArgs(userID, 1).
			WillReturnError(errMockDB)

		repo := NewSubscriptionRepo(mock)
		err = repo.DeleteSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_UpdateSubscriptionInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		info := &dto.SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 7,
			Price:     200,
			StartDate: "2025-01-15",
			EndDate:   "2025-06-30",
		}

		mock.ExpectExec(regexp.QuoteMeta(updateSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				info.StartDate,
				info.EndDate,
				info.Price,
			).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		repo := NewSubscriptionRepo(mock)
		err = repo.UpdateSubscriptionInfo(t.Context(), info)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		info := &dto.SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 1,
			Price:     1,
			StartDate: "2025-01-01",
			EndDate:   "2025-12-31",
		}

		mock.ExpectExec(regexp.QuoteMeta(updateSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				info.StartDate,
				info.EndDate,
				info.Price,
			).
			WillReturnError(errMockDB)

		repo := NewSubscriptionRepo(mock)
		err = repo.UpdateSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
