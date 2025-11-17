package oauth2

import (
	"fmt"
	"slices"
	"strings"
)

var ErrInvalidClientId = fmt.Errorf("Invalid client_id")       // TODO: Get correct err code
var ErrInvalidRedirectUri = fmt.Errorf("Invalid redirect_uri") // TODO: Get correct err code
var ErrInvalidResponseType = fmt.Errorf("invalid_request")
var ErrUnauthorizationClient = fmt.Errorf("unauthorized_client")
var ErrInvalidScope = fmt.Errorf("invalid_scope")

type ClientDetails struct { // TODO: Rename this

	ClientId             string
	ValidRedirectUris    []string
	AllowedResponseTypes []string // TODO: Consider bitmap
}

type AuthorizationReq struct {
	ResponseType string
	ClientId     string
	RedirectUri  *string
	ClientSecret string
	Scope        string
	State        string

	UsePkce             bool // Must be remembered by the authorization server between issuing the authorization code and issuing the access token
	CodeChallenge       string
	CodeChallengeMethod string
}

func (r *AuthorizationReq) VerifyAuthorizationRequest(getClientDetails func(clientId string) *ClientDetails) error {

	// Fetch client for id
	client := getClientDetails(r.ClientId)

	// Validate client_id
	if client == nil {
		return ErrInvalidClientId
	}

	// Validate redirect_uri
	// If the request contains a redirect_uri parameter, the server must confirm it is a valid redirect URL for this application. If there is no redirect_uri parameter in the request, and only one URL was registered, the server uses the redirect URL that was previously registered. Otherwise, if no redirect URL is in the request, and no redirect URL has been registered, this is an error.
	if r.RedirectUri == nil {
		if len(client.ValidRedirectUris) != 1 {
			return ErrInvalidRedirectUri
		}
		newRedirectUri := strings.Clone(client.ValidRedirectUris[0]) // TODO: Evaluate necessity of clone
		r.RedirectUri = &newRedirectUri

	} else if !slices.Contains(client.ValidRedirectUris, *r.RedirectUri) {
		return ErrInvalidRedirectUri
	}

	if r.ResponseType == "" {
		return ErrInvalidResponseType
	}

	if !slices.Contains(client.AllowedResponseTypes, r.ResponseType) {
		return ErrInvalidResponseType
	}

	return nil
}
