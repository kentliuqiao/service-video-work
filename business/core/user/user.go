package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kentliuqiao/service/foundation/logger"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	// Update(ctx context.Context, usr User) error
	// Delete(ctx context.Context, usr User) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By,
	// 	pageNumber int, rowsPerPage int) ([]User, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// Core manages the set of APIs for user access.
type Core struct {
	storer Storer
	log    *logger.Logger
}

// NewCore constructs a core for user api access.
func NewCore(log *logger.Logger, storer Storer) *Core {
	return &Core{
		storer: storer,
		log:    log,
	}
}

// Create adds a new user to the system.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return User{}, fmt.Errorf("uuid newv7: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:           id,
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		Department:   nu.Department,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := c.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
