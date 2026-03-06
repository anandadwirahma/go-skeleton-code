// Package dto defines request and response Data Transfer Objects for the Contact HTTP API.
// DTOs live in the delivery layer and must NOT be imported by the domain or usecase layers.
package dto

import "github.com/yourusername/go-skeleton-code/internal/domain/contact"

// CreateContactRequest is the expected JSON body for POST /contacts.
type CreateContactRequest struct {
	Name    string `json:"name"    binding:"required,min=2,max=100"`
	Email   string `json:"email"   binding:"required,email"`
	Message string `json:"message" binding:"required,min=10"`
}

// UpdateContactRequest is the expected JSON body for PUT /contacts/:id.
type UpdateContactRequest struct {
	Name    string `json:"name"    binding:"required,min=2,max=100"`
	Email   string `json:"email"   binding:"required,email"`
	Message string `json:"message" binding:"required,min=10"`
}

// ContactResponse is the JSON response returned to the client.
type ContactResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// ToResponse converts a domain Contact entity to a ContactResponse DTO.
func ToResponse(c *contact.Contact) ContactResponse {
	return ContactResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Message:   c.Message,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToResponseList converts a slice of domain contacts to response DTOs.
func ToResponseList(contacts []contact.Contact) []ContactResponse {
	result := make([]ContactResponse, len(contacts))
	for i := range contacts {
		result[i] = ToResponse(&contacts[i])
	}
	return result
}

// SuccessResponse is the standard success envelope returned by all endpoints.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ErrorResponse is the standard error envelope returned on failures.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
