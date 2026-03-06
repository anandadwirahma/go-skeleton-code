// Package contactrepo provides the PostgreSQL implementation of the contact.Repository interface.
package contactrepo

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/yourusername/go-skeleton-code/internal/domain/contact"
)

// postgresRepository is the GORM-backed implementation of contact.Repository.
type postgresRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

// New returns a new contact.Repository backed by PostgreSQL via GORM.
func New(db *gorm.DB, log *zap.Logger) contact.Repository {
	return &postgresRepository{db: db, log: log}
}

// Create inserts a new Contact record.
func (r *postgresRepository) Create(ctx context.Context, c *contact.Contact) error {
	if err := r.db.WithContext(ctx).Create(c).Error; err != nil {
		r.log.Error("contactrepo: Create failed", zap.Error(err))
		return fmt.Errorf("contactrepo: Create: %w", err)
	}
	return nil
}

// FindAll retrieves all contacts ordered by most recent first.
func (r *postgresRepository) FindAll(ctx context.Context) ([]contact.Contact, error) {
	var contacts []contact.Contact
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&contacts).Error; err != nil {
		r.log.Error("contactrepo: FindAll failed", zap.Error(err))
		return nil, fmt.Errorf("contactrepo: FindAll: %w", err)
	}
	return contacts, nil
}

// FindByID retrieves a single Contact by primary key.
// Returns contact.ErrNotFound when the record does not exist.
func (r *postgresRepository) FindByID(ctx context.Context, id uint) (*contact.Contact, error) {
	var c contact.Contact
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contact.ErrNotFound
		}
		r.log.Error("contactrepo: FindByID failed", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("contactrepo: FindByID: %w", err)
	}
	return &c, nil
}

// Update saves changes to an existing Contact.
func (r *postgresRepository) Update(ctx context.Context, c *contact.Contact) error {
	if err := r.db.WithContext(ctx).Save(c).Error; err != nil {
		r.log.Error("contactrepo: Update failed", zap.Uint("id", c.ID), zap.Error(err))
		return fmt.Errorf("contactrepo: Update: %w", err)
	}
	return nil
}

// Delete removes a Contact by id.
// Returns contact.ErrNotFound when no row is deleted.
func (r *postgresRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&contact.Contact{}, id)
	if result.Error != nil {
		r.log.Error("contactrepo: Delete failed", zap.Uint("id", id), zap.Error(result.Error))
		return fmt.Errorf("contactrepo: Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return contact.ErrNotFound
	}
	return nil
}
