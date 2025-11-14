include .env
export


VERSION=

STEPS ?= 0

grpcapp:
	@KAFKA_BROKERS=$(KAFKA_BROKERS) DB_URL=$(USER_SERVICE_POSTGRES_URL) PORT=$(USER_SERVICE_PORT) go run ./cmd/user_service/main.go
migrator-up:
	@DB_URL=$(USER_SERVICE_POSTGRES_URL) go run ./cmd/migrator/main.go -m up -steps $(STEPS)
migrator-down:
	@DB_URL=$(USER_SERVICE_POSTGRES_URL) go run ./cmd/migrator/main.go -m down -steps $(STEPS)
migrator-force:
	@DB_URL=$(USER_SERVICE_POSTGRES_URL) go run ./cmd/migrator/main.go -m force -version $(VERSION)
migrator-version:
	@DB_URL=$(USER_SERVICE_POSTGRES_URL) go run ./cmd/migrator/main.go -m version