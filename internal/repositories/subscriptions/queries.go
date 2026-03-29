package subscriptions

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
	`
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
