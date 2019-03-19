package rest

import (
	"b2b-go/lib"
	"github.com/gin-gonic/gin"
)

func registerSourceRoutes(r lib.SourceRepo, g *gin.Engine) {

	g.GET("/api/sources", func(c *gin.Context) {
		//TODO
	})
}
