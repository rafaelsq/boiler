# Golang Boilerplate

[![Actions Status](https://github.com/rafaelsq/boiler/workflows/tests/badge.svg)](https://github.com/rafaelsq/boiler/actions)
[![Report card](https://goreportcard.com/badge/github.com/rafaelsq/boiler)](https://goreportcard.com/report/github.com/rafaelsq/boiler)
[![GoDoc](https://godoc.org/github.com/rafaelsq/boiler?status.svg)](http://godoc.org/github.com/rafaelsq/boiler)
<a href="https://github.com/rafaelsq/boiler/blob/master/LICENSE">
<img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License MIT">
</a>

The project should;

- be as explicit as possible
- use interfaces as little as possible
- easy to understand

**Why only one service?**

- prevent import cycles
- less interfaces


# Run Dev Mode

```bash
$ make
```

if pkg/entity or pkg/iface changes, it will run `$ make gen` automatically  
if ./schema.graphql changes, it will run `$ make update-graphql-schema` automatically

more info on .wtc.yaml
