package auth

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/SergeyParamoshkin/alerts/internal/keycloak"
)

type Auth struct {
	logger *zap.Logger
	oauth2 *OAuth2Controller
}

func NewAuthController(
	logger *zap.Logger,
	oauth2 *OAuth2Controller,
) *Auth {
	return &Auth{
		logger: logger.Named("auth_controller"),
		oauth2: oauth2,
	}
}

func (c *Auth) Logout(r *http.Request) ([]*http.Cookie, error) {
	var cookies []*http.Cookie

	if c, err := c.oauth2.Logout(r); err != nil {
		if !errors.Is(err, ErrTokenNotFound) {
			return nil, err
		}
	} else {
		cookies = append(cookies, c...)
	}

	return cookies, nil
}

func (c *Auth) Authenticate(r *http.Request) (domain.User, []*http.Cookie, error) {
	user, cookies, err := c.oauth2.Authenticate(r)

	switch {
	case err == nil: // oauth2 auth success
		return user, cookies, nil

	case
		errors.Is(err, ErrTokenNotFound),
		errors.Is(err, keycloak.ErrInvalidGrant): // fallback to internal auth

		if err != nil {
			return domain.User{}, nil, err
		}

		return user, nil, nil

	default:
		// fwctx.LoggerFromCtx(r.Context(), c.logger).Error("oauth2 auth error", zap.Error(err))

		return domain.User{}, nil, err
	}
}
