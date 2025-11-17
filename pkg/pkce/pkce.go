package pkce

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const (
	CodeChallengePlain string = "plain"
	CodeChallengeS256  string = "s256"
)

var (
	ErrPkceVerification error = errors.New("pkce verification failed")
	ErrPkceMethod       error = errors.New("unsupported pkce method")
)

func Verify(codeChallengeMethod string, codeChallenge, codeVerifier string) error {
	switch strings.ToLower(codeChallengeMethod) {
	case CodeChallengePlain:
		if codeVerifier != codeChallenge {
			return ErrPkceVerification
		}
	case CodeChallengeS256:
		sum := sha256.Sum256([]byte(codeVerifier))
		encoded := base64.RawURLEncoding.EncodeToString(sum[:])
		if encoded != codeChallenge {
			return ErrPkceVerification
		}
	default:
		return fmt.Errorf("%w: %s", ErrPkceMethod, codeChallengeMethod)
	}

	return nil
}
