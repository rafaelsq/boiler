package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"boiler/pkg/entity"
	"boiler/pkg/errors"
	"boiler/pkg/store"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AddUser add a new user
func (s *Service) AddUser(ctx context.Context, user *entity.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not generate password; %w", err)
	}

	tx, err := s.store.Tx()
	if err != nil {
		return fmt.Errorf("could not begin transaction; %w", err)
	}

	user.Password = string(hash)

	err = s.store.AddUser(ctx, tx, user)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			err = fmt.Errorf("%s; %w", er, err)
		}

		return fmt.Errorf("could not add user; %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not add user; %w", err)
	}

	return nil
}

// AuthUser returns a JWT token from users credentials
func (s *Service) AuthUser(ctx context.Context, email, password string, user *entity.User, token *string) error {

	var IDs []int64
	err := s.store.FilterUsersID(ctx, store.FilterUsers{Email: email, Limit: FilterUsersDefaultLimit}, &IDs)
	if err != nil {
		return err
	}
	if len(IDs) != 1 {
		return errors.ErrNotFound
	}

	err = s.GetUserByID(ctx, IDs[0], user)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return errors.ErrInvalidPassword
	}

	t := jwt.New()

	// https://tools.ietf.org/html/rfc7519#page-9
	_ = t.Set(jwt.SubjectKey, strconv.FormatInt(user.ID, 10))
	_ = t.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(s.config.JWT.ExpireIn).Unix())
	_ = t.Set(jwt.AudienceKey, "auth")
	_ = t.Set(jwt.IssuerKey, s.config.JWT.Issuer)

	raw, err := jwt.Sign(t, jwa.RS256, s.config.JWT.PrivateKey)
	if err != nil {
		return err
	}

	*token = string(raw)

	return nil
}

// EnqueueDeleteUser enqueue user to be deleted
func (s *Service) EnqueueDeleteUser(ctx context.Context, userID int64) error {
	_, err := s.enqueuer.Enqueue(DeleteUser, map[string]interface{}{"id": userID})
	return err
}

// DeleteUser remove user by ID
func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	tx, err := s.store.Tx()
	if err != nil {
		return fmt.Errorf("could not begin delete user transaction; %w", err)
	}

	err = s.store.DeleteUser(ctx, tx, userID)
	if err != nil && err != errors.ErrNotFound {
		if er := tx.Rollback(); er != nil {
			err = fmt.Errorf("%s; %w", er, err)
		}

		return fmt.Errorf("could not delete user; %w", err)
	}

	err = s.store.DeleteEmailsByUserID(ctx, tx, userID)
	if err != nil && err != errors.ErrNotFound {
		if er := tx.Rollback(); er != nil {
			err = fmt.Errorf("%s; %w", er, err)
		}

		return fmt.Errorf("could not delete user emails; %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit delete user; %w", err)
	}

	return nil
}

// FilterUsers retrieve users
func (s *Service) FilterUsers(ctx context.Context, filter store.FilterUsers, users *[]entity.User) error {

	if filter.Limit == 0 {
		filter.Limit = FilterUsersDefaultLimit
	}

	var IDs []int64
	err := s.store.FilterUsersID(ctx, filter, &IDs)
	if err != nil {
		return err
	}

	return s.store.FetchUsers(ctx, IDs, users)
}

// GetUserByID get user by ID
func (s *Service) GetUserByID(ctx context.Context, userID int64, user *entity.User) error {
	var users []entity.User
	err := s.store.FetchUsers(ctx, []int64{userID}, &users)
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return errors.ErrNotFound
	}
	*user = users[0]
	return nil
}

// GetUserByEmail get user by Email
func (s *Service) GetUserByEmail(ctx context.Context, email string, user *entity.User) error {

	filter := store.FilterUsers{Email: email, Limit: FilterUsersDefaultLimit}

	var IDs []int64
	err := s.store.FilterUsersID(ctx, filter, &IDs)
	if err != nil {
		return err
	}
	if len(IDs) != 1 {
		return errors.ErrNotFound
	}

	return s.GetUserByID(ctx, IDs[0], user)
}
