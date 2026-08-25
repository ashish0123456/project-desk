.PHONY: tidy build test run-gateway run-auth run-project run-issue run-worker docker-up docker-down docker-logs migrate-auth migrate-project migrate-issue migrate-all

COMPOSE_FILE := deployments/docker/docker-compose.yml
MIGRATE ?= migrate

tidy:
	go work sync
	cd pkg && go mod tidy
	cd services/api-gateway && go mod tidy
	cd services/auth-service && go mod tidy
	cd services/project-service && go mod tidy
	cd services/issue-service && go mod tidy
	cd services/worker-service && go mod tidy

build:
	go build -o bin/api-gateway ./services/api-gateway/cmd/server
	go build -o bin/auth-service ./services/auth-service/cmd/server
	go build -o bin/project-service ./services/project-service/cmd/server
	go build -o bin/issue-service ./services/issue-service/cmd/server
	go build -o bin/worker-service ./services/worker-service/cmd/worker

test:
	go test ./...

run-gateway:
	go run ./services/api-gateway/cmd/server

run-auth:
	go run ./services/auth-service/cmd/server

run-project:
	go run ./services/project-service/cmd/server

run-issue:
	go run ./services/issue-service/cmd/server

run-worker:
	go run ./services/worker-service/cmd/worker

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down

docker-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

migrate-auth:
	$(MIGRATE) -path services/auth-service/migrations -database "postgres://postgres:postgres@localhost:5432/pm_auth?sslmode=disable" up

migrate-project:
	$(MIGRATE) -path services/project-service/migrations -database "postgres://postgres:postgres@localhost:5433/pm_project?sslmode=disable" up

migrate-issue:
	$(MIGRATE) -path services/issue-service/migrations -database "postgres://postgres:postgres@localhost:5434/pm_issue?sslmode=disable" up

migrate-all: migrate-auth migrate-project migrate-issue
