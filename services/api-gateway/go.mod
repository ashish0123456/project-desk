module gitlab.com/Ash3941/project-management/services/api-gateway

go 1.25.0

require (
	github.com/go-chi/chi/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	gitlab.com/Ash3941/project-management/pkg v0.0.0
)

replace gitlab.com/Ash3941/project-management/pkg => ../../pkg
