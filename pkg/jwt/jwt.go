package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

type Jwt struct {
	hmacSecret []byte
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"uid"`
}

func New(hmacSecret string) *Jwt {
	return &Jwt{
		hmacSecret: []byte(hmacSecret),
	}
}

func (j *Jwt) Encode(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gophermart",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
		},
		UserID: user.ID,
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
	if err != nil {
		return nil, err
	}

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
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "An invalid token was transferred or it has expired")
	} else {
		ctx.Set("tokenClaims", claims)
	}
}
