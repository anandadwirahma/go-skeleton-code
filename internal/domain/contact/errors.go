// Package contact defines domain-level sentinel errors.
// Placing errors here keeps them accessible to both the repository and usecase layers.
package contact

import "errors"

// ErrNotFound is returned when a requested Contact does not exist.
var ErrNotFound = errors.New("contact: not found")

// ErrEmailTaken is returned when a contact with the same email already exists.
var ErrEmailTaken = errors.New("contact: email already registered")
