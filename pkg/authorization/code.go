package authorization

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthCodeReader interface {
	AuthCode(code string) (*AuthorizationCode, error)
}

type AuthorizationCode struct {
	Id          string
	ClientId    string
	RedirectUri string
	UserId      string

	CodeChallengeMethod string
	CodeChallenge       string

	IssuedAt  time.Time
	ExpiresIn time.Duration
}

func (c *AuthorizationCode) String() string {
	// TODO: Make this a jwt
	return fmt.Sprintf("%s~~%s~~%s~~%s~~%s~~%s~~%s~~%s", c.Id, c.ClientId, c.RedirectUri, c.UserId, c.CodeChallengeMethod, c.CodeChallenge, c.IssuedAt.Format(time.RFC3339), c.ExpiresIn.String())
}

func (c *AuthorizationCode) Expired(nowFunc func() time.Time) bool {
	return nowFunc().Before(c.IssuedAt.Add(c.ExpiresIn))
}

/*
Create a "simple" authorization code. Do not use in production. Simple concats the necessary information
in a string together. Used for testing purposes
*/
// func (r *AuthorizationReq) AuthorizationCodeSimple(sub string, lifetime int) (string, error) {
// 	codeId, err := funcs.GenerateRandomString(11)
// 	if err != nil {
// 		return "", err
// 	}
// 	return r.clientId + ":" + r.redirectUri + ":" + sub + ":" + strconv.Itoa(lifetime) + ":" + codeId + ":" + r.codeChallenge + ":" + r.codeChallengeMethod, nil
// }

func (r *AuthorizationReq) AuthorizationCodeSimple(userId string) (*AuthorizationCode, error) {
	return &AuthorizationCode{
		Id:          uuid.NewString(),
		ClientId:    r.ClientId,
		RedirectUri: r.RedirectUri,
		UserId:      userId,

		CodeChallengeMethod: r.CodeChallengeMethod,
		CodeChallenge:       r.CodeChallenge,

		IssuedAt:  time.Now(),
		ExpiresIn: time.Minute * 10,
	}, nil
}

// func (c AuthorizationCode) GetExpirationTime() (time.Time, error) {
// 	return time.Time{}, nil
// }
// func (c AuthorizationCode) GetIssuedAt() (time.Time, error) {
// 	return time.Time{}, nil

// }
// func (c AuthorizationCode) GetNotBefore() (time.Time, error) {
// 	return time.Time{}, nil

// }
// func (c AuthorizationCode) GetIssuer() (string, error) {
// 	return "", nil
// }
// func (c AuthorizationCode) GetSubject() (string, error) {
// 	return "", nil

// }
// func (c AuthorizationCode) GetAudience() ([]string, error) {
// 	return nil, nil
// }
