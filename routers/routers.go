package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ozeer/sloth/config"
	"github.com/ozeer/sloth/controllers"
	"github.com/ozeer/sloth/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(c config.Conf) *gin.Engine {
	gin.SetMode(c.Gin.Mode)
	router := gin.New()
	pprof.Register(router)

	router.GET("/", controllers.Welcome)

	v1 := router.Group("v1")
	{
		v1.POST("/add_task", controllers.AddTask)
		v1.GET("/get_task", controllers.GetTask)
		v1.POST("/del_task", controllers.DelTask)
	}

	// 注册prometheus中间件
	gp := prometheus.New(router)
	router.Use(gp.Middleware())
	// metrics采样
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
