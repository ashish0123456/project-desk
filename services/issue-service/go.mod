module gitlab.com/Ash3941/project-management/services/issue-service

go 1.25.0

require (
	github.com/joho/godotenv v1.5.1
	gitlab.com/Ash3941/project-management/gen v0.0.0
	gitlab.com/Ash3941/project-management/pkg v0.0.0
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.30.0
)

replace gitlab.com/Ash3941/project-management/gen => ../../gen
replace gitlab.com/Ash3941/project-management/pkg => ../../pkg
