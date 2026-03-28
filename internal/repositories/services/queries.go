package services

func insertServiceQuery() string {
	return `
		INSERT INTO services (name)
		VALUES ($1)
		RETURNING service_id
	`
}

func getServiceQuery() string {
	return `
		SELECT service_id
		FROM services
		WHERE name = $1
	`
}

func getOrCreateServiceQuery() string {
	return `
		SELECT get_or_create_service($1)
	`
}
