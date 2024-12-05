package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpresp"
	fwerr "github.com/SergeyParamoshkin/alerts/internal/err"
	"github.com/SergeyParamoshkin/alerts/internal/keycloak"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	CookieOauth2AccessToken  = "accessToken"
	CookieOauth2RefreshToken = "refreshToken"

	Oauth2BackoffRetries         = 4
	Oauth2BackoffInitialInterval = 20 * time.Millisecond
	Oauth2BackoffMaxInterval     = 200 * time.Millisecond

	defaultOauth2CookieDomain = "localhost:8080"
	defaultOauth2CookiePath   = "/"
	defaultOauth2CookieSecure = false
)

type Oauth2CookieConfig struct {
	Domain string `yaml:"domain"`
	Path   string `yaml:"path"`
	Secure bool   `yaml:"secure"`
}

type OAuth2Config struct {
	Cookie Oauth2CookieConfig `yaml:"cookie"`
}

func NewDefaultOAuth2Config() OAuth2Config {
	return OAuth2Config{
		Cookie: Oauth2CookieConfig{
			Domain: defaultOauth2CookieDomain,
			Path:   defaultOauth2CookiePath,
			Secure: defaultOauth2CookieSecure,
		},
	}
}

type OAuth2Controller struct {
	logger *zap.Logger
	config *OAuth2Config

	keycloakClient *keycloak.Client
}

func NewOAuth2Controller(
	logger *zap.Logger,
	config *OAuth2Config,
	keycloakClient *keycloak.Client,
) *OAuth2Controller {
	return &OAuth2Controller{
		logger:         logger,
		config:         config,
		keycloakClient: keycloakClient,
	}
}

func (c *OAuth2Controller) backoff() backoff.BackOff { //nolint:ireturn // OK
	authRetryBackoff := backoff.NewExponentialBackOff()
	authRetryBackoff.InitialInterval = Oauth2BackoffInitialInterval
	authRetryBackoff.MaxInterval = Oauth2BackoffMaxInterval

	return backoff.WithMaxRetries(authRetryBackoff, Oauth2BackoffRetries)
}

func (c *OAuth2Controller) AccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieOauth2AccessToken)
	if err != nil {
		return "", fmt.Errorf("%w: access token cookie", ErrTokenNotFound)
	}

	return cookie.Value, nil
}

func (c *OAuth2Controller) RefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieOauth2RefreshToken)
	if err != nil {
		return "", fmt.Errorf("%w: refresh token cookie", ErrTokenNotFound)
	}

	return cookie.Value, nil
}

func (c *OAuth2Controller) GenRefreshTokenCookie(refreshToken string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     CookieOauth2RefreshToken,
		Value:    refreshToken,
		Domain:   c.config.Cookie.Domain,
		Path:     c.config.Cookie.Path,
		HttpOnly: true,
		Secure:   c.config.Cookie.Secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}

func (c *OAuth2Controller) GenAccessTokenCookie(accessToken string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     CookieOauth2AccessToken,
		Value:    accessToken,
		Domain:   c.config.Cookie.Domain,
		Path:     c.config.Cookie.Path,
		HttpOnly: true,
		Secure:   c.config.Cookie.Secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}

func (c *OAuth2Controller) SetAccessTokenCookie(
	w http.ResponseWriter, accessToken string, maxAge int,
) {
	http.SetCookie(w, c.GenAccessTokenCookie(accessToken, maxAge))
}

func (c *OAuth2Controller) SetRefreshTokenCookie(
	w http.ResponseWriter, refreshToken string, maxAge int,
) {
	http.SetCookie(w, c.GenRefreshTokenCookie(refreshToken, maxAge))
}

func (c *OAuth2Controller) GenDeleteCookies() []*http.Cookie {
	return []*http.Cookie{
		c.GenAccessTokenCookie("", -1),
		c.GenRefreshTokenCookie("", -1),
	}
}

