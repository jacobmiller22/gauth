package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"time"
)

type Claims interface {
	GetExpirationTime() (time.Time, error)
	GetIssuedAt() (time.Time, error)
	GetNotBefore() (time.Time, error)
	GetIssuer() (string, error)
	GetSubject() (string, error)
	GetAudience() ([]string, error)
	GetId() (string, error)
}

type Decoder interface {
	Decode(v any) error
}

type Encoder interface {
	Encode(v any) error
}

type JWTHeader struct {
	Algorithm   string `json:"alg"`
	Type        string `json:"typ"`
	ContentType string `json:"cty"`
}

func jsonify(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil

}

type JWTPayload struct {
	Issuer     string `json:"iss"`
	Subject    string `json:"sub"`
	Audience   string `json:"aud"` // OPTIONAL
	Expiration int    `json:"exp"` // OPTIONAL
	NotBefore  int    `json:"nbf"`
	IssuedAt   int    `json:"iat"`
	JwtId      string `json:"jti"`
}

type JWT struct {
	Header    JWTHeader
	Payload   JWTPayload
	Signature string
}

func (jwt *JWT) Sign(key []byte) (string, error) {

	hmac := hmac.New(sha256.New, []byte(key))

	hjson, err := jsonify(jwt.Header)
	if err != nil {
		return "", err
	}
	pjson, err := jsonify(jwt.Payload)
	if err != nil {
		return "", err
	}
	hmac.Write(append(append(hjson, "."...), pjson...))
	signature := hmac.Sum(nil)
	return string(hjson) + "." + string(pjson) + "." + string(signature), nil
}

func NewJWT(header JWTHeader, payload JWTPayload, secret string) *JWT {
	return &JWT{
		Header:    header,
		Payload:   payload,
		Signature: "",
	}
}
