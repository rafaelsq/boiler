package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

// AddUser add a new user
func (s *Service) AddUser(ctx context.Context, name string) (int64, error) {
	tx, err := s.storage.Tx()
	if err != nil {
		return 0, errors.New("could not begin transaction").SetParent(err)
	}

	ID, err := s.storage.AddUser(ctx, tx, name)
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

// DeleteUser remove user by ID
func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	tx, err := s.storage.Tx()
	if err != nil {
		return errors.New("could not begin delete user transaction").SetParent(err)
	}

	err = s.storage.DeleteUser(ctx, tx, userID)
	if err != nil && err != iface.ErrNotFound {
		if er := tx.Rollback(); er != nil {
			return errors.New("could not rollback delete user").SetParent(
				errors.New(er.Error()).SetParent(err),
			)
		}

		return errors.New("could not delete user").SetParent(err)
	}

	err = s.storage.DeleteEmailsByUserID(ctx, tx, userID)
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

// FilterUsers retrive users
func (s *Service) FilterUsers(ctx context.Context, filter iface.FilterUsers) ([]*entity.User, error) {
	IDs, err := s.storage.FilterUsersID(ctx, filter)
	if err != nil {
		return nil, err
	}

	return s.storage.FetchUsers(ctx, IDs...)
}

// GetUserByID get user by ID
func (s *Service) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	us, err := s.storage.FetchUsers(ctx, userID)
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
	IDs, err := s.storage.FilterUsersID(ctx, iface.FilterUsers{Email: email})
	if err != nil {
		return nil, err
	}
	if len(IDs) != 1 {
		return nil, iface.ErrNotFound
	}

	return s.GetUserByID(ctx, IDs[0])
}
