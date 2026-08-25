module gitlab.com/Ash3941/project-management/services/worker-service

go 1.25.0

require (
	github.com/joho/godotenv v1.5.1
	gitlab.com/Ash3941/project-management/pkg v0.0.0
	github.com/redis/go-redis/v9 v9.0.0
)

replace gitlab.com/Ash3941/project-management/pkg => ../../pkg
