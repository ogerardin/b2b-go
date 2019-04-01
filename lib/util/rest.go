package util

import (
	"b2b-go/lib/typeregistry"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

// Maps the body of the HTTP request from context c into a new struct whose type is expected to be found
// in a field named "_t" of the JSON body. If targetType is not nil, it is asserted that the type implements
// the specified targetType (which is expected to be an interface type).
// The type specified in field "_t" must have been registered using typeregistry.Register.
// In case of success, returns a pointer to the populated struct as an interface{}, and true.
// In case of error, calls c.AbortWithError, and returns nil and false
//
// To obtain the Type of an instance you could do something like:
//		targetType := reflect.TypeOf((*MyInterface)(nil)).Elem()
//
func MapBody(c *gin.Context, targetType reflect.Type) (interface{}, bool) {

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

	if targetType != nil && !t.Implements(targetType) {
		_ = c.AbortWithError(http.StatusBadRequest, errors.Errorf("Type %s does not implement %s", t, targetType))
		return nil, false
	}

	source := reflect.New(t).Interface()

	/*	// Allocate a new value for the target type
		i := reflect.New(t).Interface()
		// make sure it can be converted to the expected interface
		source, ok := i.(domain.BackupSource)
		if !ok {
			_ = c.AbortWithError(http.StatusBadRequest, errors.Errorf("Type %s does not implement %s", t, targetType))
			return nil, false
		}
	*/
	// map the map obtained from the JSON body into the new struct
	err = mapstructure.Decode(bodyAsMap, source)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return nil, false
	}

	return source, true
}
