package subscriptionservice

import (
	"efmob/internal/constants"
	"efmob/internal/dto"
	"efmob/logger"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func (s *SubService) InsertSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var response dto.Response

	decoder := json.NewDecoder(r.Body)
	var req dto.CreateSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler - Failed to parse request", zap.Error(err))
		response.Message = constants.InsertSubscriptionIncorrectBodyMessage
		w.WriteHeader(http.StatusBadRequest)
		goto response
	}

	if err = s.validator.Struct(&req); err != nil {
		logger.Log().Error("InsertSubscriptionHandler - Validation failed", zap.Error(err))
		response.Message = constants.IncorrectDateFormatMessage
		w.WriteHeader(http.StatusUnprocessableEntity)
		goto response
	}

	err = s.repo.InsertSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler", zap.Error(err))
		switch {
		case constants.IsMonthYearError(err):
			response.Message = constants.IncorrectDateFormatMessage
			w.WriteHeader(http.StatusBadRequest)
		default:
			response.Message = constants.Internal
			w.WriteHeader(http.StatusInternalServerError)
		}
		goto response
	} else {
		response.Message = constants.Ok
		w.WriteHeader(http.StatusCreated)
	}

response:
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(response)
	w.Write(b)
}

func (s *SubService) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var response dto.Response

	decoder := json.NewDecoder(r.Body)
	var req dto.DeleteSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("DeleteSubscriptionHandler - Failed to parse request", zap.Error(err))
		response.Message = constants.DeleteSubscriptionIncorrectBodyMessage
		w.WriteHeader(http.StatusBadRequest)
		goto response
	}

	if err = s.validator.Struct(&req); err != nil {
		logger.Log().Error("DeleteSubscriptionHandler - Validation failed", zap.Error(err))
		response.Message = err.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		goto response
	}

	err = s.repo.DeleteSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("DeleteSubscriptionHandler", zap.Error(err))
		switch {
		case errors.Is(err, constants.ErrSubscriptionNotFound):
			response.Message = constants.ErrSubscriptionNotFound.Error()
			w.WriteHeader(http.StatusBadRequest)
		default:
			response.Message = constants.Internal
			w.WriteHeader(http.StatusInternalServerError)
		}
		goto response
	}

	response.Message = constants.Ok
	w.WriteHeader(http.StatusOK)

response:
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(response)
	w.Write(b)
}

func (s *SubService) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var response dto.Response

	decoder := json.NewDecoder(r.Body)
	var req dto.UpdateSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler - Failed to parse request", zap.Error(err))
		response.Message = constants.UpdateSubscriptionIncorrectBodyMessage
		w.WriteHeader(http.StatusBadRequest)
		goto response
	}

	if err = s.validator.Struct(&req); err != nil {
		logger.Log().Error("UpdateSubscriptionHandler - Validation failed", zap.Error(err))
		response.Message = err.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		goto response
	}

	err = s.repo.UpdateSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler", zap.Error(err))
		switch {
		case constants.IsMonthYearError(err):
			response.Message = constants.IncorrectDateFormatMessage
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, constants.ErrSubscriptionNotFound):
			response.Message = constants.ErrSubscriptionNotFound.Error()
			w.WriteHeader(http.StatusBadRequest)
		default:
			response.Message = constants.Internal
			w.WriteHeader(http.StatusInternalServerError)
		}
		goto response
	}

	response.Message = constants.Ok
	w.WriteHeader(http.StatusOK)

response:
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(response)
	w.Write(b)
}
