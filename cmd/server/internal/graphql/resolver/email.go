package resolver

import (
	"context"
	"strconv"

	"boiler/cmd/server/internal/graphql/entity"
	"boiler/pkg/service"
	"boiler/pkg/store"
)

// NewEmail return a new Email resolver
func NewEmail(service service.Interface) *Email {
	return &Email{
		service: service,
	}
}

// Email resolver for Email
type Email struct {
	service service.Interface
}

// User resolve User by Email
func (r *Email) User(ctx context.Context, e *entity.Email) (*entity.User, error) {
	u, err := r.service.GetUserByEmail(ctx, e.Address)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, Wrap(ctx, err, "fail to get user by email")
}

// Email resolve Email by emailID
func (r *Email) Email(ctx context.Context, rawEmailID string) (*entity.Email, error) {
	emailID, err := strconv.ParseInt(rawEmailID, 10, 64)
	if err != nil || emailID == 0 {
		return nil, service.ErrInvalidID
	}

	emails, err := r.service.FilterEmails(ctx, store.FilterEmails{EmailID: emailID})
	if err != nil {
		return nil, Wrap(ctx, err, "fail to filter emails")
	}

	if len(emails) == 0 {
		return nil, store.ErrNotFound
	}

	return entity.NewEmail(emails[0]), nil
}
