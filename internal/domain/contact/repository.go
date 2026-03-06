// Package contact defines the repository interface for the Contact domain.
// Infrastructure implementations must satisfy this contract.
package contact

import "context"

// Repository defines the persistence contract for Contact entities.
// Any database technology (PostgreSQL, MySQL, in-memory) must implement this interface.
type Repository interface {
	// Create persists a new Contact and returns it with generated fields populated.
	Create(ctx context.Context, contact *Contact) error

	// FindAll returns all Contact records ordered by created_at DESC.
	FindAll(ctx context.Context) ([]Contact, error)

	// FindByID returns the Contact with the given id, or an error if not found.
	FindByID(ctx context.Context, id uint) (*Contact, error)

	// Update saves changes to an existing Contact.
	Update(ctx context.Context, contact *Contact) error

	// Delete removes the Contact with the given id.
	Delete(ctx context.Context, id uint) error
}
