package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/order_pool"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
	"github.com/k-orolevsk-y/gophermart/pkg/jwt"
	"github.com/k-orolevsk-y/gophermart/pkg/router"
	"github.com/k-orolevsk-y/gophermart/pkg/validation"
)

type handlerService struct {
	logger    *zap.Logger
	jwt       *jwt.Jwt
	pg        repository.Repository
	orderPool pool

	bindingWithValidation binding.Binding
}

type apiService interface {
	GetRouter() *router.Router
	GetLogger() *zap.Logger
	GetPg() repository.Repository
	GetOrderPool() *orderpool.OrderPool
}

type pool interface {
	AddJob(order models.Order)
}

func ConfigureHandlersService(api apiService) {
	hs := &handlerService{
		logger:    api.GetLogger(),
		jwt:       jwt.New(config.Config.HmacTokenSecret),
		pg:        api.GetPg(),
		orderPool: api.GetOrderPool(),

		bindingWithValidation: validation.NewBindingWithValidation(),
	}
	r := api.GetRouter()

	r.Use(gin.CustomRecovery(hs.Recovery))
	r.Group("/api", func(routerApi router.RouterGroup) {
		routerApi.Group("/user", func(routerUser router.RouterGroup) {
			routerUser.POST("/register", hs.Register)
			routerUser.POST("/login", hs.Login)

			routerUser.Use(hs.jwt.Middleware)

			routerUser.GET("/orders", hs.GetOrders)
			routerUser.POST("/orders", hs.NewOrder)

			routerUser.Group("/balance", func(routerUserBalance router.RouterGroup) {
				routerUserBalance.GET("/", hs.GetBalance)
				routerUserBalance.POST("/withdraw", hs.NewBalanceWithdrawn)
			})
			routerUser.GET("/withdrawals", hs.GetWithdrawals)
		})
	})

	r.HandleMethodNotAllowed = true
	r.NoRoute(hs.NoRoute)
	r.NoMethod(hs.NoMethod)
}
