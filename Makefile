.PHONY: help run seed fmt vet tidy lint dev test test-cover cover-html build docker-build docker-run clean

help:
	@echo "Доступные команды:"
	@echo "  make build          - Собрать приложение"
	@echo "  make run            - Запустить приложение"
	@echo "  make seed           - Запустить seed скрипт"
	@echo "  make test           - Запустить тесты"
	@echo "  make test-cover     - Запустить тесты с покрытием"
	@echo "  make cover-html     - Генерировать HTML отчет о покрытии"
	@echo "  make fmt            - Форматировать код"
	@echo "  make vet            - Проверить код (vet)"
	@echo "  make lint           - Проверить код (golangci-lint)"
	@echo "  make tidy           - Очистить go.mod и go.sum"
	@echo "  make dev            - Запустить в режиме разработки (air)"
	@echo "  make docker-build   - Собрать Docker образ"
	@echo "  make docker-run     - Запустить контейнер Docker"
	@echo "  make clean          - Очистить временные файлы"

build:
	go build -o bin/coworking-manager ./cmd/api/main.go

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

test:
	go test -v ./...

test-cover:
	go test -cover ./...

cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-build:
	docker build -t coworking-manager .

docker-run:
	docker run -p 8080:8080 --env-file .env coworking-manager

clean:
	rm -f bin/coworking-manager
	rm -f coverage.out
	go clean
	rm -rf vendor/
