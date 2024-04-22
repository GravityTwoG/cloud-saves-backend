
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

infra:
	cd deployments && docker-compose up

infra-detached:
	cd deployments && docker-compose up -d


dev: infra-detached migrate infra