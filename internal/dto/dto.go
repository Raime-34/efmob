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

type PatchSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id" validate:"uuid"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type Response struct {
	Message string `json:"message"`
}
