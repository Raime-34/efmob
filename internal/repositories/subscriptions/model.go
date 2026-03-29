package subscriptions

import (
	"fmt"
	"time"

	"efmob/internal/dto"
	"efmob/internal/util"
)

// SubscriptionInfo — сущность подписки в слое хранения (не HTTP DTO).
type SubscriptionInfo struct {
	ServiceID int
	Price     int
	UserID    string
	StartDate time.Time
	EndDate   *time.Time
}

// NewSubscriptionInfoFromCreate собирает модель из тела создания подписки.
func NewSubscriptionInfoFromCreate(req dto.CreateOrUpdateSubscriptionRequest, serviceID int) (*SubscriptionInfo, error) {
	start, err := util.MonthYearToTime(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("start_date: %w", err)
	}

	subInfo := SubscriptionInfo{
		ServiceID: serviceID,
		Price:     req.Price,
		UserID:    req.UserID,
		StartDate: start,
	}

	if req.EndDate != "" {
		end, err := util.MonthYearToTime(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date: %w", err)
		}
		subInfo.EndDate = &end
	}

	return &subInfo, nil
}
