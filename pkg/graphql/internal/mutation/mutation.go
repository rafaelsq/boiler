package mutation

import (
	"context"
	"fmt"
	"net/mail"
	"strconv"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/log"
	"github.com/rafaelsq/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// NewMutation return a new Mutation
func NewMutation(service iface.Service) *Mutation {
	return &Mutation{
		service: service,
	}
}

// Mutation handle service mutation
type Mutation struct {
	service iface.Service
}

// AddUser add a new User to the service
func (m *Mutation) AddUser(ctx context.Context, input entity.AddUserInput) (*entity.UserResponse, error) {
	userID, err := m.service.AddUser(ctx, input.Name, input.Password)
	if err != nil {
		log.Log(errors.New("fail to add user").SetParent(err))
		return nil, fmt.Errorf("service failed")
	}

	return &entity.UserResponse{User: &entity.User{ID: strconv.FormatInt(userID, 10)}}, nil
}

// AddEmail add a new Email to the service
func (m *Mutation) AddEmail(ctx context.Context, input entity.AddEmailInput) (*entity.EmailResponse, error) {
	userID, err := strconv.ParseInt(input.UserID, 10, 64)
	if err != nil || userID == 0 {
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

	emailID, err := m.service.AddEmail(ctx, userID, address.Address)
	if err != nil {
		if er := errors.Cause(err); er == iface.ErrAlreadyExists {
			return nil, &gqlerror.Error{
				Message: er.Error(),
				Extensions: map[string]interface{}{
					"code": er.(*errors.Error).Args["code"].(string),
				},
			}
		}

		log.Log(errors.New("fail to add email").SetParent(err))
		return nil, fmt.Errorf("service failed")
	}

	return &entity.EmailResponse{Email: &entity.Email{ID: strconv.FormatInt(emailID, 10)}}, nil
}

// AuthUser returns a JWT token
func (m *Mutation) AuthUser(ctx context.Context, input entity.AuthUserInput) (*entity.AuthUserResponse, error) {
	user, token, err := m.service.AuthUser(ctx, input.Email, input.Password)
	if err != nil {
		log.Log(errors.New("fail to authenticate user").SetParent(err))
		return nil, fmt.Errorf("service failed")
	}

	return &entity.AuthUserResponse{
		Token: token,
		User:  &entity.User{ID: strconv.FormatInt(user.ID, 10)},
	}, nil
}
