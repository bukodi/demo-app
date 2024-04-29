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
	user    User
	changed bool
}

func getAuthData(ctx context.Context) *authData {
	authCtxData, ok := ctx.Value(contextKeyAuthData).(*authData)
	if !ok || authCtxData == nil {
		return nil
	}
	return authCtxData
}

func GetUser(ctx context.Context) User {
	authCtxData := getAuthData(ctx)
	if authCtxData == nil {
		return nil
	}
	return authCtxData.user
}

func SetUser(ctx context.Context, user User) {
	authCtxData := getAuthData(ctx)
	if authCtxData == nil {
		slog.Error(fmt.Sprintf("authn.SetUser: setter not found in context"))
		return
	}
	authCtxData.user = user
	authCtxData.changed = true
	return
}
