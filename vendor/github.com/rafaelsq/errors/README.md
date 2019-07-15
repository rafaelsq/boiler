# errors [![GoDoc](https://godoc.org/github.com/rafaelsq/errors?status.svg)](http://godoc.org/github.com/rafaelsq/errors) [![Report card](https://goreportcard.com/badge/github.com/rafaelsq/errors)](https://goreportcard.com/report/github.com/rafaelsq/errors)

```golang
package main

import (
	"encoding/json"
	"fmt"

	"github.com/rafaelsq/errors"
)

var (
	ErrA = fmt.Errorf("err a")
)

func responseErrA() error {
	return ErrA
}

func responseErrB(id int) error {
	err := responseErrA()
	if err != nil {
		return errors.New("err b").SetArg("id", id).SetParent(err)
	}

	return nil
}

func main() {
	err := responseErrB(10)
	fmt.Println(err) // "err b; err a"

	fmt.Println(errors.Cause(err) == ErrA) // true

	er, _ := err.(*errors.Error)
	fmt.Printf("caller %s\n", er.Caller) // caller main.go:21

	b, _ := json.Marshal(er.Args)
	fmt.Printf("args %q\n", b) // args {"id": 10}

	fmt.Printf("parent %v\n", *er.Parent) // parent err a

	errs := errors.List(err)
	fmt.Println("len", len(errs))                   // 2
	fmt.Println(">", errs[0])                       // err a
	fmt.Println(">", errs[1])                       // err b; err a
	fmt.Println(">", (errs[1].(*errors.Error)).Msg) // err b
}
```
