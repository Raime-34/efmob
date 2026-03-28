package subscriptionservice

import (
	"efmob/internal/dto"
	"efmob/internal/serviceerrors"
	"efmob/logger"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func (s *SubService) InsertSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req dto.CreateSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler - Failed to parse request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.InsertSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler", zap.Error(err))
		switch {
		case serviceerrors.IsMonthYearError(err):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *SubService) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req dto.DeleteSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("DeleteSubscriptionHandler - Failed to parse request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.DeleteSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("DeleteSubscriptionHandler", zap.Error(err))
		switch {
		case errors.Is(err, serviceerrors.ErrSubscriptionNotFound):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func (s *SubService) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req dto.UpdateSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler - Failed to parse request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.UpdateSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler", zap.Error(err))
		switch {
		case serviceerrors.IsMonthYearError(err):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, serviceerrors.ErrSubscriptionNotFound):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}
