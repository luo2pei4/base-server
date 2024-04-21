package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luo2pei4/base-server/logger"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome Gin Server")
		logger.Info("hellow gin server")
	})
	return router
}
