package oauth2

import (
	"context"

	"github.com/jacobmiller22/gauth/pkg/oauth2"
)

type CodeRepo interface {
	Save(ctx context.Context, code, requestId string) error
}

type MemoryCodeRepo struct {
	codes map[string]string
}

func (r *MemoryCodeRepo) Save(ctx context.Context, code, requestId string) error {
	r.codes[requestId] = code
	return nil
}

type RequestRepo interface {
	Save(ctx context.Context, request oauth2.AuthorizationReq) error
}

type MemoryRequestRepo struct {
	requests map[string]oauth2.AuthorizationReq
}

func (r *MemoryRequestRepo) Save(ctx context.Context, request oauth2.AuthorizationReq) error {
	r.requests[request.Id.String()] = request
	return nil
}

type TokenRepo interface {
	Save(ctx context.Context, token *oauth2.TokenResponse) error
}

type MemoryTokenRepo struct {
	tokens map[string]oauth2.TokenResponse
}

func (r *MemoryTokenRepo) Save(ctx context.Context, token oauth2.TokenResponse) error {
	r.tokens[token.AccessToken] = token
	return nil
}
