package main

import (
	"context"
	"fmt"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/usecase"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	ucase := usecase.NewUser( /*db*/ )

	u, err := ucase.ByID(ctx, 1)
	if err != nil {
		return
	}

	fmt.Println("User", u.ID, u.Name)

	us, err := ucase.Friends(ctx, &entity.UserFriendsFilter{
		FromUserID: u.ID,
	})
	if err != nil {
		return
	}

	fmt.Println("Friends")
	for i, u := range us {
		fmt.Println("-", i, u.ID, u.Name)

	}
}
