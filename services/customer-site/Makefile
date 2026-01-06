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