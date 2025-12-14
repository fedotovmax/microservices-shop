VERSION=

STEPS ?= 0

app:
	@APP_ENV=development go run ./cmd/user_service/main.go -c .env
migrator-up:
	@go run ./cmd/migrator/main.go -m up -steps $(STEPS) -c .env
migrator-down:
	@go run ./cmd/migrator/main.go -m down -steps $(STEPS) -c .env
migrator-force:
	@go run ./cmd/migrator/main.go -m force -version $(VERSION) -c .env
migrator-version:
	@go run ./cmd/migrator/main.go -m version -c .env