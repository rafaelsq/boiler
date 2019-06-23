package mutation

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/rafaelsq/boiler/pkg/errors"
	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/vektah/gqlparser/gqlerror"
)

func NewMutation(us iface.UserService, es iface.EmailService) *Mutation {
	return &Mutation{
		userService:  us,
		emailService: es,
	}
}

type Mutation struct {
	userService  iface.UserService
	emailService iface.EmailService
}

func (m *Mutation) AddUser(ctx context.Context, input entity.AddUserInput) (*entity.User, error) {
	userID, err := m.userService.Add(ctx, input.Name)
	if err != nil {
		errors.Log(err)
		return nil, fmt.Errorf("service failed")
	}

	return &entity.User{ID: userID}, nil
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	if input.UserID <= 0 {
		return nil, &gqlerror.Error{
			Message: "invalid userID",
			Extensions: map[string]interface{}{
				"code": "-1",
			},
		}
	}

	address, err := mail.ParseAddress(input.Address)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "invalid email address",
			Extensions: map[string]interface{}{
				"code": "-2",
			},
		}
	}

	_, err = m.emailService.Add(ctx, input.UserID, address.Address)
	if err != nil {
		errors.Log(err)
		return nil, fmt.Errorf("service failed")
	}

	return &entity.User{ID: input.UserID}, nil
}
