DATABASE_NAME:=shortener

run:
	go run cmd/shortener/main.go

t:
	go test ./...

mocks:
	mockgen -source=internal/repository/repository.go Repository > internal/test/mocks/repository_mock.go

db-create:
	psql -U postgres -c "drop database if exists $(DATABASE_NAME)"
	psql -U postgres -c "create database $(DATABASE_NAME)"

db-up:
	goose -dir ./internal/database/migrations postgres "${DATABASE_DSN}" up

jet:
	jet -dsn ${DATABASE_DSN} -path=./internal/generated