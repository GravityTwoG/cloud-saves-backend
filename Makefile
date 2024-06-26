
vet:
	go vet ./...

fmt:
	go fmt ./...

app-build:
	go build -o app.bin ./cmd/cloud-saves-backend/main.go

app-run:
	go run ./cmd/cloud-saves-backend/main.go

migrate:
	go run ./cmd/migrate/main.go

test:
	go test ./...

infra:
	cd deployments && docker-compose up

infra-detached:
	cd deployments && docker-compose up -d

infra-down:
	cd deployments && docker-compose down

infra-test:
	cd deployments && docker-compose -f docker-compose.test.yaml up

infra-test-detached:
	cd deployments && docker-compose -f docker-compose.test.yaml up -d

infra-test-down:
	cd deployments && docker-compose -f docker-compose.test.yaml down

openapi:
	swag init -g ./cmd/cloud-saves-backend/main.go

dev: infra-detached migrate infra
