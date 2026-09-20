# ProjectDesk

ProjectDesk is a microservices-based project management platform built with Go. It provides authentication, project management, issue management, an API gateway, asynchronous worker processing, and service-owned persistence.

## Features

- User registration and login with bcrypt password hashing.
- JWT access and refresh tokens.
- Project creation and project membership management.
- Issue creation with status and priority defaults.
- HTTP APIs for client-facing operations.
- gRPC contracts and servers for internal service communication.
- Dedicated PostgreSQL database for each business service.
- Redis-backed worker service for asynchronous jobs.
- Structured logging, health endpoints, graceful shutdown, and automatic migrations.
- Docker Compose setup for local development.

## Architecture

```text
Client

API Gateway :8080
    |--------> Auth Service    HTTP :8081 | gRPC :50051 | PostgreSQL
    |--------> Project Service HTTP :8082 | gRPC :50052 | PostgreSQL
    |--------> Issue Service   HTTP :8083 | gRPC :50053 | PostgreSQL

Worker Service <-------- Redis :6379
```

Each business service owns its application logic, HTTP handlers, gRPC handlers, database store, and migrations. The shared `pkg` module contains PostgreSQL, logging, and migration helpers. Protobuf source contracts are stored in `proto`, with generated Go code in `gen`.

## Repository Structure

```text
├── deployments/docker/       Docker Compose deployment definition
├── gen/                      Generated protobuf and gRPC Go code
├── pkg/                      Shared Go packages
├── proto/                    Protobuf service contracts
├── services/
│   ├── api-gateway/          Public HTTP gateway and reverse proxies
│   ├── auth-service/         Authentication and JWT service
│   ├── issue-service/        Issue management service
│   ├── project-service/      Project and membership service
│   └── worker-service/       Redis-backed worker process
├── docker-compose.yml        Root Compose entry point
├── go.work                   Go workspace definition
└── Makefile                  Common development commands
```

## Prerequisites

- Docker Desktop with Docker Compose.
- Go 1.26 or later for local Go development.
- Buf and protobuf Go plugins only when regenerating gRPC code.

## Run With Docker Compose

The root Compose file is the recommended entry point. It includes the complete deployment definition in `deployments/docker/docker-compose.yml`, so you do not need to start each service separately.

For Oracle VM deployment through GitLab CI/CD, see [`deployments/README.md`](deployments/README.md). The production Compose definition is separate from the local development definition and pulls images from the GitLab Container Registry.

From the repository root:

```powershell
docker compose up --build
```

Run in the background:

```powershell
docker compose up --build -d
```

View logs:

```powershell
docker compose logs -f
```

Stop the stack:

```powershell
docker compose down
```

Stop the stack and remove persisted PostgreSQL volumes:

```powershell
docker compose down -v
```

Each service applies its own SQL migrations during startup. The PostgreSQL data volumes preserve the databases between normal restarts.

## Service Ports

| Component | Host Port |
| --- | ---: |
| API Gateway | 8080 |
| Auth Service HTTP | 8081 |
| Auth Service gRPC | 50051 |
| Project Service HTTP | 8082 |
| Project Service gRPC | 50052 |
| Issue Service HTTP | 8083 |
| Issue Service gRPC | 50053 |
| Auth PostgreSQL | 5432 |
| Project PostgreSQL | 5433 |
| Issue PostgreSQL | 5434 |
| Redis | 6379 |

## Health Checks

```powershell
curl.exe http://localhost:8080/health/live
curl.exe http://localhost:8080/health/ready
curl.exe http://localhost:8080/api/v1/version
```

## HTTP API Examples

### Register a user

```powershell
curl.exe -X POST http://localhost:8081/register `
    -H "Content-Type: application/json" `
    -d '{"email":"user@example.com","password":"password123"}'
```

### Log in

```powershell
curl.exe -X POST http://localhost:8081/login `
    -H "Content-Type: application/json" `
    -d '{"email":"user@example.com","password":"password123"}'
```

The login response contains an access token, refresh token, and access-token expiration time.

### Create a project

```powershell
curl.exe -X POST http://localhost:8082/projects `
    -H "Content-Type: application/json" `
    -d '{"name":"Demo Project"}'
```

### Add a project member

Replace `PROJECT_ID` and `USER_ID` with values returned by the relevant operations.

```powershell
curl.exe -X POST http://localhost:8082/projects/members `
    -H "Content-Type: application/json" `
    -d '{"project_id":"PROJECT_ID","user_id":"USER_ID","role":"member"}'
```

### Create an issue

```powershell
curl.exe -X POST http://localhost:8083/issues `
    -H "Content-Type: application/json" `
    -d '{"project_id":"PROJECT_ID","title":"First issue","body":"Issue details"}'
```

Issues are created with `open` status and `normal` priority by default.

## Configuration

Copy `.env.example` to `.env` for local non-Docker execution:

```powershell
Copy-Item .env.example .env
```

Docker Compose supplies container-specific database hostnames and Redis configuration. The main configuration variables include:

- `APP_ENV`, `LOG_LEVEL`, and `VERSION`.
- `GATEWAY_PORT`.
- `AUTH_*`, `PROJECT_*`, and `ISSUE_*` database settings.
- `AUTH_GRPC_PORT`, `PROJECT_GRPC_PORT`, and `ISSUE_GRPC_PORT`.
- `WORKER_COUNT` and `REDIS_ADDR`.
- `AUTH_JWT_SECRET` for JWT signing.

## Local Go Commands

The Makefile contains shortcuts for building and running individual modules:

```powershell
make build
make run-gateway
make run-auth
make run-project
make run-issue
make run-worker
```

When running services directly with Go, PostgreSQL and Redis must be available separately, and the environment variables must point to their local host ports. Docker Compose is the complete local runtime and does not require separate service terminals.

## Protobuf Generation

After changing a `.proto` contract, regenerate the Go code from the repository root:

```powershell
buf generate
```

Generated clients and server interfaces are written under `gen/proto`.
