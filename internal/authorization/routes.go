package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/jacobmiller22/gauth/pkg/authorization"
)

type Routes struct {
	S *Service
}

func (rte *Routes) Authorize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		clientId := q.Get("client_id")
		redirectUri := q.Get("redirect_uri")
		responseType := q.Get("response_type")
		scope := q.Get("scope")
		state := q.Get("state")
		codeChallengeMethod := q.Get("code_challenge_method")
		codeChallenge := q.Get("code_challenge")

		authReq, err := authorization.NewAuthorizationRequest(
			clientId,
			responseType,
			redirectUri,
			scope,
			state,
			codeChallengeMethod, codeChallenge,
		)

		if err != nil {
			w.WriteHeader(500) // TODO: Finish this
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		userId := "30917f3d-69c3-4bf1-883c-3a4b0f061b1e" // TODO: Figure out how to get this

		authRes, err := rte.S.Authorize(r.Context(), authReq, userId)
		if err != nil {
			w.WriteHeader(500) // TODO: Finish this
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		resLocation, err := authRes.Location()
		if err != nil {
			w.WriteHeader(500) // TODO: Finish this
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Location", resLocation.String())

		if authRes.ResponseType == authorization.ResponseTypeToken {
			w.Header().Set("Cache-Control", "no-store")
		}
		w.WriteHeader(http.StatusFound)
	}
}

func (rte *Routes) Token() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		tokenReq := authorization.TokenRequest{
			GrantType:    q.Get("grant_type"),
			Code:         q.Get("code"),
			RedirectUri:  q.Get("redirect_uri"),
			Username:     q.Get("username"),
			Password:     q.Get("password"),
			ClientId:     q.Get("client_id"),
			ClientSecret: q.Get("client_secret"),
			Scope:        q.Get("scope"),
			CodeVerifier: q.Get("code_verifier"),
		}

		tokenRes, err := rte.S.Token(r.Context(), &tokenReq)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(500) // TODO: Finish this
			return
		}

		w.Header().Set("Cache-Control", "no-store")

		jsonRes, err := json.Marshal(tokenRes)
		if err != nil {
			w.WriteHeader(500) // TODO: Finish this
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		if n, err := w.Write(jsonRes); err != nil || n < len(jsonRes) {
			// log error
		}
	}
}
