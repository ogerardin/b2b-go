package rest

import (
	"b2b-go/app/domain"
	"b2b-go/app/repo"
	"b2b-go/lib/log4go"
	"b2b-go/lib/util"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

var logger = log4go.GetPackageLogger()

func RegisterSourceRoutes(r repo.SourceRepo, g *gin.Engine) {

	targetType := reflect.TypeOf((*domain.BackupSource)(nil)).Elem()

	g.GET("/api/sources", func(c *gin.Context) {
		sources, err := r.GetAll()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError,
				errors.Wrap(err, "Failed to retrieve sources from repository"))
			return
		}

		jsonOK(c, sources)
	})

	g.GET("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		source, err := r.GetById(id)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		jsonOK(c, source)
	})

	g.DELETE("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		err := r.Delete(id)
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

		source, ok := util.MapBody(c, targetType)
		if !ok {
			return
		}

		err := r.Update(id, source.(domain.BackupSource))
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

}

func jsonOK(c *gin.Context, data interface{}) {
	res, err := jsonapi.MarshalToStruct(data, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "Failed to marshal response"))
		return
	}
	c.JSON(http.StatusOK, res)
}
