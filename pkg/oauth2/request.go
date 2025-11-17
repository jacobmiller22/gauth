package oauth2

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jacobmiller22/gauth/pkg/clog"
)

const (
	ResponseTypeCode  string = "code"
	ResponseTypeToken string = "token"
)

var (
	ErrInvalidClientId         error = fmt.Errorf("access_denied")
	ErrInvalidRedirectUri      error = fmt.Errorf("access_denied")
	ErrInvalidResponseType     error = fmt.Errorf("invalid_request")
	ErrUnsupportedResponseType error = fmt.Errorf("unsupported_response_type")
	ErrUnauthorizationClient   error = fmt.Errorf("unauthorized_client")
	ErrInvalidScope            error = fmt.Errorf("invalid_scope")
)

type AuthReqReader interface {
	GetAuthReq() (*AuthorizationReq, error)
}

// The Authorization Request. Construct with NewAuthorizationRequest
type AuthorizationReq struct {
	Id           uuid.UUID
	ClientId     string
	ResponseType string
	RedirectUri  string
	Scope        string
	State        string

	CodeChallengeMethod string
	CodeChallenge       string

	IssuedAt  time.Time
	ExpiresIn time.Duration
}

func NewAuthorizationRequest(clientId, responseType, redirectUri, scope, state, codeChallengeMethod, codeChallenge string) (*AuthorizationReq, error) {
	codeChallengeMethod = strings.ToLower(codeChallengeMethod)
	switch codeChallengeMethod {
	case CodeChallengePlain, CodeChallengeS256:
		break
	case "":
		if codeChallenge != "" {
			return nil, ErrPkceMethod
		}
	default:
		return nil, ErrPkceMethod
	}

	return &AuthorizationReq{
		Id:                  uuid.New(),
		ClientId:            clientId,
		ResponseType:        responseType,
		RedirectUri:         redirectUri,
		Scope:               scope,
		State:               state,
		CodeChallengeMethod: codeChallengeMethod,
		CodeChallenge:       codeChallenge,
		IssuedAt:            time.Now(),
		ExpiresIn:           time.Minute * 10,
	}, nil
}

func (r *AuthorizationReq) Expired(nowFunc func() time.Time) bool {
	return nowFunc().Before(r.IssuedAt.Add(r.ExpiresIn))
}

func (r *AuthorizationReq) Verify(ctx context.Context, clients ClientReader) error {
	l := clog.FromContext(ctx)

	l.DebugContext(ctx, "VERIFY_AUTHORIZATION_REQUEST")

	// Fetch client for id

	client, err := clients.ClientById(ctx, r.ClientId)
	if err != nil {
		l.ErrorContext(ctx, "client-by-id-error", "err", err)
		return ErrInvalidClientId
	}

	if err := r.validateClientId(client); err != nil {
		return err
	}

	if err := r.validateRedirectUri(client.RedirectUris); err != nil {
		return err
	}

	if err := r.validateResponseType(client.ResponseTypes); err != nil {
		return err
	}

	return nil
}

// Validate the redirect uri of the authorization request according to the client
func (r *AuthorizationReq) validateClientId(client *ClientDetails) error {
	if client == nil || client.Id == "" {
		return ErrInvalidClientId
	}
	return nil
}

// Validate the redirect uri of the authorization request according to the client
func (r *AuthorizationReq) validateRedirectUri(clientRedirectUris []string) error {
	// If the request contains a redirect_uri parameter, the server must confirm it is a valid redirect URL for this application.
	if r.RedirectUri != "" {
		if !slices.Contains(clientRedirectUris, r.RedirectUri) {
			return ErrInvalidRedirectUri
		}
		return nil // Provided redirect uri is valid
	}

	// If there is no redirect_uri parameter in the request, and more than one URL was registered, this is an error
	if len(clientRedirectUris) != 1 {
		return ErrInvalidRedirectUri
	}

	// Otherwise, the server uses the redirect URL that was previously registered.
	r.RedirectUri = clientRedirectUris[0]
	return nil
}

// Validate the response type of the authorization request according to the client
func (r *AuthorizationReq) validateResponseType(clientResponseTypes []string) error {

	if r.ResponseType == "" {
		return ErrInvalidResponseType
	}

	if !slices.Contains(clientResponseTypes, r.ResponseType) {
		return ErrUnsupportedResponseType
	}

	return nil
}
