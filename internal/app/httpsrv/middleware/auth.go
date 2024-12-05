package middleware

import (
	"context"
	"net/http"

	"github.com/SergeyParamoshkin/alerts/internal/app/auth"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpresp"
	fwerr "github.com/SergeyParamoshkin/alerts/internal/err"
	"go.uber.org/zap"
)

func Auth(
	logger *zap.Logger, a *auth.Auth,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, cookies, err := a.Authenticate(r)
			if err != nil {
				httpresp.Render(w, r, httpresp.NewErrResponse(
					fwerr.Wrap(fwerr.CodeAuthError, "authorization failed", err)),
				)

				return
			}

			for _, cookie := range cookies {
				http.SetCookie(w, cookie)
			}

			ctx := context.WithValue(r.Context(), auth.CtxKeyUser, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
