package handlers

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
	"github.com/k-orolevsk-y/gophermart/pkg/jwt"
	"github.com/k-orolevsk-y/gophermart/pkg/router"
)

type handlerService struct {
	logger *zap.Logger
	jwt    *jwt.Jwt
	pg     *repository.Pg
}

type apiService interface {
	GetRouter() *router.Router
	GetLogger() *zap.Logger
	GetPg() *repository.Pg
}

func ConfigureHandlersService(api apiService) {
	hs := &handlerService{
		logger: api.GetLogger(),
		jwt:    jwt.New(config.Config.HmacTokenSecret),
		pg:     api.GetPg(),
	}
	r := api.GetRouter()

	r.Group("/api", func(routerApi router.RouterGroup) {
		routerApi.Group("/user", func(routerUser router.RouterGroup) {
			routerUser.POST("/register", hs.Register)
			routerUser.POST("/login", hs.Login)

			routerUser.Use(hs.jwt.Middleware)

			routerUser.POST("/orders", nil)
			routerUser.Group("/balance", func(routerUserBalance router.RouterGroup) {
				routerUserBalance.GET("/", nil)
				routerUserBalance.POST("/withdraws", nil)
			})
			routerUser.GET("/withdraws", nil)
		})
	})

	r.HandleMethodNotAllowed = true
	r.NoRoute(hs.NoRoute)
	r.NoMethod(hs.NoMethod)
}
