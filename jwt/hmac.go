package jwt

import (
	"errors"
	"fmt"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type epicHmac struct {
	algorithm gojwt.SigningMethod
	secret    string
}

func NewHMAC(secret string, algorithm gojwt.SigningMethod) EpicJWT {
	return &epicHmac{secret: secret, algorithm: algorithm}
}

func (h *epicHmac) Sign(claims gojwt.Claims) (string, error) {
	token := gojwt.NewWithClaims(h.algorithm, claims)
	signedToken, err := token.SignedString([]byte(h.secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (h *epicHmac) Verify(tokenStr string) error {
	token, err := gojwt.Parse(tokenStr, func(token *gojwt.Token) (any, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.secret), nil
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

func (h *epicHmac) Decode(tokenStr string, claims gojwt.Claims) error {
	token, err := gojwt.ParseWithClaims(tokenStr, claims, func(token *gojwt.Token) (any, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.secret), nil
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
