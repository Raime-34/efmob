package subscriptionservice

import (
	"efmob/internal/constants"
	"efmob/internal/dto"
	"efmob/logger"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// InsertSubscriptionHandler создаёт подписку.
// @Summary Создать подписку
// @Tags subscription
// @Accept json
// @Produce json
// @Param body body dto.CreateOrUpdateSubscriptionRequest true "Тело запроса"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /subscription [post]
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

// DeleteSubscriptionHandler удаляет подписку.
// @Summary Удалить подписку
// @Tags subscription
// @Accept json
// @Produce json
// @Param body body dto.DeleteSubscriptionRequest true "Тело запроса"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /subscription [delete]
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

// UpdateSubscriptionHandler полностью обновляет подписку.
// @Summary Обновить подписку (PUT)
// @Tags subscription
// @Accept json
// @Produce json
// @Param body body dto.CreateOrUpdateSubscriptionRequest true "Тело запроса"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /subscription [put]
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

// SubscriptionFilteredSumHandler возвращает сумму price по фильтру (все query опциональны, даты MM-YYYY).
// @Summary Сумма подписок по фильтру
// @Tags subscription
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Имя сервиса"
// @Param start_date query string false "MM-YYYY, подписки с start_date не раньше"
// @Param end_date query string false "MM-YYYY, подписки с end_date не позже"
// @Success 200 {object} dto.SubscriptionFilteredSumResponse
// @Failure 400 {object} dto.SubscriptionFilteredSumResponse
// @Failure 422 {object} dto.SubscriptionFilteredSumResponse
// @Failure 500 {object} dto.SubscriptionFilteredSumResponse
// @Router /subscription/sum [get]
func (s *SubService) SubscriptionFilteredSumHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := strings.TrimSpace(q.Get("user_id"))
	serviceName := strings.TrimSpace(q.Get("service_name"))
	startDate := strings.TrimSpace(q.Get("start_date"))
	endDate := strings.TrimSpace(q.Get("end_date"))

	if userID != "" {
		if err := s.validator.Var(userID, "uuid"); err != nil {
			logger.Log().Error("SubscriptionFilteredSumHandler - invalid user_id", zap.Error(err))
			resp := dto.SubscriptionFilteredSumResponse{Message: err.Error(), SumPrice: 0}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			b, _ := json.Marshal(resp)
			w.Write(b)
			return
		}
	}

	sum, err := s.repo.SumSubscriptionPriceByFilter(r.Context(), userID, serviceName, startDate, endDate)
	if err != nil {
		logger.Log().Error("SubscriptionFilteredSumHandler", zap.Error(err))
		resp := dto.SubscriptionFilteredSumResponse{Message: constants.Internal, SumPrice: 0}
		code := http.StatusInternalServerError
		if constants.IsMonthYearError(err) {
			resp.Message = constants.IncorrectDateFormatMessage
			code = http.StatusBadRequest
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	resp := dto.SubscriptionFilteredSumResponse{Message: constants.Ok, SumPrice: sum}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
}

// ListSubscriptionsHandler возвращает список подписок пользователя.
// @Summary Список подписок по user_id
// @Tags subscription
// @Produce json
// @Param user_id query string true "UUID пользователя"
// @Param page query int false "Номер страницы (с 1), по умолчанию 1"
// @Param per_page query int false "Размер страницы 1–100, по умолчанию 20"
// @Success 200 {object} dto.ListSubscriptionsResponse
// @Failure 400 {object} dto.ListSubscriptionsResponse
// @Failure 422 {object} dto.ListSubscriptionsResponse
// @Failure 500 {object} dto.ListSubscriptionsResponse
// @Router /subscription [get]
func (s *SubService) ListSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	empty := func(msg string, code int) {
		resp := dto.ListSubscriptionsResponse{
			Message:         msg,
			Subscriptions:   []dto.SubscriptionListItem{},
			Page:            0,
			PerPage:         0,
			Total:           0,
			TotalPages:      0,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		b, _ := json.Marshal(resp)
		w.Write(b)
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		empty(constants.ListSubscriptionMissingUserIDMessage, http.StatusBadRequest)
		return
	}

	if err := s.validator.Var(userID, "uuid"); err != nil {
		logger.Log().Error("ListSubscriptionsHandler - invalid user_id", zap.Error(err))
		empty(err.Error(), http.StatusUnprocessableEntity)
		return
	}

	page, perPage, ok := parseListPagination(r.URL.Query())
	if !ok {
		empty(constants.ListSubscriptionPaginationInvalidMessage, http.StatusUnprocessableEntity)
		return
	}
	offset := (page - 1) * perPage

	items, total, err := s.repo.ListSubscriptionsByUserID(r.Context(), userID, perPage, offset)
	if err != nil {
		logger.Log().Error("ListSubscriptionsHandler", zap.Error(err))
		empty(constants.Internal, http.StatusInternalServerError)
		return
	}

	resp := dto.ListSubscriptionsResponse{
		Message:         constants.Ok,
		Subscriptions:   items,
		Page:            page,
		PerPage:         perPage,
		Total:           total,
		TotalPages:      totalPages(total, perPage),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
}

// PatchSubscriptionHandler частично обновляет подписку.
// @Summary Частичное обновление (PATCH)
// @Tags subscription
// @Accept json
// @Produce json
// @Param body body dto.PatchSubscriptionRequest true "Тело запроса"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /subscription [patch]
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
