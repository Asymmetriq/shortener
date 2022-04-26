run:
	go run cmd/shortener/main.go

t:
	go test ./...

mocks:
	mockgen -source=internal/app/service/service.go Repository > internal/app/test/mocks/repository_mock.go