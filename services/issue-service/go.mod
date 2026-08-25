module gitlab.com/Ash3941/project-management/services/issue-service

go 1.25.0

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	gitlab.com/Ash3941/project-management/gen v0.0.0
	gitlab.com/Ash3941/project-management/pkg v0.0.0
	google.golang.org/grpc v1.59.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace gitlab.com/Ash3941/project-management/gen => ../../gen

replace gitlab.com/Ash3941/project-management/pkg => ../../pkg
