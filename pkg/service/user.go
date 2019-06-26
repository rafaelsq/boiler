package service

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func (s *Service) AddUser(ctx context.Context, name string) (int, error) {
	tx, err := s.storage.Tx()
	if err != nil {
		return 0, multierr.Combine(err, fmt.Errorf("could not begin transaction"))
	}

	ID, err := s.storage.AddUser(ctx, tx, name)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			return 0, multierr.Combine(err, er, fmt.Errorf("could not add user"))
		}

		return 0, multierr.Combine(err, fmt.Errorf("could not add user"))
	}

	if err := tx.Commit(); err != nil {
		return 0, multierr.Combine(err, fmt.Errorf("could not add user"))
	}

	return ID, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID int) error {
	return s.storage.DeleteUser(ctx, userID)
}

func (s *Service) FilterUsers(ctx context.Context, filter iface.FilterUsers) ([]*entity.User, error) {
	return s.storage.FilterUsers(ctx, filter)
}

func (s *Service) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	us, err := s.storage.FilterUsers(ctx, iface.FilterUsers{UserID: userID})
	if err != nil {
		return nil, err
	}
	if len(us) != 1 {
		return nil, iface.ErrNotFound
	}
	return us[0], nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	us, err := s.storage.FilterUsers(ctx, iface.FilterUsers{Email: email})
	if err != nil {
		return nil, err
	}
	if len(us) != 1 {
		return nil, iface.ErrNotFound
	}
	return us[0], nil
}
