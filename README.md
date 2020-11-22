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
├─■ .golangci.yml
├─■ .wtc.yaml            // watch settings
│
├─┐cmd
│ ├─■ cmd.go             // common funciont
│ │
│ ├─┐worker              // handle async operation
│ │ ├─■ worker.go
│ │ └─┐internal
│ │   └─┐handle
│ │     └─■ handle.go
│ │
│ └─┐server                // HTTP server
│   ├─■ server.go          // entrypoint
│   └─┐internal
│     ├─┐router
│     │ ├─■ middleware.go
│     │ └─■ router.go      // route www, rest, graphql, etc..
│     │
│     ├─┐www
│     │ ├─■ handle.go
│     │ └─┐ static
│     │   └─■ *.*
│     │
│     ├─┐rest
│     │ ├─■ handle.go
│     │ └─┐entity
│     │   └─■ <handle_name>.go  // payload and response definitions
│     │
│     └─┐graphql
│       ├─■ handle.go
│       ├─■ schema.graphql
│       ├─■ gqlgen.yml
│       ├─■ query.go
│       ├─■ mutation.go
│       ├─■ resolver.go
│       └─┐entity
│         └─■ *.go
│
└─┐pkg
  │
  ├─┐service             // business logic
  │ ├─■ service.go       // new & interface
  │ ├─■ *.go
  │ └─┐mock              // service mock
  │   └─■ *.go
  │
  ├─┐entity
  │ └─■ *.go
  │
  └─┐store
    ├─■ interface.go
    │
    ├─┐mock
    │ └─■ *.go
    │
    ├─┐cfg
    │ └─■ cfg.go
    │
    ├─┐log
    │ └─■ log.go
    │
    └─┐<db>
      └─■ <db>.go
```

# Requirements

Worker requires a running Redis server.  

You can easily start a redis server using docker;
`docker run -d --name=redis -p 6379:6379  redis:6`


# Run Dev Mode

```bash
$ make
```

more info on .wtc.yaml
