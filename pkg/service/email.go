package service

import (
	"context"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/errors"
)

func (s *Service) AddEmail(ctx context.Context, userID int, address string) (int, error) {
	tx, err := s.storage.Tx()
	if err != nil {
		return 0, errors.New("could not begin transaction").SetParent(err)
	}

	ID, err := s.storage.AddEmail(ctx, tx, userID, address)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			return 0, errors.New("could not add email").SetParent(errors.New(er.Error()).SetParent(err))
		}

		return 0, errors.New("could not add email").SetParent(err)
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.New("could not add email").SetParent(err)
	}

	return ID, nil
}

func (s *Service) DeleteEmail(ctx context.Context, emailID int) error {
	return s.storage.DeleteEmail(ctx, emailID)
}

func (s *Service) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	return s.storage.FilterEmails(ctx, filter)
}
