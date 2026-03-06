// Package contact defines the Contact domain entity.
// This is the core business object — it must not depend on any framework or infrastructure.
package contact

import "time"

// Contact represents a contact form submission.
type Contact struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null"  json:"name"`
	Email     string    `gorm:"type:varchar(150);not null;index" json:"email"`
	Message   string    `gorm:"type:text;not null"          json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime"              json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"              json:"updated_at"`
}

// TableName overrides the default GORM table name.
func (Contact) TableName() string {
	return "contacts"
}
