package dto

// CreateSubscriptionRequest — тело POST /subscription (создание подписки).
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// DeleteSubscriptionRequest — тело DELETE /subscription.
type DeleteSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	UserID      string `json:"user_id"`
}

// UpdateSubscriptionRequest — тело PUT /subscription.
type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type ServiceInfo struct {
	Id   *int
	Name string
}
