package exceptions

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")

	ErrUserExists = errors.New("user already exists")

	ErrDiscordInteraction = errors.New("discord interaction error")

	ErrDatabaseOp = errors.New("database operation failed")

	ErrInvalidInput = errors.New("invalid input")
)

type UserNotFoundError struct {
	UserID int
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.UserID)
}

func (e UserNotFoundError) Is(target error) bool {
	return target == ErrUserNotFound
}

type UserExistsError struct {
	UserID   int
	Username string
}

func (e UserExistsError) Error() string {
	return fmt.Sprintf("user %s (ID: %d) is already registered", e.Username, e.UserID)
}

func (e UserExistsError) Is(target error) bool {
	return target == ErrUserExists
}

// DatabaseError wraps database-related errors
type DatabaseError struct {
	Operation string
	Err       error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

func (e DatabaseError) Is(target error) bool {
	return target == ErrDatabaseOp
}

// DiscordError represents errors when interacting with Discord API
type DiscordError struct {
	Action string
	Err    error
}

func (e DiscordError) Error() string {
	return fmt.Sprintf("discord error during %s: %v", e.Action, e.Err)
}

func (e DiscordError) Unwrap() error {
	return e.Err
}

func (e DiscordError) Is(target error) bool {
	return target == ErrDiscordInteraction
}

// ValidationError represents input validation errors
type ValidationError struct {
	Field string
	Value interface{}
	Msg   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid %s (%v): %s", e.Field, e.Value, e.Msg)
}

func (e ValidationError) Is(target error) bool {
	return target == ErrInvalidInput
}
