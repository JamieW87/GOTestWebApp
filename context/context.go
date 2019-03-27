package context

import (
	"context"
	"webapp/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

//WithUser accepts an existing context and a user, and returns a new context with that user set as the value
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

//User looks up a user from a given context
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
