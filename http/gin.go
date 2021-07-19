package http

import (
	"http-probe/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartGin(c *config.Config) {
	r := gin.Default()
	Routes(r)
	r.Run(c.HttpListenAddr)
}

func Routes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/probe/http", HttpProbe)
	api.GET("/v1", func(c *gin.Context) {
		c.String(http.StatusOK, "你好我是 http prober")
	})
}
