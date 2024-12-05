package auth

import (
	"context"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
)

type CtxKey uint8

const (
	CtxKeyUser CtxKey = iota
)

func UserFromContext(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(CtxKeyUser).(domain.User)

	return user, ok
}
