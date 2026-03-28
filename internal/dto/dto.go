package dto

type CreateSubscriptionDTO struct {
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

type SubscriptionInfo struct {
	ServiceID int
	Price     int
	UserID    string
	StartDate string
	EndDate   string
}
