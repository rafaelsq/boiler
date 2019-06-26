package service

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
)

func (s *Service) AddEmail(ctx context.Context, userID int, address string) (int, error) {
	tx, err := s.storage.Tx()
	if err != nil {
		return 0, multierr.Combine(err, fmt.Errorf("could not begin transaction"))
	}

	ID, err := s.storage.AddEmail(ctx, tx, userID, address)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			return 0, multierr.Combine(err, er, fmt.Errorf("could not add email"))
		}

		return 0, multierr.Combine(err, fmt.Errorf("could not add email"))
	}

	if err := tx.Commit(); err != nil {
		return 0, multierr.Combine(err, fmt.Errorf("could not add email"))
	}

	return ID, nil
}

func (s *Service) DeleteEmail(ctx context.Context, emailID int) error {
	return s.storage.DeleteEmail(ctx, emailID)
}

func (s *Service) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	return s.storage.FilterEmails(ctx, filter)
}
