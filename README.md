# Golang Boilerplate

[![GoDoc](https://godoc.org/github.com/rafaelsq/boiler?status.svg)](http://godoc.org/github.com/rafaelsq/boiler)
[![Report card](https://goreportcard.com/badge/github.com/rafaelsq/boiler)](https://goreportcard.com/report/github.com/rafaelsq/boiler)
[![Actions Status](https://github.com/rafaelsq/boiler/workflows/tests/badge.svg)](https://github.com/rafaelsq/boiler/actions)
<a href="https://opensource.org/licenses/MIT">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License MIT">
</a>

### Why pkg/iface
so more than one place can implement the interface

### Why another boilerplate?
"Now is better than never."  
"There should be one-- and preferably only one --obvious way to do it."

### Why only one service?
"Simple is better than complex."
> so I don't need to worry about import cycle

### Why injecting storage and service all around?
"Explicit is better than implicit."

### Why only one storage?
"If the implementation is easy to explain, it may be a good idea."
> who uses storage don't need to know what happens inside the storage

### Why cache implements storage?
"Readability counts."

### Why context on error?
"In the face of ambiguity, refuse the temptation to guess."

### Where did I read all those quotes before?
The Zen of Python

# Run

You will need memcache and mysql/mariaDB running.  
You can start using Docker running; `$ make start-deps`.

To start watching your files for modification;
```bash
$ make
```

# Watch
You can use the watch out of this project
```bash
$ go get -u github.com/rafaelsq/boiler/cmd/watch
$ cd my_go_project
$ watch "main.go" "./my_go_project"
// default
watch "cmd/server/server.go" "./server"
```

### Dependencies
MySQL and Memcache  
You can run it with `$ make start-deps`

You'll need Docker and Protobuf.  
To install Protobuf; https://github.com/protocolbuffers/protobuf/releases and `$ go get github.com/gogo/protobuf/protoc-gen-gofast`  
put the binary anywhere in your path  
put `include` folder under ~/go/  

Watch will run [golangci-lint](https://github.com/golangci/golangci-lint) for each file modified  

pkg/entity was changed?
```bash
$ make proto
```

pkg/iface was changed?
```bash
$ make gen
```

./schema.graphql was changed?
```bash
$ make update-graphql-schema
```
> ps; `$ make` will watch and run it automatically
