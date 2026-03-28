package subscriptions

func insertSubscriptionData() string {
	return `
		INSERT INTO subscriptions (user_id, service_id, price, start_data, end_data)
		VALUES ($1, $2, $3, $4, $5)
	`
}

func getSubscriptionData() string {
	return `
		SELECT *
		FROM subscriptions
		WHERE user_id = $1
			AND service_id = $2
	`
}
