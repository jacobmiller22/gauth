package authorization

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
)

type mockClientReader struct {
	clients []ClientDetails
}

func (r mockClientReader) ClientById(ctx context.Context, clientId string) (*ClientDetails, error) {
	found := slices.IndexFunc(r.clients, func(c ClientDetails) bool {
		return c.Id == clientId
	})
	if found < 0 {
		return nil, fmt.Errorf("Client with id '%s' not found", clientId)
	}
	return &r.clients[found], nil
}

var mockClients = []ClientDetails{
	{
		Id:            "1",
		RedirectUris:  []string{"myapp.com"},
		ResponseTypes: []string{"code"},
	},
	{
		Id:            "2",
		RedirectUris:  []string{"app1.com", "app2.com"},
		ResponseTypes: []string{"code", "token"},
	},
}

func TestVerifyAuthorizationRequest(t *testing.T) {

	testCases := []struct {
		name    string
		clients []ClientDetails
		req     AuthorizationReq
		wantErr error
	}{
		{
			name:    "valid authorization request",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "1",
				ResponseType: "code",
				RedirectUri:  "myapp.com",
			},
			wantErr: nil,
		},
		{
			name:    "missing redirect uri multiple valid redirect uri",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "2",
				ResponseType: "code",
			},
			wantErr: ErrInvalidRedirectUri,
		},
		{
			name:    "missing redirect uri single valid redirect uri",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "1",
				ResponseType: "code",
			},
			wantErr: nil,
		},
		{
			name:    "invalid redirect uri single valid redirect uri",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "1",
				ResponseType: "code",
				RedirectUri:  "someinvaliduri.com",
			},
			wantErr: nil,
		},
		{
			name:    "invalid redirect uri multiple valid redirect uri",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "2",
				ResponseType: "code",
				RedirectUri:  "someinvaliduri.com",
			},
			wantErr: ErrInvalidRedirectUri,
		},
		{
			name:    "valid response type code",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "1",
				ResponseType: "code",
				RedirectUri:  "myapp.com",
			},
			wantErr: nil,
		},
		{
			name:    "valid response type token",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "2",
				ResponseType: "token",
				RedirectUri:  "app1.com",
			},
			wantErr: nil,
		},
		{
			name:    "invalid response type",
			clients: mockClients,
			req: AuthorizationReq{
				ClientId:     "1",
				ResponseType: "token",
				RedirectUri:  "myapp.com",
			},
			wantErr: ErrInvalidResponseType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clients := mockClientReader{clients: tc.clients}
			gotErr := tc.req.Verify(context.TODO(), clients)
			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("+want: %v; -got: %v", tc.wantErr, gotErr)
			}
		})
	}
}
