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
	var req dto.CreateOrUpdateSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("InsertSubscriptionHandler - Failed to parse request", zap.Error(err))
		response.Message = constants.InsertSubscriptionIncorrectBodyMessage
		w.WriteHeader(http.StatusBadRequest)
		goto response
	}

	if err = s.validator.Struct(&req); err != nil {
		logger.Log().Error("InsertSubscriptionHandler - Validation failed", zap.Error(err))
		response.Message = constants.ValidationFailedMessage
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
	var req dto.CreateOrUpdateSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("UpdateSubscriptionHandler - Failed to parse request", zap.Error(err))
		response.Message = constants.UpdateSubscriptionIncorrectBodyMessage
		w.WriteHeader(http.StatusBadRequest)
		goto response
	}

	if err = s.validator.Struct(&req); err != nil {
		logger.Log().Error("UpdateSubscriptionHandler - Validation failed", zap.Error(err))
		response.Message = constants.ValidationFailedMessage
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

func (s *SubService) ListSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		resp := dto.ListSubscriptionsResponse{
			Message:       constants.ListSubscriptionMissingUserIDMessage,
			Subscriptions: []dto.SubscriptionListItem{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	if err := s.validator.Var(userID, "uuid"); err != nil {
		logger.Log().Error("ListSubscriptionsHandler - invalid user_id", zap.Error(err))
		resp := dto.ListSubscriptionsResponse{
			Message:       err.Error(),
			Subscriptions: []dto.SubscriptionListItem{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	items, err := s.repo.ListSubscriptionsByUserID(r.Context(), userID)
	if err != nil {
		logger.Log().Error("ListSubscriptionsHandler", zap.Error(err))
		resp := dto.ListSubscriptionsResponse{
			Message:       constants.Internal,
			Subscriptions: []dto.SubscriptionListItem{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	resp := dto.ListSubscriptionsResponse{
		Message:       constants.Ok,
		Subscriptions: items,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
}

func (s *SubService) PatchSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var response dto.Response

	decoder := json.NewDecoder(r.Body)
	var req dto.PatchSubscriptionRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Log().Error("PatchSubscriptionHandler - Failed to parse request", zap.Error(err))
		response.Message = constants.PatchSubscriptionIncorrectBodyMessage
		w.WriteHeader(http.StatusBadRequest)
		goto response
	}

	if err = s.validator.Struct(&req); err != nil {
		logger.Log().Error("PatchSubscriptionHandler - Validation failed", zap.Error(err))
		response.Message = constants.ValidationFailedMessage
		w.WriteHeader(http.StatusUnprocessableEntity)
		goto response
	}

	if req.Price == nil && req.StartDate == nil && req.EndDate == nil {
		response.Message = constants.ErrPatchNothingToUpdate.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)
		goto response
	}

	err = s.repo.PatchSubscriptionInfo(r.Context(), req)
	if err != nil {
		logger.Log().Error("PatchSubscriptionHandler", zap.Error(err))
		switch {
		case errors.Is(err, constants.ErrPatchNothingToUpdate):
			response.Message = err.Error()
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, constants.ErrSubscriptionNotFound):
			response.Message = constants.ErrSubscriptionNotFound.Error()
			w.WriteHeader(http.StatusBadRequest)
		case constants.IsMonthYearError(err):
			response.Message = constants.IncorrectDateFormatMessage
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
