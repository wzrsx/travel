.PHONY: all build run

all: build run
# Сборка
build:
	@echo "Building binary..."
	go build -o bin/app_route ./cmd/app_route
	@echo "Binary built."

# Запуск сервера
run:
	@echo "Starting server..."
	./bin/app_route