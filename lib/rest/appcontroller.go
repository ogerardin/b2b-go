package rest

import (
	"b2b-go/meta"
	"github.com/gin-gonic/gin"
)

func RegisterAppRoutes(g *gin.Engine) gin.IRoutes {
	return g.GET("/app/version", func(c *gin.Context) {
		c.String(200, meta.Version)
	})
}
