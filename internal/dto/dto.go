package dto

type CreateSubscriptionDTO struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartData   string `json:"start_date"`
	EndData     string `json:"end_data"`
}

type ServiceInfo struct {
	Id   *int
	Name string
}

type SubscriptionInfo struct {
	ServiceID int
	Price     int
	UserID    string
	StartData string
	EndData   string
}
