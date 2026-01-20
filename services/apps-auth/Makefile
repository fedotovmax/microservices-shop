NAME?=
TYPE?=


NAME_ARG=$(if $(NAME),-name $(NAME),)
TYPE_ARG=$(if $(TYPE),-type $(TYPE),)

auth-app:
	@APP_ENV=development go run ./cmd/apps-auth/main.go -c .env


cli-app-dev:
	@APP_ENV=development go run ./cmd/cli-app/main.go -c .env $(NAME_ARG) $(TYPE_ARG)

cli-app-prod:
	@APP_ENV=production go run ./cmd/cli-app/main.go $(NAME_ARG) $(TYPE_ARG)