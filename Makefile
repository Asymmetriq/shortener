run:
	go run cmd/shortener/main.go

t:
	go test ./...

mocks:
	mockgen -source=internal/repository/interface.go Repository > internal/test/mocks/repository_mock.go