package util

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func VerifiedTokenAuth(tokenAuth string, secretKey string) (string, bool, string) {

	authorized := true
	jti := ""
	message := ""

	if tokenAuth != "" {

		// validate jwt token
		token, err := jwt.Parse(tokenAuth, func(token *jwt.Token) (any, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// var j jwtToken

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secretKey), nil
		})

		if err != nil {
			return "", false, "token tidak valid"
		}

		if !token.Valid {
			return "", false, "token tidak valid"
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return "", false, "token tidak valid"
		}

		if claims["jti"] == nil {
			return "", false, "token tidak valid"
		} else {
			jti = claims["jti"].(string)
		}

	} else {

		return "", false, "token kosong"

	}

	return jti, authorized, message
}
