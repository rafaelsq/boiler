watch: godeps
ifeq (, $(shell which wtc))
	go get github.com/rafaelsq/wtc
endif
	@wtc

run: godeps
	@go run cmd/server/server.go

gen: godeps
	go generate ./...

update-graphql-schema: godeps
	gqlgen --config cmd/server/internal/graphql/gqlgen.yml

godeps:
ifeq (, $(shell which msgp))
	go get github.com/tinylib/msgp
endif
ifeq (, $(shell which gqlgen))
	go get github.com/99designs/gqlgen
endif
ifeq (, $(shell which mockgen))
	go get github.com/golang/mock/mockgen
endif
