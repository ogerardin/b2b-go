package rest

import (
	"b2b-go/app/domain"
	"b2b-go/app/repo"
	"b2b-go/lib/util"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

func RegisterSourceRoutes(r repo.SourceRepo, g *gin.Engine) {

	targetType := reflect.TypeOf((*domain.BackupSource)(nil)).Elem()

	g.GET("/api/sources", func(c *gin.Context) {
		sources, err := r.GetAll()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, sources)
	})

	g.GET("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		objectId := bson.ObjectIdHex(id)
		source, err := r.GetById(objectId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

	g.DELETE("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		objectId := bson.ObjectIdHex(id)
		err := r.Delete(objectId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusOK)
	})

	g.POST("/api/sources", func(c *gin.Context) {

		source, ok := util.MapBody(c, targetType)
		if !ok {
			return
		}

		// Save into the repository
		_, err := r.SaveNew(source.(domain.BackupSource))
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Finally reply with newly created object
		c.JSON(http.StatusOK, source)
	})

	g.PUT("/api/sources/:id", func(c *gin.Context) {

		id := c.Params.ByName("id")
		if !bson.IsObjectIdHex(id) {
			_ = c.AbortWithError(http.StatusBadRequest, errors.Errorf("Not a valid object ID: %s", id))
			return
		}
		objectId := bson.ObjectIdHex(id)

		source, ok := util.MapBody(c, targetType)
		if !ok {
			return
		}

		err := r.Update(objectId, source.(domain.BackupSource))
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

}
