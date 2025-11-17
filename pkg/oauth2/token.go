package oauth2

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jacobmiller22/gauth/pkg/pkce"
)

const (
	GrantTypeAuthorizationCode string = "authorization_code"
	GrantTypePassword          string = "password"
	GrantTypeClientCredentials string = "client_credentials"
)

var (
	ErrInvalidTokenRequest    error = errors.New("token: invalid_request")
	ErrInvalidTokenExpiration error = errors.New("token: invalid_expiration")
)

var (
	ErrInvalidRequest       error = errors.New("invalid_request")
	ErrInvalidClient        error = errors.New("invalid_client")
	ErrInvalidGrant         error = errors.New("invalid_grant")
	ErrUnauthorizedClient   error = errors.New("unauthorized_client")
	ErrUnsupportedGrantType error = errors.New("unknown_grant_type")
)

type TokenRequest struct {
	Context context.Context

	GrantType string

	Code        string
	RedirectUri string

	Username     string
	Password     string
	ClientId     string
	ClientSecret string

	Scope string

	CodeVerifier string
}

func NewTokenRequest(
	ctx context.Context,
	grantType,
	code,
	redirectUri,
	username,
	password,
	clientId,
	clientSecret,
	scope,
	codeVerifier string,
) (*TokenRequest, error) {

	r := TokenRequest{
		Context:      ctx,
		GrantType:    grantType,
		Code:         code,
		RedirectUri:  redirectUri,
		Username:     username,
		Password:     password,
		ClientId:     clientId,
		ClientSecret: clientId,
		Scope:        scope,
		CodeVerifier: codeVerifier,
	}

	if err := r.validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTokenRequest, err)
	}

	return &r, nil
}

func (r *TokenRequest) validate() error {
	if r.Context == nil {
		r.Context = context.Background()
	}
	r.GrantType = strings.ToLower(r.GrantType)
	switch r.GrantType {
	case GrantTypeAuthorizationCode:
		return r.validateGrantAuthorizationCode()
	case GrantTypePassword:
		return r.validateGrantPassword()
	case GrantTypeClientCredentials:
		return r.validateGrantClientCredentials()
	default:
		return ErrInvalidTokenRequest
	}
}

func (r *TokenRequest) validateGrantAuthorizationCode() error {
	if r.GrantType != GrantTypeAuthorizationCode {
		return fmt.Errorf("invalid grant type")
	}
	if r.ClientId == "" {
		return fmt.Errorf("missing client_id")
	}
	return nil
}

func (r *TokenRequest) validateGrantPassword() error {
	if r.GrantType != GrantTypePassword {
		return fmt.Errorf("invalid grant type")
	}
	if r.Username == "" {
		return fmt.Errorf("missing username")
	}
	if r.Password == "" {
		return fmt.Errorf("missing password")
	}
	return nil
}
func (r *TokenRequest) validateGrantClientCredentials() error {
	if r.GrantType != GrantTypeClientCredentials {
		return fmt.Errorf("invalid grant type")
	}
	if r.ClientId == "" {
		return fmt.Errorf("missing client_id")
	}
	if r.ClientSecret == "" {
		return fmt.Errorf("missing client_secret")
	}
	return nil
}

type tokenGrantFunc func(*TokenRequest) error

type AccessTokenMinter interface {
	AccessToken(r *TokenRequest) (string, error)
}

type RefreshTokenMinter interface {
	RefreshToken(r *TokenRequest) (string, error)
}

type TokenRequestProcessor struct {
	codeReader          AuthCodeReader
	clientReader        ClientReader
	clientAuther        ClientAuther
	userAuther          UserAuther
	accessTokenMinter   AccessTokenMinter
	refreshTokenMinter  RefreshTokenMinter
	tokenExpirationFunc func(*TokenRequest) time.Duration
}

func (p *TokenRequestProcessor) verify(r *TokenRequest) error {

	var grantFunc tokenGrantFunc

	switch r.GrantType {
	case GrantTypeAuthorizationCode:
		grantFunc = p.tokenGrantAuthorizationCode
	case GrantTypePassword:
		grantFunc = p.tokenGrantPassword
	case GrantTypeClientCredentials:
		grantFunc = p.tokenGrantClientCredentials
	default:
		return ErrUnsupportedGrantType
	}

	if err := grantFunc(r); err != nil {
		return err
	}

	return nil
}

// The authorization code grant is used when an application exchanges an
// authorization code for an access token.
func (p *TokenRequestProcessor) tokenGrantAuthorizationCode(r *TokenRequest) error {
	c, err := p.codeReader.AuthCode(r.Code)
	if err != nil {
		return err
	}

	if c.Expired(func() time.Time { return time.Now() }) {
		// TODO: Correct error
		return fmt.Errorf("token has expired")
	}

	if r.ClientId != c.ClientId {
		// TODO: Correct error
		return fmt.Errorf("client id mismatch")
	}

	if r.RedirectUri != c.RedirectUri {
		// TODO: Correct error
		return fmt.Errorf("redirect uri mismatch")
	}

	if c.CodeChallenge != "" {
		if err := pkce.Verify(c.CodeChallengeMethod, c.CodeChallenge, r.CodeVerifier); err != nil {
			// TODO: Correct error
			return err
		}
	}

	return nil
}

// The Password grant is used when the application exchanges the user’s username
// and password for an access token. This is exactly the thing OAuth was created
// to prevent in the first place, so you should never allow third-party apps to
// use this grant.
func (p *TokenRequestProcessor) tokenGrantPassword(r *TokenRequest) error {

	// If a client was issued a secret, client authentication must be done
	clientAuthed, err := p.clientAuther.Auth(r.Context, r.ClientId, r.ClientSecret)
	if err != nil || !clientAuthed {
		// TODO: Error handle
		return fmt.Errorf("invalid client")
	}

	userAuthed, err := p.userAuther.Auth(r.Context, r.Username, r.Password)
	if err != nil {
		// TODO: Error handle
		return err
	}

	if !userAuthed {
		// TODO: Error handle
		return fmt.Errorf("invalid")
	}

	return nil
}

// The Client Credentials grant is used when applications request an access
// token to access their own resources, not on behalf of a user.
func (p *TokenRequestProcessor) tokenGrantClientCredentials(r *TokenRequest) error {
	clientAuthed, err := p.clientAuther.Auth(r.Context, r.ClientId, r.ClientSecret)

	if err != nil || !clientAuthed {
		// TODO: Error handle
		return fmt.Errorf("invalid client")
	}

	return nil
}

func (p *TokenRequestProcessor) BearerToken(r *TokenRequest) (*TokenResponse, error) {
	p.verify(r)

	accessToken, err := p.accessTokenMinter.AccessToken(r)
	if err != nil {
		// TODO: correct error
		return nil, err
	}
	refreshToken, err := p.refreshTokenMinter.RefreshToken(r)
	if err != nil {
		// TODO: correct error
		return nil, err
	}

	expiresIn := p.tokenExpirationFunc(r)

	if expiresIn < 0 {
		return nil, ErrInvalidTokenExpiration
	}

	return &TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scope:        r.Scope,
		ExpiresIn:    expiresIn,
	}, nil
}

type TokenResponse struct {
	TokenType    string
	AccessToken  string
	ExpiresIn    time.Duration
	RefreshToken string
	Scope        string
}
