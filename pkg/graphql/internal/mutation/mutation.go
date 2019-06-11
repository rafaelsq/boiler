package mutation

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
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
		return nil, err
	}

	return &entity.User{ID: userID}, nil
}

func (m *Mutation) AddMail(ctx context.Context, input entity.AddMailInput) (*entity.User, error) {
	_, err := m.emailService.Add(ctx, input.UserID, input.Address)
	if err != nil {
		return nil, err
	}

	return &entity.User{ID: input.UserID}, nil
}
