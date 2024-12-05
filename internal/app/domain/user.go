package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthType string

const (
	AuthTypeSystem   AuthType = "system"
	AuthTypeInternal AuthType = "internal"
	AuthTypeOauth2   AuthType = "oauth2"

	DefaultUserDomain = "domain1" // TODO: implement this!!!
)

type User struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AuthType     AuthType
	Username     string
	Domain       string // TODO: implement this !!!
	PasswordHash []byte

	Roles []Role
}

func NewDeviceUser(esn string) User {
	return User{
		Username: esn,
		Domain:   DefaultUserDomain,
	}
}

func NewInternalUser() User {
	return User{
		AuthType: AuthTypeSystem,
		Username: "_internal_",
		Domain:   DefaultUserDomain,
	}
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
}

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", nil
	}

	passwordBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(passwordBytes), err
}
