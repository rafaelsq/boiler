watch: godeps
	@go run cmd/watch/watch.go

run: godeps
	@go run cmd/server/server.go

gen: godeps
	go generate ./...

update-graphql-schema: godeps
	gqlgen

godeps:
ifeq (, $(shell which msgp))
	go get github.com/tinylib/msgp/msgp
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
	@if [ ! "`docker ps -q -f name=boilerdb`"  ]; then \
		if [ "`docker ps -aq -f status=exited -f name=boilerdb`" ]; then \
			docker rm boilerdb; \
		fi; \
		docker run -d -p 3307:3306 --name=boilerdb \
			-v ${PWD}/../boilerdb:/var/lib/mysql \
			-v ${PWD}/schema.sql:/docker-entrypoint-initdb.d/schema.sql \
			-e MYSQL_ROOT_PASSWORD=boiler \
			mariadb; \
	fi
