LOCAL_DB_NAME:=shortener

run:
	go run cmd/shortener/main.go

t:
	go test ./...

mocks:
	mockgen -source=internal/repository/interface.go Repository > internal/test/mocks/repository_mock.go

db-create:
	psql -U postgres -c "drop database if exists $(LOCAL_DB_NAME)"
	psql -U postgres -c "create database $(LOCAL_DB_NAME)"

	