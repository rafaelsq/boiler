package mutation

import (
	"context"
	"fmt"
	"net/mail"
	"strconv"

	"boiler/cmd/server/internal/graphql/entity"
	lentity "boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/service"
)

// NewMutation return a new Mutation
func NewMutation(service service.Interface) *Mutation {
	return &Mutation{
		service: service,
	}
}

// Mutation handle service mutation
type Mutation struct {
	service service.Interface
}

// AddUser add a new User to the service
func (m *Mutation) AddUser(ctx context.Context, input entity.AddUserInput) (*entity.UserResponse, error) {
	user := lentity.User{
		Name:     input.Name,
		Password: input.Password,
	}

	err := m.service.AddUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("fail to add user; %w", err)
	}

	return &entity.UserResponse{User: &entity.User{ID: strconv.FormatInt(user.ID, 10)}}, nil
}

// AddEmail add a new Email to the service
func (m *Mutation) AddEmail(ctx context.Context, input entity.AddEmailInput) (*entity.EmailResponse, error) {
	userID, err := strconv.ParseInt(input.UserID, 10, 64)
	if err != nil || userID == 0 {
		return nil, errors.ErrInvalidID
	}

	address, err := mail.ParseAddress(input.Address)
	if err != nil {
		return nil, errors.ErrInvalidEmailAddress
	}

	email := lentity.Email{
		UserID:  userID,
		Address: address.Address,
	}

	err = m.service.AddEmail(ctx, &email)
	if err != nil {
		return nil, err
	}

	return &entity.EmailResponse{Email: &entity.Email{ID: strconv.FormatInt(email.ID, 10)}}, nil
}

// AuthUser returns a JWT token
func (m *Mutation) AuthUser(ctx context.Context, input entity.AuthUserInput) (*entity.AuthUserResponse, error) {

	var user lentity.User
	var token string

	err := m.service.AuthUser(ctx, input.Email, input.Password, &user, &token)
	if err != nil {
		return nil, fmt.Errorf("fail to authenticate user; %w", err)
	}

	return &entity.AuthUserResponse{
		Token: token,
		User:  &entity.User{ID: strconv.FormatInt(user.ID, 10)},
	}, nil
}
