package service

import (
	"context"
	"strconv"
	"time"

	"boiler/pkg/entity"
	"boiler/pkg/iface"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rafaelsq/errors"
	"golang.org/x/crypto/bcrypt"
)

// AddUser add a new user
func (s *Service) AddUser(ctx context.Context, name, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return 0, errors.New("could not generate password").SetParent(err)
	}

	tx, err := s.store.Tx()
	if err != nil {
		return 0, errors.New("could not begin transaction").SetParent(err)
	}

	ID, err := s.store.AddUser(ctx, tx, name, string(hash))
	if err != nil {
		if er := tx.Rollback(); er != nil {
			return 0, errors.New("could not add user").SetParent(
				errors.New(er.Error()).SetParent(err),
			)
		}

		return 0, errors.New("could not add user").SetParent(err)
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.New("could not add user").SetParent(err)
	}

	return ID, nil
}

// AuthUser returns a JWT token from users credentials
func (s *Service) AuthUser(ctx context.Context, email, password string) (*entity.User, string, error) {
	var token string

	IDs, err := s.store.FilterUsersID(ctx, iface.FilterUsers{Email: email})
	if err != nil {
		return nil, token, err
	}
	if len(IDs) != 1 {
		return nil, token, iface.ErrNotFound
	}

	user, err := s.GetUserByID(ctx, IDs[0])
	if err != nil {
		return nil, token, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, token, iface.ErrInvalidPassword
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
		return nil, token, err
	}

	token = string(raw)

	return user, token, nil
}

// EnqueueDeleteUser enqueue user to be deleted
func (s *Service) EnqueueDeleteUser(ctx context.Context, userID int64) error {
	_, err := s.enqueue.Enqueue(iface.DeleteUser, map[string]interface{}{"id": userID})
	return err
}

// DeleteUser remove user by ID
func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	tx, err := s.store.Tx()
	if err != nil {
		return errors.New("could not begin delete user transaction").SetParent(err)
	}

	err = s.store.DeleteUser(ctx, tx, userID)
	if err != nil && err != iface.ErrNotFound {
		if er := tx.Rollback(); er != nil {
			return errors.New("could not rollback delete user").SetParent(
				errors.New(er.Error()).SetParent(err),
			)
		}

		return errors.New("could not delete user").SetParent(err)
	}

	err = s.store.DeleteEmailsByUserID(ctx, tx, userID)
	if err != nil && err != iface.ErrNotFound {
		if er := tx.Rollback(); er != nil {
			return errors.New("could not rollback delete emails by user ID").SetParent(
				errors.New(er.Error()).SetParent(err),
			)
		}

		return errors.New("could not delete user emails").SetParent(err)
	}

	if err := tx.Commit(); err != nil {
		return errors.New("could not commit delete user").SetParent(err)
	}

	return nil
}

// FilterUsers retrieve users
func (s *Service) FilterUsers(ctx context.Context, filter iface.FilterUsers) ([]*entity.User, error) {
	IDs, err := s.store.FilterUsersID(ctx, filter)
	if err != nil {
		return nil, err
	}

	return s.store.FetchUsers(ctx, IDs...)
}

// GetUserByID get user by ID
func (s *Service) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	us, err := s.store.FetchUsers(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(us) != 1 {
		return nil, iface.ErrNotFound
	}
	return us[0], nil
}

// GetUserByEmail get user by Email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	IDs, err := s.store.FilterUsersID(ctx, iface.FilterUsers{Email: email})
	if err != nil {
		return nil, err
	}
	if len(IDs) != 1 {
		return nil, iface.ErrNotFound
	}

	return s.GetUserByID(ctx, IDs[0])
}
