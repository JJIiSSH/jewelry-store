up:
	docker compose -f deployments/docker-compose.yml up -d
down:
	docker compose -f deployments/docker-compose.yml down
run:
	go run cmd/api/main.go
logs:
	docker compose -f deployments/docker-compose.yml logs -f
test:
	go test ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "postgres://IvanDev:1111@localhost:5432/mydb?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://IvanDev:1111@localhost:5432/mydb?sslmode=disable" down 1

proto:
	protoc --go_out=gen --go-grpc_out=gen proto/notification.proto
