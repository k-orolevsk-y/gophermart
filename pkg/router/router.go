package router

import "github.com/gin-gonic/gin"

type Router struct {
	*gin.Engine
}

type routeGroupFunc func(r RouterGroup)

func New() *Router {
	return &Router{Engine: gin.New()}
}

func (r Router) Group(relativePath string, f routeGroupFunc) {
	group := r.Engine.Group(relativePath)
	f(RouterGroup{group})
}

type RouterGroup struct {
	*gin.RouterGroup
}

func (rg RouterGroup) Group(relativePath string, f routeGroupFunc) {
	group := rg.RouterGroup.Group(relativePath)
	f(RouterGroup{group})
}
