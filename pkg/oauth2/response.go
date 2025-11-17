package oauth2

import (
	"strconv"

	"net/url"
)

type AuthorizationRes struct {
	State        string
	Code         string
	Token        string
	ResponseType string
	RedirectUri  *url.URL
}

// func (req *AuthorizationReq) ToAuthorizationRes() (*AuthorizationRes, error) {

// 	var code string
// 	var token string
// 	var err error

// 	switch req.ResponseType {
// 	case "code":
// 		code, err = req.AuthorizationCodeSimple("sub", 86400)
// 	case "token":
// 		token = ""
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("Error creating code or token") // TODO: Custom error
// 	}

// 	locationUri, err := url.Parse(*req.RedirectUri)

// 	if err != nil {
// 		return nil, fmt.Errorf("Invalid redirect_uri") // TODO: Custom error
// 	}

// 	return &AuthorizationRes{
// 		State:       req.State,
// 		Code:        code,
// 		Token:       token,
// 		RedirectUri: locationUri,
// 	}, nil
// }

// Taken from net/http
func cloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	u2 := new(url.URL)
	*u2 = *u
	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}
	return u2
}

func (r *AuthorizationRes) Location() (*url.URL, error) {
	locationUri := cloneURL(r.RedirectUri)
	existingQuery := locationUri.Query()
	existingQuery.Set("state", r.State)
	switch r.ResponseType {
	case "code":
		existingQuery.Set("code", r.Code)
	case "token":
		existingQuery.Set("token_type", "Bearer")
		existingQuery.Set("token", r.Token)
		existingQuery.Set("expires_in", strconv.Itoa(86400))
	}
	r.RedirectUri.RawQuery = existingQuery.Encode()
	return locationUri, nil
}
