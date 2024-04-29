package authn

import (
	"context"
	"fmt"
	"log/slog"
)

type User interface {
	HasRole(role string) bool
}

const contextKeyAuthData = "CONTEXT_KEY_AUTH_DATA"

type authData struct {
	user User
}

func GetUser(ctx context.Context) User {
	authCtxData, ok := ctx.Value(contextKeyAuthData).(*authData)
	if !ok || authCtxData == nil {
		return nil
	}
	return authCtxData.user
}

func SetUser(ctx context.Context, user User) {
	authCtxData, ok := ctx.Value(contextKeyAuthData).(*authData)
	if !ok || authCtxData == nil {
		slog.Error(fmt.Sprintf("authn.SetUser: setter not found in context"))
		return
	}
	authCtxData.user = user
	return
}
