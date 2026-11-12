package storage

import "context"

// CreateAdminInput is a validated admin with an already hashed password.
type CreateAdminInput struct {
	Name, Email, PasswordHash, Level string
	IPAddress, UserAgent             string
}

type CreateAdminPersistence interface {
	// CreateAdmin returns the new user's id; ErrEmailAlreadyExists if taken.
	CreateAdmin(ctx context.Context, in CreateAdminInput) (string, error)
}
