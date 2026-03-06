// Package contact defines the use-case interface for Contact business logic.
// The delivery layer depends on this interface, not on concrete implementations.
package contact

import "context"

// Usecase defines the application business logic contract for Contact operations.
type Usecase interface {
	// Create validates and persists a new contact.
	Create(ctx context.Context, input CreateInput) (*Contact, error)

	// GetAll retrieves all contacts.
	GetAll(ctx context.Context) ([]Contact, error)

	// GetByID retrieves a single contact by its id.
	GetByID(ctx context.Context, id uint) (*Contact, error)

	// Update modifies an existing contact.
	Update(ctx context.Context, id uint, input UpdateInput) (*Contact, error)

	// Delete removes a contact by id.
	Delete(ctx context.Context, id uint) error
}

// CreateInput carries the validated data needed to create a contact.
type CreateInput struct {
	Name    string
	Email   string
	Message string
}

// UpdateInput carries the validated data for updating a contact.
type UpdateInput struct {
	Name    string
	Email   string
	Message string
}
