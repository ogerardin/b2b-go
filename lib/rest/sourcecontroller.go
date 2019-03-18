package rest

import (
	"b2b-go/lib/domain/repo"
	"github.com/gin-gonic/gin"
)

func registerSourceRoutes(r repo.SourceRepo, g *gin.Engine) {

	g.GET("/api/sources", func(c *gin.Context) {
		//TODO
	})
}
