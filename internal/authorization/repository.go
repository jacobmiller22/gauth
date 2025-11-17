package authorization

import (
	"context"

	"github.com/jacobmiller22/gauth/pkg/authorization"
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
	Save(ctx context.Context, request authorization.AuthorizationReq) error
}

type MemoryRequestRepo struct {
	requests map[string]authorization.AuthorizationReq
}

func (r *MemoryRequestRepo) Save(ctx context.Context, request authorization.AuthorizationReq) error {
	r.requests[request.Id.String()] = request
	return nil
}

type TokenRepo interface {
	Save(ctx context.Context, token *authorization.TokenResponse) error
}

type MemoryTokenRepo struct {
	tokens map[string]authorization.TokenResponse
}

func (r *MemoryTokenRepo) Save(ctx context.Context, token authorization.TokenResponse) error {
	r.tokens[token.AccessToken] = token
	return nil
}
