package rest

import (
	"b2b-go/app/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterSourceRoutes(r repo.SourceRepo, g *gin.Engine) {

	g.GET("/api/sources", func(c *gin.Context) {
		sources, err := r.GetAll()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, sources)
	})
}
