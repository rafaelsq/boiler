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

- less interfaces
- import cycles

**Why only one storage?**

- less interfaces

# Run

To start watching your files for modification;

```bash
$ make start-deps
$ make
```

### Dependencies

MySQL and Memcache

pkg/entity or pkg/iface was changed?

```bash
$ make gen
```

./schema.graphql was changed?

```bash
$ make update-graphql-schema
```

> ps; `$ make` will watch and run `make gen and update-graphql-schema` automatically
