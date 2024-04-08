# Cloud Saves Backend

## Swagger

Swagger located by url ![localhost:post/swagger/index.html](localhost:port/swagger/index.html)

## Scripts

```sh
go mod tidy

go run .\cmd\cloud-saves-backend\main.go

go run .\cmd\migrate\main.go

go install github.com/swaggo/swag/cmd/swag@latest

swag init -g .\cmd\cloud-saves-backend\main.go
```
