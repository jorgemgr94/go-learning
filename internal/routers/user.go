package routers

import (
	userhandlers "go-learning/internal/handlers"

	"github.com/gin-gonic/gin"
)

func UserRouter(routerGroup *gin.RouterGroup) *gin.RouterGroup {
	users := routerGroup.Group("/users")
	users.POST("", userhandlers.New())
	users.GET("", userhandlers.GetList())
	users.GET("/:id", userhandlers.GetById())

	return users
}
