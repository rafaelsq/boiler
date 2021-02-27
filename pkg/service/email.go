package service

import (
	"context"
	"fmt"

	"boiler/pkg/entity"
	"boiler/pkg/store"
)

// AddEmail add a new email
func (s *Service) AddEmail(ctx context.Context, email *entity.Email) error {
	tx, err := s.store.Tx()
	if err != nil {
		return fmt.Errorf("could not begin transaction; %w", err)
	}

	err = s.store.AddEmail(ctx, tx, email)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			err = fmt.Errorf("%s; %w", er, err)
		}

		return fmt.Errorf("could not add email; %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not add email; %w", err)
	}

	return nil
}

// EnqueueDeleteEmail enqueue email to be deleted
func (s *Service) EnqueueDeleteEmail(ctx context.Context, emailID int64) error {
	_, err := s.enqueuer.Enqueue(DeleteEmail, map[string]interface{}{"id": emailID})
	return err
}

// DeleteEmail remove an email
func (s *Service) DeleteEmail(ctx context.Context, emailID int64) error {
	tx, err := s.store.Tx()
	if err != nil {
		return fmt.Errorf("could not begin delete email transaction; %w", err)
	}

	err = s.store.DeleteEmail(ctx, tx, emailID)
	if err != nil {
		if er := tx.Rollback(); er != nil {
			err = fmt.Errorf("%s; %w", er, err)
		}

		return fmt.Errorf("could not delete email; %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit delete email; %w", err)
	}

	return nil
}

// FilterEmails retrieve emails
func (s *Service) FilterEmails(ctx context.Context, filter store.FilterEmails, emails *[]entity.Email) error {
	return s.store.FilterEmails(ctx, filter, emails)
}
