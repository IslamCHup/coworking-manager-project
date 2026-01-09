run:
	go run cmd/api/main.go

seed:
	go run cmd/seed/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	golangci-lint run

dev:
	air
# Запуск тестов
test:
	go test -v ./...

# Запуск тестов с покрытием
test-cover:
	go test -cover ./...

# Генерация HTML-отчёта о покрытии
cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
