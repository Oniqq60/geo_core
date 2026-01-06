.PHONY: build run test clean docker-build docker-up docker-down migrate

# Сборка приложения
build:
	go build -o bin/server ./cmd/server

# Запуск приложения
run:
	go run ./cmd/server/main.go

# Запуск тестов
test:
	go test -v ./...

# Очистка
clean:
	rm -rf bin/

# Сборка Docker образа
docker-build:
	docker-compose build

# Запуск через Docker Compose
docker-up:
	docker-compose up -d

# Остановка Docker Compose
docker-down:
	docker-compose down

# Применение миграций
migrate:
	psql -U postgres -d geo_system -f migrations/001_init.sql

# Установка зависимостей
deps:
	go mod download
	go mod tidy

