# [WIP] Golang Boilerplate

[![GoDoc](https://godoc.org/github.com/rafaelsq/boiler?status.svg)](http://godoc.org/github.com/rafaelsq/boiler) [![Report card](https://goreportcard.com/badge/github.com/rafaelsq/boiler)](https://goreportcard.com/report/github.com/rafaelsq/boiler)

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

To start the server
```bash
$ make start-deps
$ make
```

### Dependencies
MySQL and Memcache  
You can run it with `$ make start-deps`

You'll need Docker and Protobuf.  
To install Protobuf; https://github.com/protocolbuffers/protobuf/releases and `$ go get github.com/gogo/protobuf/protoc-gen-gofast`  
put the binary anywhere in your path  
put `include` folder under ~/go/

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
