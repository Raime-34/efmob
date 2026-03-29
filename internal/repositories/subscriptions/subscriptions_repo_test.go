package subscriptions

import (
	"efmob/internal/constants"
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

		start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		info := &SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 42,
			Price:     999,
			StartDate: start,
			EndDate:   &end,
		}

		mock.ExpectExec(regexp.QuoteMeta(insertSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				info.Price,
				start,
				&end,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		repo := NewSubscriptionRepo(mock)
		err = repo.CreateSubscriptionInfo(t.Context(), info)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_without_end_date", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		start := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
		info := &SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 42,
			Price:     500,
			StartDate: start,
			EndDate:   nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(insertSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				info.Price,
				start,
				(*time.Time)(nil),
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

		info := &SubscriptionInfo{
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

		info := &SubscriptionInfo{
			UserID:    userID,
			ServiceID: 999,
		}

		mock.ExpectQuery(regexp.QuoteMeta(getSubscriptionData())).
			WithArgs(userID, 999).
			WillReturnError(pgx.ErrNoRows)

		repo := NewSubscriptionRepo(mock)
		err = repo.GetSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, constants.ErrSubscriptionNotFound))
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
		info := &SubscriptionInfo{
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
		info := &SubscriptionInfo{UserID: userID, ServiceID: 1}

		mock.ExpectExec(regexp.QuoteMeta(deleteSubscriptionData())).
			WithArgs(userID, 1).
			WillReturnError(errMockDB)

		repo := NewSubscriptionRepo(mock)
		err = repo.DeleteSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		userID := "550e8400-e29b-41d4-a716-446655440000"
		info := &SubscriptionInfo{UserID: userID, ServiceID: 99}

		mock.ExpectExec(regexp.QuoteMeta(deleteSubscriptionData())).
			WithArgs(userID, 99).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		repo := NewSubscriptionRepo(mock)
		err = repo.DeleteSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, constants.ErrSubscriptionNotFound))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriptionRepo_UpdateSubscriptionInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		start := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
		info := &SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 7,
			Price:     200,
			StartDate: start,
			EndDate:   &end,
		}

		mock.ExpectExec(regexp.QuoteMeta(updateSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				start,
				&end,
				info.Price,
			).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		repo := NewSubscriptionRepo(mock)
		err = repo.UpdateSubscriptionInfo(t.Context(), info)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_without_end_date", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		start := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
		info := &SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 7,
			Price:     300,
			StartDate: start,
			EndDate:   nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(updateSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				start,
				(*time.Time)(nil),
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

		start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		info := &SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 1,
			Price:     1,
			StartDate: start,
			EndDate:   &end,
		}

		mock.ExpectExec(regexp.QuoteMeta(updateSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				start,
				&end,
				info.Price,
			).
			WillReturnError(errMockDB)

		repo := NewSubscriptionRepo(mock)
		err = repo.UpdateSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		info := &SubscriptionInfo{
			UserID:    "550e8400-e29b-41d4-a716-446655440000",
			ServiceID: 99,
			Price:     1,
			StartDate: start,
			EndDate:   &end,
		}

		mock.ExpectExec(regexp.QuoteMeta(updateSubscriptionData())).
			WithArgs(
				info.UserID,
				info.ServiceID,
				start,
				&end,
				info.Price,
			).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		repo := NewSubscriptionRepo(mock)
		err = repo.UpdateSubscriptionInfo(t.Context(), info)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, constants.ErrSubscriptionNotFound))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuildSumPriceQuery(t *testing.T) {
	uid := "550e8400-e29b-41d4-a716-446655440000"
	sn := "Netflix"
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	t.Run("nil_filters", func(t *testing.T) {
		q, args := buildSumPriceQuery(nil)
		assert.Contains(t, q, "COALESCE(SUM(sub.price)")
		assert.Contains(t, q, "WHERE TRUE")
		assert.Empty(t, args)
	})

	t.Run("user_only", func(t *testing.T) {
		q, args := buildSumPriceQuery(&PriceSumFilters{UserID: &uid})
		assert.Contains(t, q, "sub.user_id = $1")
		assert.NotContains(t, q, "JOIN services")
		assert.Equal(t, []any{uid}, args)
	})

	t.Run("service_join_and_dates", func(t *testing.T) {
		q, args := buildSumPriceQuery(&PriceSumFilters{
			UserID:      &uid,
			ServiceName: &sn,
			StartDate:   &start,
			EndDate:     &end,
		})
		assert.Contains(t, q, "INNER JOIN services sv")
		assert.Contains(t, q, "sv.name = $2")
		assert.Contains(t, q, "sub.start_date >= $3::date")
		assert.Contains(t, q, "sub.end_date <= $4::date")
		assert.Equal(t, []any{uid, sn, start, end}, args)
	})
}

func TestSubscriptionRepo_SumPrice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		filters := &PriceSumFilters{}
		q, args := buildSumPriceQuery(filters)

		mock.ExpectQuery(regexp.QuoteMeta(q)).
			WithArgs(args...).
			WillReturnRows(
				pgxmock.NewRows([]string{"sum"}).AddRow(int64(12345)),
			)

		repo := NewSubscriptionRepo(mock)
		sum, err := repo.SumPrice(t.Context(), filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(12345), sum)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query_error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		uid := "550e8400-e29b-41d4-a716-446655440000"
		filters := &PriceSumFilters{UserID: &uid}
		q, args := buildSumPriceQuery(filters)

		mock.ExpectQuery(regexp.QuoteMeta(q)).
			WithArgs(args...).
			WillReturnError(errMockDB)

		repo := NewSubscriptionRepo(mock)
		_, err = repo.SumPrice(t.Context(), filters)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
