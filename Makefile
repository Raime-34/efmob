up:
	docker compose down -v
	docker compose up -d --build

test_html:
	go test -coverprofile=cover ./internal/...
	go tool cover -html=cover

test_total_cover:
	go test -coverprofile=cover ./internal/...
	go tool cover -func=cover