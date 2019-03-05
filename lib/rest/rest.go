package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func StartApi() {
	g := gin.Default()

	registerAppRoutes(g)

	err := g.Run(":8080")
	if err != nil {
		panic(errors.Wrap(err, "Failed to start HTTP server"))
	}
}
