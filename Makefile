install:
	go get ./...

start:
	go run main.go

test:
	go test -race ./...

cover:
	go test -cover ./...

local-up:
	docker-compose up

local-down:
	docker-compose down --remove-orphans

configure:
	cp internal/config/example.env internal/config/.env
	swag init