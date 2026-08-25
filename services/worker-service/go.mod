module gitlab.com/Ash3941/project-management/services/worker-service

go 1.25.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.0.0
	gitlab.com/Ash3941/project-management/pkg v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace gitlab.com/Ash3941/project-management/pkg => ../../pkg
