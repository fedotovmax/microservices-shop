css: #--watch
	npx @tailwindcss/cli -i ./web/dev/tailwind/styles.css -o ./web/public/css/styles.css --minify

templ: 
	templ generate

build:
	go build -o bin/app ./cmd/site/main.go

run:
	./bin/app

dev:
	go run ./cmd/site/main.go

esbuild:
	node esbuild.config.js


.PHONY: openapi
openapi:
	curl -sSL https://raw.githubusercontent.com/fedotovmax/microservices-shop/main/services/api-gateway/docs/swagger.yaml -o openapi/swagger.yaml

	swagger generate client -f openapi/swagger.yaml -A api -t ./internal