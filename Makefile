tailwind-build:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/css/style.min.css --minify

build:
	templ generate && npx tailwindcss -i ./assets/css/input.css -o ./assets/css/style.min.css --minify && go build -o ./tmp/main ./cmd/server

create-db:
	psql -h localhost -U postgres -w -c "create database echo_scaffold;"
