package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luo2pei4/base-server/logger"
	"github.com/luo2pei4/base-server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.GET("/", middleware.CollectAPIStats("main"), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome Gin Server")
		logger.Info("hellow gin server")
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return router
}
