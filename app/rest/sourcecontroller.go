package rest

import (
	"b2b-go/app/domain"
	"b2b-go/app/repo"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
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

	g.GET("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		objectId := bson.ObjectIdHex(id)
		source, err := r.GetById(objectId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

	g.DELETE("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		objectId := bson.ObjectIdHex(id)
		err := r.Delete(objectId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusOK)
	})

	g.POST("/api/sources", func(c *gin.Context) {
		//var source domain.BackupSource
		var source domain.FilesystemSource
		err := c.Bind(&source)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		_, err = r.SaveNew(source)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

	g.PUT("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		objectId := bson.ObjectIdHex(id)

		var source domain.BackupSource
		err := c.Bind(&source)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		err = r.Update(objectId, source)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

}
