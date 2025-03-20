package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type epicRSA struct {
	algorithm  gojwt.SigningMethod
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRSA(private *rsa.PrivateKey, public *rsa.PublicKey, algorithm gojwt.SigningMethod) EpicJWT {
	return &epicRSA{
		algorithm:  algorithm,
		privateKey: private,
		publicKey:  public,
	}
}

func (r *epicRSA) Sign(claims gojwt.Claims) (string, error) {
	if r.privateKey == nil {
		return "", errors.New("requires private key to sign jwt")
	}
	token := gojwt.NewWithClaims(r.algorithm, claims)
	signedToken, err := token.SignedString(r.privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (r *epicRSA) Verify(tokenStr string) error {
	token, err := gojwt.Parse(tokenStr, func(token *gojwt.Token) (any, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return r.publicKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, gojwt.ErrTokenExpired):
			return ErrTokenExpired
		case errors.Is(err, gojwt.ErrTokenMalformed):
			return ErrTokenMalformed
		case errors.Is(err, gojwt.ErrTokenNotValidYet):
			return ErrTokenNotValidYet
		case errors.Is(err, gojwt.ErrTokenSignatureInvalid):
			return ErrTokenSignatureInvalid
		case errors.Is(err, ErrUnexpectedSigningMethod):
			return ErrUnexpectedSigningMethod
		default:
			return fmt.Errorf("failed to parse token: %w", err)
		}
	}

	if !token.Valid {
		return ErrInvalidToken
	}
	return nil
}

func (r *epicRSA) Decode(tokenStr string, claims gojwt.Claims) error {
	token, err := gojwt.ParseWithClaims(tokenStr, claims, func(token *gojwt.Token) (any, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return r.publicKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, gojwt.ErrTokenExpired):
			return ErrTokenExpired
		case errors.Is(err, gojwt.ErrTokenMalformed):
			return ErrTokenMalformed
		case errors.Is(err, gojwt.ErrTokenNotValidYet):
			return ErrTokenNotValidYet
		case errors.Is(err, gojwt.ErrTokenSignatureInvalid):
			return ErrTokenSignatureInvalid
		case errors.Is(err, ErrUnexpectedSigningMethod):
			return ErrUnexpectedSigningMethod
		default:
			return fmt.Errorf("failed to parse token: %w", err)
		}
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}
