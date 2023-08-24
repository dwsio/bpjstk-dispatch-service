ENGINE=main.go
BUILD_DIR=build

run:
	go run ${ENGINE} run -c jmo -p normal
.PHONY: run

run-high:
	go run ${ENGINE} run -c jmo -p high
.PHONY: run-high

build:
	@echo "Building app"
	go build -o ${BUILD_DIR}/app ${ENGINE}
	@echo "Success build app. Your app is ready to use in 'build/' directory."
.PHONY: build

dependency:
	@echo "Downloading all Go dependencies needed"
	go mod download
	go mod verify
	go mod tidy
	@echo "All Go dependencies was downloaded. you can run 'make debug' to compile locally or 'make build' to build app."
.PHONY: dependency

tidy:
	@echo "Synchronize dependency"
	go mod tidy
	@echo "Finish Synchronize dependency"
.PHONY: tidy

lint:
	golangci-lint run ./...
.PHONY: lint

local:
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up -d --build
.PHONY: local

docker-build:
	@echo "Building and run service in container"
	docker build -t dispatch-service .
.PHONY: docker-build