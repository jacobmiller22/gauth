package oauth2

import "context"

type UserReader interface {
	UserById(ctx context.Context, userId string) (*UserDetails, error)
}

type UserDetails struct {
	Username string
}

type UserAuther interface {
	Auth(ctx context.Context, userId, userPassword string) (bool, error)
}
