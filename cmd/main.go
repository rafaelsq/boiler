package main

import (
	"fmt"

	"github.com/rafaelsq/boiler/pkg/usecase"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	ucase := usecase.NewUser( /*db*/ )

	u, err := ucase.ByID(1)
	if err != nil {
		return
	}

	fmt.Println("User", u.ID, u.Name)
}
