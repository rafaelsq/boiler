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
	gqlgen

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

start-deps:
	@if [ ! "`docker ps -q -f name=memcached`"  ]; then \
		docker run --name memcached -p 11211:11211 -d --restart=always memcached; \
	fi
