package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/max-gui/logagent/pkg/routerutil"
	// nethttp "net/http"
)

func SetupRouter() *gin.Engine {
	// gin.New()

	r := gin.New()                      //.Default()
	r.Use(routerutil.GinHeaderMiddle()) // ginHeaderMiddle())
	r.Use(routerutil.GinLogger())       //LoggerWithConfig())
	r.Use(routerutil.GinErrorMiddle())  //ginErrorMiddle())

	// r.Any("/eurekaagent/:appname/:env/*path", regagent.Eurekaagent)
	// r.Any("/consulagent/:appname/:env/*path", regagent.Consulagent)
	r.GET("/actuator/health", health)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "online")
}
