package resolver

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

// NewResponse return a new Response resolver
func NewResponse(service iface.Service) *Response {
	return &Response{
		service: service,
	}
}

// Response is a response resolver
type Response struct {
	service iface.Service
}

// User return an User
func (r *Response) User(ctx context.Context, ur *entity.UserResponse) (*entity.User, error) {
	return (&User{r.service}).User(ctx, ur.User.ID)
}

// Email return an Email
func (r *Response) Email(ctx context.Context, ur *entity.EmailResponse) (*entity.Email, error) {
	return (&Email{r.service}).Email(ctx, ur.Email.ID)
}
