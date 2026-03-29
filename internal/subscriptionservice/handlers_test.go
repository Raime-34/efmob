package subscriptionservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"efmob/internal/constants"
	"efmob/internal/dto"
	mockrepo "efmob/mocks"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
)

const testUserUUID = "550e8400-e29b-41d4-a716-446655440000"

func newTestService(ctrl *gomock.Controller) (*SubService, *mockrepo.MockRepositories) {
	m := mockrepo.NewMockRepositories(ctrl)
	return &SubService{
		repo:      m,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}, m
}

func TestInsertSubscriptionHandler(t *testing.T) {
	validBody := `{"service_name":"Yandex","price":400,"user_id":"` + testUserUUID + `","start_date":"01-2025"}`
	validReq := dto.CreateOrUpdateSubscriptionRequest{
		ServiceName: "Yandex",
		Price:       400,
		UserID:      testUserUUID,
		StartDate:   "01-2025",
	}

	t.Run("bad_json_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{`))
		rec := httptest.NewRecorder()
		srv.InsertSubscriptionHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var resp dto.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, constants.InsertSubscriptionIncorrectBodyMessage, resp.Message)
	})

	t.Run("validation_422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		body := `{"service_name":"","price":0,"user_id":"bad","start_date":""}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.InsertSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("created_201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			InsertSubscriptionInfo(gomock.Any(), gomock.Eq(validReq)).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.InsertSubscriptionHandler(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		var resp dto.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, constants.Ok, resp.Message)
	})

	t.Run("month_year_error_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		err := fmt.Errorf("wrap: %w", constants.ErrMonthYearInvalidFormat)
		repo.EXPECT().
			InsertSubscriptionInfo(gomock.Any(), gomock.Eq(validReq)).
			Return(err)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.InsertSubscriptionHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var resp dto.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, constants.IncorrectDateFormatMessage, resp.Message)
	})

	t.Run("internal_500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			InsertSubscriptionInfo(gomock.Any(), gomock.Eq(validReq)).
			Return(fmt.Errorf("db down"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.InsertSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDeleteSubscriptionHandler(t *testing.T) {
	bodyJSON := `{"service_name":"Yandex","user_id":"` + testUserUUID + `"}`
	wantReq := dto.DeleteSubscriptionRequest{ServiceName: "Yandex", UserID: testUserUUID}

	t.Run("bad_json_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewBufferString(`x`))
		rec := httptest.NewRecorder()
		srv.DeleteSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ok_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			DeleteSubscriptionInfo(gomock.Any(), gomock.Eq(wantReq)).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewBufferString(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.DeleteSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp dto.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, constants.Ok, resp.Message)
	})

	t.Run("not_found_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			DeleteSubscriptionInfo(gomock.Any(), gomock.Eq(wantReq)).
			Return(fmt.Errorf("del: %w", constants.ErrSubscriptionNotFound))

		req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewBufferString(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.DeleteSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestUpdateSubscriptionHandler(t *testing.T) {
	validBody := `{"service_name":"Yandex","price":400,"user_id":"` + testUserUUID + `","start_date":"01-2025"}`
	validReq := dto.CreateOrUpdateSubscriptionRequest{
		ServiceName: "Yandex",
		Price:       400,
		UserID:      testUserUUID,
		StartDate:   "01-2025",
	}

	t.Run("validation_422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.UpdateSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("ok_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			UpdateSubscriptionInfo(gomock.Any(), gomock.Eq(validReq)).
			Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.UpdateSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not_found_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			UpdateSubscriptionInfo(gomock.Any(), gomock.Eq(validReq)).
			Return(fmt.Errorf("up: %w", constants.ErrSubscriptionNotFound))

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.UpdateSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("month_year_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			UpdateSubscriptionInfo(gomock.Any(), gomock.Eq(validReq)).
			Return(fmt.Errorf("w: %w", constants.ErrMonthYearInvalidFormat))

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(validBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.UpdateSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPatchSubscriptionHandler(t *testing.T) {
	price := 99
	t.Run("nothing_to_patch_422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		body := `{"service_name":"Yandex","user_id":"` + testUserUUID + `"}`
		req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.PatchSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("ok_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		want := dto.PatchSubscriptionRequest{
			ServiceName: "Yandex",
			UserID:      testUserUUID,
			Price:       &price,
		}
		repo.EXPECT().
			PatchSubscriptionInfo(gomock.Any(), gomock.Eq(want)).
			Return(nil)

		body := `{"service_name":"Yandex","user_id":"` + testUserUUID + `","price":99}`
		req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.PatchSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("not_found_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		want := dto.PatchSubscriptionRequest{
			ServiceName: "Yandex",
			UserID:      testUserUUID,
			Price:       ptrInt(1),
		}
		repo.EXPECT().
			PatchSubscriptionInfo(gomock.Any(), gomock.Eq(want)).
			Return(fmt.Errorf("p: %w", constants.ErrSubscriptionNotFound))

		body := `{"service_name":"Yandex","user_id":"` + testUserUUID + `","price":1}`
		req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		srv.PatchSubscriptionHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func ptrInt(v int) *int { return &v }

func TestListSubscriptionsHandler(t *testing.T) {
	t.Run("missing_user_id_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		srv.ListSubscriptionsHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid_uuid_422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/?user_id=not-a-uuid", nil)
		rec := httptest.NewRecorder()
		srv.ListSubscriptionsHandler(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("bad_pagination_422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+testUserUUID+"&page=0", nil)
		rec := httptest.NewRecorder()
		srv.ListSubscriptionsHandler(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("ok_200_with_defaults", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		items := []dto.SubscriptionListItem{{ServiceName: "A", Price: 1, StartDate: "01-2025"}}
		repo.EXPECT().
			ListSubscriptionsByUserID(gomock.Any(), testUserUUID, defaultListPerPage, 0).
			Return(items, 1, nil)

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+testUserUUID, nil)
		rec := httptest.NewRecorder()
		srv.ListSubscriptionsHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp dto.ListSubscriptionsResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, constants.Ok, resp.Message)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, defaultListPerPage, resp.PerPage)
		assert.Equal(t, 1, resp.Total)
		assert.Equal(t, 1, resp.TotalPages)
		assert.Len(t, resp.Subscriptions, 1)
	})

	t.Run("ok_page2", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			ListSubscriptionsByUserID(gomock.Any(), testUserUUID, 10, 10).
			Return([]dto.SubscriptionListItem{}, 25, nil)

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+testUserUUID+"&page=2&per_page=10", nil)
		rec := httptest.NewRecorder()
		srv.ListSubscriptionsHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp dto.ListSubscriptionsResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, 2, resp.Page)
		assert.Equal(t, 10, resp.PerPage)
		assert.Equal(t, 25, resp.Total)
		assert.Equal(t, 3, resp.TotalPages)
	})

	t.Run("repo_error_500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			ListSubscriptionsByUserID(gomock.Any(), testUserUUID, defaultListPerPage, 0).
			Return(nil, 0, context.Canceled)

		req := httptest.NewRequest(http.MethodGet, "/?user_id="+testUserUUID, nil)
		rec := httptest.NewRecorder()
		srv.ListSubscriptionsHandler(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSubscriptionFilteredSumHandler(t *testing.T) {
	t.Run("invalid_user_id_422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, _ := newTestService(ctrl)

		req := httptest.NewRequest(http.MethodGet, "/?user_id=bad", nil)
		rec := httptest.NewRecorder()
		srv.SubscriptionFilteredSumHandler(rec, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("month_year_err_400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			SumSubscriptionPriceByFilter(gomock.Any(), "", "", "xx", "").
			Return(int64(0), fmt.Errorf("s: %w", constants.ErrMonthYearInvalidFormat))

		req := httptest.NewRequest(http.MethodGet, "/?start_date=xx", nil)
		rec := httptest.NewRecorder()
		srv.SubscriptionFilteredSumHandler(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ok_200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			SumSubscriptionPriceByFilter(gomock.Any(), "", "", "", "").
			Return(int64(42), nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		srv.SubscriptionFilteredSumHandler(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp dto.SubscriptionFilteredSumResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, constants.Ok, resp.Message)
		assert.Equal(t, int64(42), resp.SumPrice)
	})

	t.Run("internal_500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srv, repo := newTestService(ctrl)

		repo.EXPECT().
			SumSubscriptionPriceByFilter(gomock.Any(), "", "", "", "").
			Return(int64(0), fmt.Errorf("db"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		srv.SubscriptionFilteredSumHandler(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
