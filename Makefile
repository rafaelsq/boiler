run:
	@go run -mod=vendor cmd/server/server.go

gen:
	go generate ./...

update-graphql-schema:
	go run github.com/99designs/gqlgen

proto:
	@mkdir -p pkg/entity
	@cd pkg/entity && protoc --gogo_out=. *.proto

start-deps:
	@if [ ! "`docker ps -q -f name=memcached`"  ]; then \
		docker run --name memcached -p 11211:11211 -d --restart=always memcached; \
	fi
	@if [ ! "`docker ps -q -f name=boilerdb`"  ]; then \
		if [ "`docker ps -aq -f status=exited -f name=boilerdb`" ]; then \
			docker rm boilerdb; \
		fi; \
		docker run -d -p 3307:3306 --name=boilerdb \
			-v ${PWD}/db:/var/lib/mysql \
			-v ${PWD}/schema.sql:/docker-entrypoint-initdb.d/schema.sql \
			-e MYSQL_ROOT_PASSWORD=boiler \
			mariadb; \
	fi
