package subscriptions

import (
	"fmt"
	"time"

	"efmob/internal/dto"
	"efmob/internal/util"
)

// ListRow — строка списка подписок пользователя (до маппинга в dto).
type ListRow struct {
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     *time.Time
}

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

type PriceSumFilters struct {
	ServiceName *string
	UserID      *string
	StartDate   *time.Time
	EndDate     *time.Time
}

func ComposePriceSumFilters(
	userID string,
	serviceName string,
	startDate string,
	endDate string,
) (*PriceSumFilters, error) {
	filters := PriceSumFilters{}

	if userID != "" {
		filters.UserID = &userID
	}
	if serviceName != "" {
		filters.ServiceName = &serviceName
	}
	if startDate != "" {
		start, err := util.MonthYearToTime(startDate)
		if err != nil {
			return nil, err
		}
		filters.StartDate = &start
	}
	if endDate != "" {
		end, err := util.MonthYearToTime(endDate)
		if err != nil {
			return nil, err
		}
		filters.EndDate = &end
	}

	return &filters, nil
}
