package rest

import (
	"b2b-go/app/domain"
	"b2b-go/app/repo"
	"b2b-go/lib/typeregistry"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

func RegisterSourceRoutes(r repo.SourceRepo, g *gin.Engine) {

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

		source, ok := mapBody(c)
		if !ok {
			return
		}

		// Save into the repository
		_, err := r.SaveNew(source)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Finally reply with newly created object
		c.JSON(http.StatusOK, source)
	})

	g.PUT("/api/sources/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		objectId := bson.ObjectIdHex(id)

		var source domain.BackupSource
		err := c.BindJSON(&source)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		err = r.Update(objectId, source)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, source)
	})

}

func mapBody(c *gin.Context) (domain.BackupSource, bool) {
	// map the request body expected to be JSON into a generic map
	// we can't map directly into a struct because we don't knwow which struct yet
	var bodyAsMap map[string]interface{}
	err := c.BindJSON(&bodyAsMap)
	if err != nil {
		return nil, false
	}

	// We expect a "_t" field that holds am identifier of the concrete type (as registered using typeregistry.Register)
	// Check presence of _t field
	typefield := bodyAsMap["_t"]
	if typefield == nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("missing field '_t' in body"))
		return nil, false
	}

	// Get corresponding concrete type from the type registry
	typeref := typefield.(string)
	t := typeregistry.GetType(typeref)
	if t == nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Errorf("Unknown type '%s'", typeref))
		return nil, false
	}

	// Allocate a new value for the target type
	i := reflect.New(t).Interface()
	// make sure it can be converted to the expected interface
	source, ok := i.(domain.BackupSource)
	if !ok {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Errorf("Type %s does not implement %s", t, reflect.TypeOf((*domain.BackupSource)(nil)).Elem()))
		return nil, false
	}

	// map the map obtained from the JSON body into the new struct
	err = mapstructure.Decode(bodyAsMap, source)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return nil, false
	}

	return source, true
}
