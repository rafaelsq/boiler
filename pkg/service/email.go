package service

import (
	"context"

	"boiler/pkg/entity"
	"boiler/pkg/iface"

	"github.com/rafaelsq/errors"
)

// AddEmail add a new email
func (s *Service) AddEmail(ctx context.Context, userID int64, address string) (int64, error) {
	tx, err := s.store.Tx()
	if err != nil {
		return 0, errors.New("could not begin transaction").SetParent(err)
	}

	ID, err := s.store.AddEmail(ctx, tx, userID, address)
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

// DeleteEmail remove an email
func (s *Service) DeleteEmail(ctx context.Context, emailID int64) error {
	tx, err := s.store.Tx()
	if err != nil {
		return errors.New("could not begin delete email transaction").SetParent(err)
	}

	err = s.store.DeleteEmail(ctx, tx, emailID)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			return errors.New("could not rollback delete email").SetParent(
				errors.New(er.Error()).SetParent(err),
			)
		}

		return errors.New("could not delete email").SetParent(err)
	}

	if err := tx.Commit(); err != nil {
		return errors.New("could not commit delete email").SetParent(err)
	}

	return nil
}

// FilterEmails retrieve emails
func (s *Service) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	return s.store.FilterEmails(ctx, filter)
}
