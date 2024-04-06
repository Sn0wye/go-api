package jwt

import (
	"strings"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWT struct {
	key []byte
}

type MyCustomClaims struct {
	UserId string `json:"user_id"`
	jwtLib.RegisteredClaims
}

func NewJwt(conf *viper.Viper) *JWT {
	return &JWT{key: []byte(conf.GetString("security.jwt_secret"))}
}

func (j *JWT) GenToken(userId string, expiresAt time.Time) (string, error) {
	claims := MyCustomClaims{
		UserId: userId,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(expiresAt),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
			Issuer:    "go-api",
			Subject:   userId,
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	strippedToken := strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwtLib.ParseWithClaims(strippedToken, &MyCustomClaims{}, func(token *jwtLib.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
