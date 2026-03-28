package services

import (
	"efmob/internal/dto"
	"errors"
	"regexp"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func TestServiceRepo_CreateService(t *testing.T) {
	t.Run("correct insert", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		serviceId := 230
		serviceName := "test-service"
		inputData := dto.ServiceInfo{
			Name: serviceName,
		}
		expServiceData := dto.ServiceInfo{
			Id:   &serviceId,
			Name: serviceName,
		}

		mock.ExpectQuery(regexp.QuoteMeta(insertServiceQuery())).
			WithArgs(
				expServiceData.Name,
			).
			WillReturnRows(
				pgxmock.NewRows([]string{"service_id"}).
					AddRow(*expServiceData.Id),
			)

		repo := NewServiceRepo(mock)
		err = repo.CreateService(
			t.Context(),
			&inputData,
		)
		assert.Nil(t, err)
		assert.NotNil(t, inputData.Id)
		assert.Equal(t, serviceId, *inputData.Id)
		assert.Equal(t, serviceName, inputData.Name)
	})
}

func TestServiceRepo_GetService(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		serviceID := 230
		serviceName := "test-service"
		inputData := dto.ServiceInfo{Name: serviceName}

		mock.ExpectQuery(regexp.QuoteMeta(getServiceQuery())).
			WithArgs(serviceName).
			WillReturnRows(
				pgxmock.NewRows([]string{"service_id"}).AddRow(serviceID),
			)

		repo := NewServiceRepo(mock)
		err = repo.GetService(t.Context(), &inputData)
		assert.NoError(t, err)
		assert.NotNil(t, inputData.Id)
		assert.Equal(t, serviceID, *inputData.Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		unknownName := "no-such-service"
		inputData := dto.ServiceInfo{Name: unknownName}

		mock.ExpectQuery(regexp.QuoteMeta(getServiceQuery())).
			WithArgs(unknownName).
			WillReturnError(pgx.ErrNoRows)

		repo := NewServiceRepo(mock)
		err = repo.GetService(t.Context(), &inputData)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pgx.ErrNoRows), "expected wrapped pgx.ErrNoRows")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestServiceRepo_GetOrCreate(t *testing.T) {
	t.Run("returns_service_id", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}

		serviceID := 42
		serviceName := "billing"

		mock.ExpectQuery(regexp.QuoteMeta(getOrCreateServiceQuery())).
			WithArgs(serviceName).
			WillReturnRows(
				pgxmock.NewRows([]string{"get_or_create_service"}).AddRow(serviceID),
			)

		info := &dto.ServiceInfo{Name: serviceName}
		repo := NewServiceRepo(mock)
		err = repo.GetOrCreate(t.Context(), info)
		assert.NoError(t, err)
		assert.NotNil(t, info.Id)
		assert.Equal(t, serviceID, *info.Id)
		assert.Equal(t, serviceName, info.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
