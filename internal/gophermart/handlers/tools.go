package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/gophermart/pkg/jwt"
)

func (hs *handlerService) GetTokenClaims(ctx *gin.Context) (*jwt.Claims, bool) {
	tokenClaimsAny, exists := ctx.Get("tokenClaims")
	if !exists {
		return nil, false
	}

	tokenClaims, ok := tokenClaimsAny.(*jwt.Claims)
	if !ok {
		return nil, false
	}

	return tokenClaims, true
}
