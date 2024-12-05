package keycloak

import (
	"errors"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

var (
	ErrInvalidGrant              = errors.New("invalid grant")
	ErrAccessTokenNotIsNotActive = errors.New("access token is not active")
)

func FromGoCloakError(clientErr error, message string) (err error) {
	defer func() {
		if err != nil && message != "" {
			err = fmt.Errorf("%s: %w", message, err)
		}
	}()

	if clientErr == nil {
		return nil
	}

	var apiErr *gocloak.APIError

	if errors.As(clientErr, &apiErr) {
		switch {
		case apiErr.Type == gocloak.APIErrTypeInvalidGrant:
			return fmt.Errorf(" %w: %s", ErrInvalidGrant, apiErr.Type)
		default:
			return clientErr
		}
	}

	return clientErr
}
