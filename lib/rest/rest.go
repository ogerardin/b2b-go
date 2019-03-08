package rest

import (
	"github.com/gin-gonic/gin"
)

func GinProvider() *gin.Engine {
	g := gin.Default()
	registerAppRoutes(g)
	return g
}
