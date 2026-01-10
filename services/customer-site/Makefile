css: #--watch
	npx @tailwindcss/cli -i ./web/dev/tailwind/styles.css -o ./web/public/css/styles.css --minify

templ-gen: 
	templ generate

build:
	go build -o bin/app ./cmd/site/main.go

prod:
	./bin/app

run:
	go run ./cmd/site/main.go

esbuild:
	node esbuild.config.js

resources:
	node generate-resources-paths.js

sprite:
	node generate-svg-sprite.js

dev: templ-gen css esbuild resources sprite run


.PHONY: openapi
openapi:

	@if [ -d ./internal/client ]; then rm -rf ./internal/client; fi
	@if [ -d ./internal/models ]; then rm -rf ./internal/models; fi

	curl -sSL https://raw.githubusercontent.com/fedotovmax/microservices-shop/main/services/api-gateway/docs/swagger.yaml -o openapi/swagger.yaml

	swagger generate client -f openapi/swagger.yaml -A api -t ./internal