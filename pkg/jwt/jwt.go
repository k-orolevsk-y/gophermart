package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	hmacSecret []byte
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
}

func New(hmacSecret string) *Jwt {
	return &Jwt{
		hmacSecret: []byte(hmacSecret),
	}
}

func (j *Jwt) Encode(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           userID,
	})

	return token.SignedString(j.hmacSecret)
}

func (j *Jwt) Decode(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.hmacSecret, nil
	})

	if claims, ok := token.Claims.(Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (j *Jwt) Middleware(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	if claims, err := j.Decode(tokenString); err != nil {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
	} else {
		ctx.Set("tokenClaims", claims)
	}
}
