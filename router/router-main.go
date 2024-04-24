package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luo2pei4/base-server/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome Gin Server")
		logger.Info("hellow gin server")
	})
	router.GET("/metrics", prometheusHandler())
	return router
}

// gin方式的metrics handler
func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
