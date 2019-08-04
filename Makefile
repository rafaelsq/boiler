run:
	@go run cmd/server/server.go

gen:
	go generate ./...

update-graphql-schema:
	go run github.com/99designs/gqlgen

start-db:
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
