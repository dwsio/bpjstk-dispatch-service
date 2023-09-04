.PHONY:

ENGINE=main.go
BUILD_DIR=build
REGISTRY_URL=172.28.108.245:8082/cns-dev/dispatch-service:v1.0

run:
	go run ${ENGINE} run -c jmo -p normal

run-high:
	go run ${ENGINE} run -c jmo -p high

build:
	@echo "Building app"
	go build -o ${BUILD_DIR}/app ${ENGINE}
	@echo "Success build app. Your app is ready to use in 'build/' directory."

dependency:
	@echo "Downloading all Go dependencies needed"
	go mod download
	go mod verify
	go mod tidy
	@echo "All Go dependencies was downloaded. you can run 'make debug' to compile locally or 'make build' to build app."

tidy:
	@echo "Synchronize dependency"
	go mod tidy
	@echo "Finish Synchronize dependency"

lint:
	golangci-lint run ./...

docker-local:
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up -d --build

docker-dev:
	@echo "Running services in container"
	docker compose -f docker-compose.yaml up -d --build

docker-build:
	@echo "Building service in container"
	docker build -t dispatch_service -f docker/dispatch_service.Dockerfile .

docker-tag:
	@echo Start tagging docker image
	docker tag ${SERVICE_NAME}:latest ${REGISTRY_URL}

docker-push:
	@echo Start push docker image
	docker push ${REGISTRY_URL}
