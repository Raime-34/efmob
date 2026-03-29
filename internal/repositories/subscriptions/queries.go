package subscriptions

import (
	"fmt"
	"strings"
)

func insertSubscriptionData() string {
	return `
		INSERT INTO subscriptions (user_id, service_id, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`
}

func getSubscriptionData() string {
	return `
		SELECT user_id, service_id, price, start_date, end_date
		FROM subscriptions
		WHERE 
			user_id = $1
			AND 
			service_id = $2
	`
}

func listSubscriptionsByUserID() string {
	return `
		SELECT sv.name, s.price, s.start_date, s.end_date
		FROM subscriptions s
		JOIN services sv ON s.service_id = sv.service_id
		WHERE s.user_id = $1
		ORDER BY sv.name
		LIMIT $2 OFFSET $3
	`
}

func countSubscriptionsByUserID() string {
	return `SELECT COUNT(*) FROM subscriptions WHERE user_id = $1`
}

func deleteSubscriptionData() string {
	return `
		DELETE FROM subscriptions
		WHERE
			user_id = $1
			AND
			service_id = $2
	`
}

func updateSubscriptionData() string {
	return `
		UPDATE subscriptions
		SET
			start_date = $3,
			end_date = $4,
			price = $5
		WHERE
			user_id = $1
			AND 
			service_id = $2
	`
}

// buildSumPriceQuery строит SELECT COALESCE(SUM(sub.price), 0) с опциональными фильтрами.
// Фильтры: user_id, имя сервиса (JOIN services), start_date (>=), end_date (<=).
func buildSumPriceQuery(filters *PriceSumFilters) (string, []any) {
	if filters == nil {
		filters = &PriceSumFilters{}
	}

	var b strings.Builder
	args := []any{}
	pos := 1

	b.WriteString(`SELECT COALESCE(SUM(sub.price), 0) FROM subscriptions AS sub`)
	if filters.ServiceName != nil {
		b.WriteString(` INNER JOIN services sv ON sub.service_id = sv.service_id`)
	}
	b.WriteString(` WHERE `)

	if filters.UserID != nil {
		fmt.Fprintf(&b, "sub.user_id = $%d", pos)
		args = append(args, *filters.UserID)
		pos++
	} else {
		b.WriteString("TRUE")
	}

	if filters.ServiceName != nil {
		fmt.Fprintf(&b, " AND sv.name = $%d", pos)
		args = append(args, *filters.ServiceName)
		pos++
	}
	if filters.StartDate != nil {
		fmt.Fprintf(&b, " AND sub.start_date >= $%d::date", pos)
		args = append(args, *filters.StartDate)
		pos++
	}
	if filters.EndDate != nil {
		fmt.Fprintf(&b, " AND sub.end_date <= $%d::date", pos)
		args = append(args, *filters.EndDate)
		pos++
	}

	return b.String(), args
}
