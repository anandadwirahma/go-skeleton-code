// Package contactusecase implements the application business logic for Contact operations.
package contactusecase

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/yourusername/go-skeleton-code/internal/domain/contact"
)

// usecase is the concrete implementation of contact.Usecase.
type usecase struct {
	repo contact.Repository
	log  *zap.Logger
}

// New returns a new contact.Usecase.
func New(repo contact.Repository, log *zap.Logger) contact.Usecase {
	return &usecase{repo: repo, log: log}
}

// Create creates a new contact after validating business rules.
func (u *usecase) Create(ctx context.Context, input contact.CreateInput) (*contact.Contact, error) {
	// Business rule: check if email is already registered.
	existing, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.Create: %w", err)
	}
	for _, c := range existing {
		if c.Email == input.Email {
			return nil, contact.ErrEmailTaken
		}
	}

	c := &contact.Contact{
		Name:    input.Name,
		Email:   input.Email,
		Message: input.Message,
	}

	if err := u.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("usecase.Create: %w", err)
	}

	u.log.Info("contact created", zap.Uint("id", c.ID), zap.String("email", c.Email))
	return c, nil
}

// GetAll returns all contacts.
func (u *usecase) GetAll(ctx context.Context) ([]contact.Contact, error) {
	contacts, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.GetAll: %w", err)
	}
	return contacts, nil
}

// GetByID returns a single contact by id.
func (u *usecase) GetByID(ctx context.Context, id uint) (*contact.Contact, error) {
	c, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, contact.ErrNotFound) {
			return nil, contact.ErrNotFound
		}
		return nil, fmt.Errorf("usecase.GetByID: %w", err)
	}
	return c, nil
}

// Update modifies the name, email, and message of an existing contact.
func (u *usecase) Update(ctx context.Context, id uint, input contact.UpdateInput) (*contact.Contact, error) {
	c, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.Name = input.Name
	c.Email = input.Email
	c.Message = input.Message

	if err := u.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("usecase.Update: %w", err)
	}

	u.log.Info("contact updated", zap.Uint("id", id))
	return c, nil
}

// Delete removes a contact by id.
func (u *usecase) Delete(ctx context.Context, id uint) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("usecase.Delete: %w", err)
	}
	u.log.Info("contact deleted", zap.Uint("id", id))
	return nil
}
