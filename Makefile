tailwind-build:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/css/style.min.css --minify

build:
	~/go/bin/templ generate && go build -o ./tmp/main ./cmd/server
