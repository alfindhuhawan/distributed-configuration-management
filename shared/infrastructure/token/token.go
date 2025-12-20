package token

import (
	"distributed-configuration-management/shared/infrastructure/config"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTToken interface {
	// CreateToken create a token with a content
	CreateToken(content []byte, expired time.Duration) (string, error)

	// VerifyToken verify and return the content
	VerifyToken(tokenString string) ([]byte, error)

	// Create pgw token
	CreatePGWToken(content []byte, expired time.Duration) (string, error)

	// CREATE pgw token
	CreatePGWTokenNew(apiToken string, expired time.Duration, secret string) (string, error)
}

const fieldContent = "content"

type jwtToken struct {
	secretKey string
}

type JwtClaims struct {
	// add as necessary
	jwt.StandardClaims
}

func NewJWTToken(cfg *config.Config) (JWTToken, error) {

	secretKey := cfg.JWTSecretKey
	if strings.TrimSpace(secretKey) == "" {
		return nil, fmt.Errorf("SecretKey must not empty")
	}

	return &jwtToken{
		secretKey: secretKey,
	}, nil
}

func (j jwtToken) CreateToken(content []byte, expired time.Duration) (string, error) {

	contentBase64 := base64.StdEncoding.EncodeToString(content)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(expired).Unix(),
		fieldContent: contentBase64,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (j jwtToken) VerifyToken(tokenString string) ([]byte, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Claims is can not asserted")
	}

	decodeStringInBytes, err := base64.StdEncoding.DecodeString(claims[fieldContent].(string))
	if err != nil {
		return nil, err
	}

	return decodeStringInBytes, nil
}

func (j jwtToken) CreatePGWToken(content []byte, expired time.Duration) (string, error) {

	contentBase64 := base64.StdEncoding.EncodeToString(content)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp": time.Now().Add(expired).Unix(),
		"jti": contentBase64,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("122d63d5847deab4966dcd5dcfd27759dc945ea8c398e41d2bcbe70f9a267b90"))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (j *jwtToken) CreatePGWTokenNew(apiToken string, expired time.Duration, secret string) (string, error) {
	claims := JwtClaims{
		jwt.StandardClaims{
			Id:        apiToken,
			ExpiresAt: time.Now().Add(2160 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, err
}
