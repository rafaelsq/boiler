package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

func (s *Service) AddUser(ctx context.Context, name string) (int, error) {
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
