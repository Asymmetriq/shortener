run:
	go run cmd/shortener/main.go

t:
	go test ./...

mocks:
	mockgen -source=internal/repositories/repository.go Repository > internal/test/mocks/repository_mock.go