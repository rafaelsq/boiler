watch:
	@go run github.com/rafaelsq/wtc

run:
	@go run cmd/server/server.go

gen:
	go generate ./...

update-graphql-schema:
	go run github.com/99designs/gqlgen --config cmd/server/internal/graphql/gqlgen.yml
