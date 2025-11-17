package oauth2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestVerifyAuthorizationRequestValid(t *testing.T) {

	valid_redirect_uri := "myapp.com"

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp.com"},
			AllowedResponseTypes: []string{"code"},
		}
	}

	req := AuthorizationReq{
		ResponseType: "code",
		ClientId:     "ValidClientId",
		RedirectUri:  &valid_redirect_uri,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(nil, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

// Authorization Request is missing redirect uri, application has multiple valid redirect uris
func TestVerifyAuthorizationRequestMissingRedirectUriMultipleValidUri(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com", "myapp2.com"},
			AllowedResponseTypes: []string{"code"},
		}
	}

	req := AuthorizationReq{
		ResponseType: "code",
		ClientId:     "ValidClientId",
		RedirectUri:  nil,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(ErrInvalidRedirectUri, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

// Authorization Request is missing redirect uri, application has a single valid redirect uri
func TestVerifyAuthorizationRequestMissingRedirectUriOneValidUri(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com"},
			AllowedResponseTypes: []string{"code"},
		}
	}

	req := AuthorizationReq{
		ResponseType: "code",
		ClientId:     "ValidClientId",
		RedirectUri:  nil,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(nil, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

func TestVerifyAuthorizationRequestInvalidRedirectUriOneValidUri(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com"},
			AllowedResponseTypes: []string{"code"},
		}
	}

	redirectUri := "myapp1.com"

	req := AuthorizationReq{
		ResponseType: "code",
		ClientId:     "ValidClientId",
		RedirectUri:  &redirectUri,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(nil, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

func TestVerifyAuthorizationRequestInvalidRedirectUriMultipleValidUri(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com", "myapp2.com", "myapp3.com"},
			AllowedResponseTypes: []string{"code"},
		}
	}

	redirectUri := "myapp2.com"

	req := AuthorizationReq{
		ResponseType: "code",
		ClientId:     "ValidClientId",
		RedirectUri:  &redirectUri,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(nil, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

func TestVerifyAuthorizationRequestInvalidResponseType(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com", "myapp2.com", "myapp3.com"},
			AllowedResponseTypes: []string{"code", "token"},
		}
	}

	redirectUri := "myapp2.com"

	req := AuthorizationReq{
		ResponseType: "notcodeortoken",
		ClientId:     "ValidClientId",
		RedirectUri:  &redirectUri,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(ErrInvalidResponseType, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

func TestVerifyAuthorizationRequestValidResponseTypeCode(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com", "myapp2.com", "myapp3.com"},
			AllowedResponseTypes: []string{"code", "token"},
		}
	}

	redirectUri := "myapp2.com"

	req := AuthorizationReq{
		ResponseType: "code",
		ClientId:     "ValidClientId",
		RedirectUri:  &redirectUri,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(nil, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}

func TestVerifyAuthorizationRequestValidResponseTypeToken(t *testing.T) {

	getClients := func(clientId string) *ClientDetails {
		if clientId != "ValidClientId" {
			return nil
		}
		return &ClientDetails{
			ClientId:             "ValidClientId",
			ValidRedirectUris:    []string{"myapp1.com", "myapp2.com", "myapp3.com"},
			AllowedResponseTypes: []string{"code", "token"},
		}
	}

	redirectUri := "myapp2.com"

	req := AuthorizationReq{
		ResponseType: "token",
		ClientId:     "ValidClientId",
		RedirectUri:  &redirectUri,
		ClientSecret: "secret",
		Scope:        "",
		State:        "mystate",

		UsePkce:             false,
		CodeChallenge:       "",
		CodeChallengeMethod: "s256",
	}

	got := req.VerifyAuthorizationRequest(getClients)

	if diff := cmp.Diff(nil, got, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("TestVerifyAuthorizationRequest() mimatch (-want +got:\n%s)", diff)
	}
}
