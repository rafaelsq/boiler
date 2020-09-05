# Boiler

[![Actions Status](https://github.com/rafaelsq/boiler/workflows/tests/badge.svg)](https://github.com/rafaelsq/boiler/actions)
[![Report card](https://goreportcard.com/badge/github.com/rafaelsq/boiler)](https://goreportcard.com/report/github.com/rafaelsq/boiler)
[![GoDoc](https://godoc.org/github.com/rafaelsq/boiler?status.svg)](http://godoc.org/github.com/rafaelsq/boiler)
<a href="https://github.com/rafaelsq/boiler/blob/master/LICENSE">
<img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License MIT">
</a>

## Project

```
Architeture Layers

┌───────────────┐
├─■ transport ■─┤  e.g. Rest, GraphQL, etc...
├──■ service ■──┤  i.e. business logic
├───■ store ■───┤  i.e. external APIs, MySQL, redis, fileSystem, etc...
└───────────────┘

Project

┌─■ README.md
├─■ Makefile
├─■ mock.sh              // script to generate mock from pkg/ifaces
├─■ .golangci.yml
├─■ .wtc.yaml            // watch settings
│
├─┐cmd
│ └─┐server              // HTTP server
│   ├─■ server.go        // entrypoint
│   └─┐internal
│     ├─┐router
│     │ ├─■ middleware.go
│     │ └─■ router.go      // route to www, rest, graphql, etc..
│     │
│     ├─┐www
│     │ ├─■ handle.go
│     │ └─┐ static
│     │   └─■ *.*
│     │
│     ├─┐rest
│     │ ├─■ handle.go
│     │ └─┐entity
│     │   └─■ *.go
│     │
│     └─┐graphql
│       ├─■ handle.go
│       └─┐graphql
│         ├─■ schema.graphql
│         ├─■ gqlgen.yml
│         ├─■ query.go
│         ├─■ mutation.go
│         ├─■ resolver.go
│         └─┐entity
│           └─■ *.go
│
└─┐pkg
  │
  ├─┐service             // business logic
  │ └─■ *.go
  │
  ├─┐entity              // entities used by service and store
  │ └─■ *.go
  │
  ├─┐iface               // interfaces
  │ ├─■ service.go
  │ ├─■ store.go
  │ └─■ *.go
  │
  ├─┐mock                // mock pkg/iface
  │ └─■ *.go
  │
  └─┐store
    ├─┐config
    │ └─■ config.go
    │
    ├─┐log
    │ └─■ log.go
    │
    └─┐database
      └─■ database.go

```

# Run Dev Mode

```bash
$ make
```

if pkg/entity or pkg/iface changes, it will run `$ make gen` automatically  
if ./schema.graphql changes, it will run `$ make update-graphql-schema` automatically

more info on .wtc.yaml
