package authorization

import "context"

type ClientReader interface {
	ClientById(ctx context.Context, clientId string) (*ClientDetails, error)
}

type ClientDetails struct {
	Id            string
	RedirectUris  []string
	ResponseTypes []string
}

type ClientAuther interface {
	Auth(ctx context.Context, clientId, clientSecret string) (bool, error)
}
