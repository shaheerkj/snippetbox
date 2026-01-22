package models

import "errors"

// ErrNoRecord is returned when a database query finds no matching record
// This is used to distinguish "not found" from other database errors
var ErrNoRecord = errors.New("models: Record not found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateEmail = errors.New("models: duplicate email")
