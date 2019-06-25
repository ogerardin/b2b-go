package rest

import (
	"b2b-go/lib/log4go"
	"b2b-go/lib/log4go/logadapters"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinProvider() *gin.Engine {
	log := log4go.GetDefaultLogger()

	gin.DefaultWriter = &logadapters.WriterAdapter{
		Level:  logrus.InfoLevel,
		Logger: log,
	}
	gin.DefaultErrorWriter = &logadapters.WriterAdapter{
		Level:  logrus.ErrorLevel,
		Logger: log,
	}

	// default instance without any middleware
	engine := gin.New()

	// add default logger (will write to gin.DefaultWriter)
	engine.Use(gin.Logger())

	// add recovery (will write to gin.DefaultErrorWriter)
	engine.Use(gin.Recovery())

	return engine
}