func (c *OAuth2Controller) Logout(r *http.Request) ([]*http.Cookie, error) {
	// ctx := r.Context()

	// if user, ok := UserFromContext(ctx); ok {
	// 	c.eventTrigger.AuthLogout(ctx, user.Username)
	// }

	// refreshToken, err := c.RefreshToken(r)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get refresh token: %w", err)
	// }

	// err = c.keycloakClient.Logout(ctx, refreshToken)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to logout: %w", err)
	// }

	return c.GenDeleteCookies(), nil
}

func (c *OAuth2Controller) Authenticate(
	r *http.Request,
) (user domain.User, cookies []*http.Cookie, err error) {
	accessToken, err := c.AccessToken(r)
	if err != nil {
		return domain.User{}, nil, err
	}

	refreshToken, err := c.RefreshToken(r)
	if err != nil {
		return domain.User{}, nil, err
	}

	initialTokenPair := keycloak.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx := r.Context()

	var (
		newTokenPair keycloak.TokenPair
		retry        int
	)

	err = backoff.Retry(func() error {
		defer func() { retry++ }()

		newTokenPair, _, err = c.keycloakClient.ValidateAndRefresh(ctx, initialTokenPair)
		if err != nil {
			return fmt.Errorf("validate and refresh failed: retry %d: %w", retry, err)
		}

		return nil
	}, c.backoff())
	if err != nil {
		return domain.User{}, nil, err
	}

	user, err = c.userFromAccessToken(newTokenPair.AccessToken)
	if err != nil {
		return domain.User{}, nil, fmt.Errorf("failed to construct user from claims: %w", err)
	}

	if initialTokenPair.AccessToken != newTokenPair.AccessToken {
		cookies = append(cookies, c.GenAccessTokenCookie(newTokenPair.AccessToken, 0))
	}

	if initialTokenPair.RefreshToken != newTokenPair.RefreshToken {
		cookies = append(cookies, c.GenRefreshTokenCookie(newTokenPair.RefreshToken, 0))
	}

	return
}

func (c *OAuth2Controller) LoginHandler(w http.ResponseWriter, r *http.Request) {
	originalRedirectURL := r.URL.Query().Get("redirect_url")

	http.Redirect(w, r, c.keycloakClient.LoginURL(originalRedirectURL), http.StatusFound)
}

func (c *OAuth2Controller) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization code from the query parameters
	authorizationCode := r.URL.Query().Get("code")
	redirectURL := r.URL.Query().Get("state")

	// Exchange the authorization code for an access jwt
	jwt, err := c.keycloakClient.ExchangeAuthorizationCode(r.Context(), authorizationCode)
	if err != nil {
		errMsg := "Failed to exchange authorization code for token"

		var apiErr *gocloak.APIError

		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case http.StatusBadRequest:
				httpresp.Error(w, r, fwerr.Wrap(fwerr.CodeBadRequest, errMsg, apiErr))

			default:
				httpresp.Error(w, r, fwerr.Wrap(fwerr.CodeInternalError, errMsg, apiErr))
			}
		} else {
			httpresp.Error(w, r, fwerr.Wrap(fwerr.CodeInternalError, errMsg, err))
		}

		return
	}

	c.SetAccessTokenCookie(w, jwt.AccessToken, 0)
	c.SetRefreshTokenCookie(w, jwt.RefreshToken, 0)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (c *OAuth2Controller) userFromAccessToken(accessToken string) (domain.User, error) {
	_, claims, err := c.keycloakClient.DecodeAccessTokenUnverified(accessToken)
	if err != nil {
		return domain.User{}, fmt.Errorf("decode access token failed: %w", err)
	}

	sub, _ := claims["sub"].(string)
	username, _ := claims["preferred_username"].(string)

	userID, err := uuid.Parse(sub)
	if err != nil {
		return domain.User{}, err
	}

	rawRoles := c.keycloakClient.DecodeAllRolesFromClaims(claims, domain.RolePrefix)

	roles := make([]domain.Role, len(rawRoles))
	for i, role := range rawRoles {
		roles[i] = domain.Role(role)
	}

	return domain.User{
		ID:       userID,
		Username: username,
		Domain:   domain.DefaultUserDomain,
		Roles:    roles,
	}, nil
}
