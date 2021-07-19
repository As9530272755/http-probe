package http

import (
	"fmt"
	"http-probe/probe"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpProbe(c *gin.Context) {
	host := c.Query("host")
	isHttps := c.Query("is_https")

	if host == "" {
		c.String(http.StatusBadRequest, "空的主机")
		return
	}

	schema := "http"
	if isHttps == "1" {
		schema = "https"
	}

	url := fmt.Sprintf("%s://%s", schema, host)
	res := probe.DoHttpProbe(url)
	c.String(http.StatusOK, res)

}
