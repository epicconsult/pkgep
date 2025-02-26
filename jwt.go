package pkgep

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type MetaToken struct {
	ID            string
	Email         string
	ExpiredAt     time.Time
	Authorization bool
}

type AccessToken struct {
	Claims MetaToken
}

type SubClaims struct {
	UserID      int    `json:"user_id"`
	ParID       int    `json:"par_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	DeviceID    string `json:"device_id"`
	AuthType    string `json:"auth_type"`
}

// MyCustomClaims defines the structure of the entire JWT payload
type VerifiedToken struct {
	Sub  SubClaims `json:"sub"`
	Jti  string    `json:"jti"`
	Type string    `json:"type"`
	jwt.RegisteredClaims
}

func Sign(Data map[string]interface{}, SecretPublicKeyEnvName string, ExpiredAt time.Duration) (string, error) {
	expiredAt := time.Now().Add(time.Duration(time.Minute) * ExpiredAt).Unix()

	jwtSecretKey := GodotEnv(SecretPublicKeyEnvName)

	claims := jwt.MapClaims{}
	claims["exp"] = expiredAt
	claims["authorization"] = true

	for i, v := range Data {
		claims[i] = v
	}

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(jwtSecretKey))

	if err != nil {
		logrus.Error(err.Error())
		return accessToken, err
	}

	return accessToken, nil
}

func VerifyTokenHeader(ctx *fiber.Ctx, SecretPublicKeyEnvName string) (*jwt.Token, *VerifiedToken, error) {

	tokenHeader := ctx.Get("Authorization")
	bearerToken := strings.Split(tokenHeader, "Bearer")
	if len(bearerToken) < 2 {
		return nil, nil, errors.New("invalid token format")
	}
	accessToken := bearerToken[1]

	pubkey, err := loadRSAPublicKey("certs/public.key")
	if err != nil {
		return nil, nil, err
	}

	token, claims, err := verifyRS256Token(strings.TrimSpace(accessToken), pubkey)
	if err != nil {
		logrus.Error(err.Error())
		return nil, nil, err
	}

	return token, claims, nil

}

func loadRSAPublicKey(filePath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pubKey, nil
}

func verifyRS256Token(tokenString string, pubKey *rsa.PublicKey) (*jwt.Token, *VerifiedToken, error) {

	claims := &VerifiedToken{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	return token, claims, err
}

func JWTProtected() fiber.Handler {

	return func(c *fiber.Ctx) error {
		if c.Get("Authorization") == "" {
			log.Println("Missing or malformed JWT")
			return HandleErrorResponse(c, InvalidToken, "")
		}

		_, _, err := VerifyTokenHeader(c, "JWT_SECRET_KEY")
		if err != nil {
			log.Println("invalid token")
			return HandleErrorResponse(c, InvalidToken, "")
		}
		return c.Status(fiber.StatusOK).Next()
	}
}

func JwtLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {

		NewLogger().LogInformation(HTTPREQEST, c)

		NewHelpers(*NewLogger())

		return c.Next()
	}
}
