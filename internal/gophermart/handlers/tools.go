package handlers

import (
	"errors"
	"strconv"

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

var (
	errInvalidNumberForAlgorithmLuna       = errors.New("invalid number")
	errInvalidTypeOfNumberForAlgorithmLuna = errors.New("invalid type")
)

func (hs *handlerService) CheckNumberAlgorithmLuna(n any) (int64, error) {
	var number int64
	switch num := n.(type) {
	case int64:
		number = num
	case string:
		var err error
		number, err = strconv.ParseInt(num, 10, 64)

		if err != nil {
			return 0, errInvalidTypeOfNumberForAlgorithmLuna
		}
	}

	if number == 0 {
		return 0, errInvalidNumberForAlgorithmLuna
	}

	checkSum := func(number int64) int64 {
		var controlSum int64

		for i := 0; number > 0; i++ {
			cur := number % 10

			if i%2 == 0 {
				cur *= 2
				if cur > 9 {
					cur = cur%10 + cur/10
				}
			}

			controlSum += cur
			number /= 10
		}
		return controlSum % 10
	}

	if (number%10+checkSum(number/10))%10 == 0 {
		return number, nil
	}

	return 0, errInvalidNumberForAlgorithmLuna
}
