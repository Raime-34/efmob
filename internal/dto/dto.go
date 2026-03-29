package dto

// CreateOrUpdateSubscriptionRequest — тело POST /subscription (создание подписки), PUT /subscription (полное изменение подписки)
type CreateOrUpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
}

// DeleteSubscriptionRequest — тело DELETE /subscription.
type DeleteSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required"`
	UserID      string `json:"user_id" validate:"required,uuid"`
}

// PatchSubscriptionRequest — тело PATCH /subscription (частичное изменение подписки).
// Хотя бы одно из полей price, start_date, end_date должно быть передано.
type PatchSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

// SubscriptionListItem — элемент списка подписок (GET /subscription).
type SubscriptionListItem struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// ListSubscriptionsResponse — ответ GET /subscription?user_id=…&page=&per_page=
type ListSubscriptionsResponse struct {
	Message         string                 `json:"message"`
	Subscriptions   []SubscriptionListItem `json:"subscriptions"`
	Page            int                    `json:"page"`
	PerPage         int                    `json:"per_page"`
	Total           int                    `json:"total"`
	TotalPages      int                    `json:"total_pages"`
}

// SubscriptionFilteredSumResponse — ответ GET /subscription/sum (сумма price по фильтру).
type SubscriptionFilteredSumResponse struct {
	Message  string `json:"message"`
	SumPrice int64  `json:"sum_price"`
}

type Response struct {
	Message string `json:"message"`
}
