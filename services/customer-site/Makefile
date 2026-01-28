css: #--watch
	npx @tailwindcss/cli -i ./web/dev/tailwind/styles.css -o ./web/public/css/styles.css --minify

templ-gen: 
	templ generate

build:
	go build -o bin/app ./cmd/site/main.go

prod:
	./bin/app

run:
	@APP_ENV=development go run ./cmd/site/main.go  -c .env

esbuild:
	node esbuild.config.js

resources:
	node generate-resources-paths.js

sprite:
	node generate-svg-sprite.js

dev: templ-gen css esbuild resources sprite run





.PHONY: openapi
openapi:

	@if [ -d ./internal/openapiclient ]; then rm -rf ./internal/openapiclient; fi

# 	curl -sSL https://raw.githubusercontent.com/fedotovmax/microservices-shop/main/services/api-gateway/docs/swagger.yaml -o openapi/swagger.yaml

	openapi-generator generate \
		-i openapi/swagger.yaml \
		-g go \
		-o internal/openapiclient\
  --additional-properties=packageName=openapiclient,withGoMod=false

	find ./internal/openapiclient -type f ! -name "*.go" -exec rm -v {} \;

	rm -rf ./internal/openapiclient/test

	@echo "Openapi client generated successfully! ✅"

