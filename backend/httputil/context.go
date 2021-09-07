package httputil

import (
	"context"

	"github.com/is0405/model"
	"github.com/pkg/errors"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func SetClaimsToContext(ctx context.Context, c *model.Claims) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, c)
}

func GetClaimsFromContext(ctx context.Context) (*model.Claims, error) {
	v := ctx.Value(ClaimsContextKey)
	claims, ok := v.(*model.Claims)
	if !ok {
		return nil, errors.New("user not found from context value")
	}
	return claims, nil
}
