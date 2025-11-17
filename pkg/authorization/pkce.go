package authorization

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	CodeChallengePlain string = "plain"
	CodeChallengeS256  string = "s256"
)

var (
	ErrPkceVerification error = errors.New("pkce verification failed")
	ErrPkceMethod       error = errors.New("unsupported pkce method")
)

func (r *TokenRequest) PkceVerify(codeChallenge string) error {
	CODE_CHALLENGE_METHOD := "s256"
	switch CODE_CHALLENGE_METHOD {
	case CodeChallengePlain:
		if r.CodeVerifier != codeChallenge {
			return ErrPkceVerification
		}
		return nil
	case CodeChallengeS256:
		sum := sha256.Sum256([]byte(r.CodeVerifier))
		encoded := base64.RawURLEncoding.EncodeToString(sum[:])
		if encoded != codeChallenge {
			return ErrPkceVerification
		}
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrPkceMethod, CODE_CHALLENGE_METHOD)
	}
}
