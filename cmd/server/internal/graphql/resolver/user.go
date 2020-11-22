package resolver

import (
	"context"
	"strconv"

	"boiler/cmd/server/internal/graphql/entity"
	"boiler/pkg/service"
	"boiler/pkg/store"
)

// NewUser return a new user resolver
func NewUser(srv service.Interface) *User {
	return &User{
		service: srv,
	}
}

// User is the user resolver
type User struct {
	service service.Interface
}

// User return an user by ID
func (r *User) User(ctx context.Context, rawUserID string) (*entity.User, error) {
	userID, err := strconv.ParseInt(rawUserID, 10, 64)
	if err != nil || userID == 0 {
		return nil, service.ErrInvalidID
	}

	u, err := r.service.GetUserByID(ctx, userID)
	if err == nil {
		return entity.NewUser(u), nil
	}
	return nil, Wrap(ctx, err, "fail to get user")
}

// Users return a slice of User
func (r *User) Users(ctx context.Context, limit uint) ([]*entity.User, error) {
	us, err := r.service.FilterUsers(ctx, store.FilterUsers{Limit: limit})
	if err == nil {
		users := make([]*entity.User, 0, len(us))
		for _, u := range us {
			users = append(users, entity.NewUser(u))
		}
		return users, nil
	}
	return nil, Wrap(ctx, err, "fail to filter users")
}

// Emails return a slice of email
func (r *User) Emails(ctx context.Context, u *entity.User) ([]*entity.Email, error) {
	userID, err := strconv.ParseInt(u.ID, 10, 64)
	if err != nil || userID == 0 {
		return nil, service.ErrInvalidID
	}

	es, err := r.service.FilterEmails(ctx, store.FilterEmails{UserID: userID})
	if err == nil {
		emails := make([]*entity.Email, 0, len(es))
		for _, e := range es {
			emails = append(emails, entity.NewEmail(e))
		}
		return emails, nil
	}

	return nil, Wrap(ctx, err, "fail to filter emails")
}
