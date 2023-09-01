# Notification Dispatch Service

How to run in local:

1. Copy the `.env.example` file to `.env` and update the values.
2. Run docker local with rhis command `docker-compose -f docker-compose.local.yaml up -d --build`
3. Run the service `go run main.go`. If want to run the priority use `go run main.go run --priority`.
