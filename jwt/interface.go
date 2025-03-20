package jwt

import (
	"crypto/rsa"
	"errors"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Algorithm  gojwt.SigningMethod
	Secret     string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// type Algorithm int

var (
	HS256 = gojwt.SigningMethodHS256
	RS256 = gojwt.SigningMethodRS256
)

func New(cfg Config) (EpicJWT, error) {
	switch cfg.Algorithm {
	case HS256:
		if cfg.Secret == "" {
			return nil, errors.New("requires secret")
		}
		return NewHMAC(cfg.Secret, cfg.Algorithm), nil
	case RS256:
		if cfg.PublicKey == nil {
			return nil, errors.New("requires public key")
		}
		return NewRSA(cfg.PrivateKey, cfg.PublicKey, cfg.Algorithm), nil
	default:
		return nil, errors.New("unsupported algorithm")
	}
}

var (
	ErrTokenExpired            = errors.New("token is expired")
	ErrInvalidToken            = errors.New("token is invalid")
	ErrTokenMalformed          = errors.New("token is malformed")
	ErrTokenNotValidYet        = errors.New("token is not valid yet")
	ErrTokenSignatureInvalid   = errors.New("token signature is invalid")
	ErrUnexpectedSigningMethod = errors.New("token signing method deos not match")
)

type EpicJWT interface {
	Sign(claims gojwt.Claims) (string, error)
	Verify(token string) error
	Decode(token string, claims gojwt.Claims) error
}
