package service

import "errors"

// Common errors
var (
	ErrFailedToStartTransaction  = errors.New("failed to start transaction")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction")
)

// User errors
var (
	ErrFailedToCreateUser   = errors.New("failed to create user")
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrFailedChangePassword = errors.New("failed to change password")
	ErrUserNotFound         = errors.New("user not found")
	ErrFailedToFindUser     = errors.New("failed to find user")
)
