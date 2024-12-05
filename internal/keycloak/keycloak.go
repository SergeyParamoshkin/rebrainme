package keycloak

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Client struct {
	config *Config
	client *gocloak.GoCloak
}

func New(params Params) (*Client, error) {
	keycloakClient := gocloak.NewClient(params.Config.BaseURL)

	return &Client{
		config: params.Config,
		client: keycloakClient,
	}, nil
}

func (c *Client) ClientID() string {
	return c.config.ClientID
}

func (c *Client) ClientRolesFromClaims(claims jwt.MapClaims) []string {
	resourceAccess, _ := claims["resource_access"].(map[string]interface{})
	if resourceAccess == nil {
		return nil
	}

	clientAccess, _ := resourceAccess[c.config.ClientID].(map[string]interface{})
	if clientAccess == nil {
		return nil
	}

	oauth2Roles, _ := clientAccess["roles"].([]interface{})
	if oauth2Roles == nil {
		return nil
	}

	roles := make([]string, len(oauth2Roles))

	for i, r := range oauth2Roles {
		role, _ := r.(string)

		roles[i] = role
	}

	return roles
}

func (c *Client) LoginURL(redirectURL string) string {
	if redirectURL == "" {
		redirectURL = c.config.DefaultLoginRedirectURL
	}

	return fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/auth?"+
			"client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=%s",
		c.config.BaseURL,
		c.config.Realm,
		c.config.ClientID,
		c.config.RedirectURI,
		redirectURL,
		"offline_access", // scopes
	)
}

func (c *Client) ExchangeAuthorizationCode(ctx context.Context, code string) (*gocloak.JWT, error) {
	grantType := "authorization_code"

	return c.client.GetToken(ctx, c.config.Realm, gocloak.TokenOptions{
		GrantType:    &grantType,
		Code:         &code,
		ClientID:     &c.config.ClientID,
		ClientSecret: &c.config.ClientSecret,
		RedirectURI:  &c.config.RedirectURI,
	})
}

func (c *Client) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	jwt, err := c.client.Login(
		ctx, c.config.ClientID, c.config.ClientSecret, c.config.Realm, username, password,
	)

	return jwt, FromGoCloakError(err, "login")
}

func (c *Client) Logout(ctx context.Context, refreshToken string) error {
	err := c.client.Logout(
		ctx, c.config.ClientID, c.config.ClientSecret, c.config.Realm, refreshToken,
	)

	return FromGoCloakError(err, "logout")
}

func (c *Client) Refresh(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	jwt, err := c.client.RefreshToken(
		ctx, refreshToken, c.config.ClientID, c.config.ClientSecret, c.config.Realm,
	)

	return jwt, FromGoCloakError(err, "refresh")
}

func (c *Client) Validate(ctx context.Context, accessToken string) error {
	tokenResult, err := c.client.RetrospectToken(
		ctx, accessToken, c.config.ClientID, c.config.ClientSecret, c.config.Realm,
	)
	if err != nil {
		return FromGoCloakError(err, "validate")
	}

	if tokenResult.Active == nil || !*tokenResult.Active {
		return ErrAccessTokenNotIsNotActive
	}

	return nil
}

func (c *Client) ValidateAndRefresh(
	ctx context.Context, tokenPair TokenPair,
) (newTokenPair TokenPair, isNew bool, err error) {
	err = c.Validate(ctx, tokenPair.AccessToken)
	if err == nil {
		return tokenPair, false, nil
	}

	if !errors.Is(err, ErrAccessTokenNotIsNotActive) {
		return TokenPair{}, false, err
	}

	jwt, err := c.Refresh(ctx, tokenPair.RefreshToken)
	if err != nil {
		return TokenPair{}, false, err
	}

	return TokenPair{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	}, true, nil
}

func (c *Client) DecodeAccessToken(
	ctx context.Context, accessToken string,
) (*jwt.Token, jwt.MapClaims, error) {
	token, claims, err := c.client.DecodeAccessToken(ctx, accessToken, c.config.Realm)
	if err != nil {
		return nil, nil, err
	}

	return token, *claims, nil
}

func (c *Client) DecodeAccessTokenUnverified(
	accessToken string,
) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, _, err := jwt.NewParser().ParseUnverified(accessToken, claims)
	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}

func (c *Client) DecodeRealmRolesFromClaims(claims jwt.MapClaims, prefix string) []string {
	var roles []string

	if ra, ok := claims["realm_access"].(map[string]interface{}); ok {
		if raw, ok := ra["roles"].([]interface{}); ok {
			for _, v := range raw {
				r, ok := v.(string)
				if !ok || !strings.HasPrefix(r, prefix) {
					continue
				}

				roles = append(roles, r)
			}
		}
	}

	return roles
}

func (c *Client) DecodeClientFromClaims(claims jwt.MapClaims, prefix string) []string {
	var roles []string

	if ra, ok := claims["resource_access"].(map[string]interface{}); ok {
		if ca, ok := ra[c.ClientID()].(map[string]interface{}); ok {
			if raw, ok := ca["roles"].([]interface{}); ok {
				for _, v := range raw {
					r, ok := v.(string)
					if !ok || !strings.HasPrefix(r, prefix) {
						continue
					}

					roles = append(roles, r)
				}
			}
		}
	}

	return roles
}

func (c *Client) DecodeAllRolesFromClaims(claims jwt.MapClaims, prefix string) []string {
	var roles []string

	roles = append(roles, c.DecodeRealmRolesFromClaims(claims, prefix)...)
	roles = append(roles, c.DecodeClientFromClaims(claims, prefix)...)

	return roles
}
