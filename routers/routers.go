package routers

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/luo2pei4/base-server/internal/hanlders"
	"github.com/luo2pei4/base-server/internal/logger"
	"github.com/luo2pei4/base-server/metrics"
	"github.com/luo2pei4/base-server/middleware"
)

func InitRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	// 创建router
	router := gin.New()
	// 注册pprof，默认地址/debug/pprof
	pprof.Register(router)
	// 注册监控指标路由
	router.GET("/metrics", gin.WrapH(metrics.PrometheusHandler()))
	// 注册根路径路由处理
	router.GET("/", middleware.CollectAPIStats("roog"), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome Gin Server")
		logger.Info("hellow gin server")
	})

	// login
	router.POST("/login", middleware.CollectAPIStats("login"), hanlders.LoginHandler)

	return router
}
