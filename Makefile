up:
	docker compose down -v
	docker compose up -d --build

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/subscriptionservice/main.go -o docs --parseDependency --parseInternal

test:
	go test ./...

test_html:
	go test -coverprofile=cover ./internal/...
	go tool cover -html=cover

test_total_cover:
	go test -coverprofile=cover ./internal/...
	go tool cover -func=cover

mock:
	go generate ./...