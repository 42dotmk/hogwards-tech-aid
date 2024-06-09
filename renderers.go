package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func render(c *gin.Context, data gin.H, templateName string) {
	// extendData(c, &data)
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}

var DefaultLayout = "home.html"

func hxRender(c *gin.Context, data gin.H, templateName string) {
	hxRenderWithLayout(c, data, templateName, DefaultLayout)
}

func hxRenderWithLayout(c *gin.Context, data gin.H, templateName string, layout string) {
	if c.Request.Header.Get("Hx-Request") == "true" {
		c.HTML(http.StatusOK, templateName, data)
	} else {
		c.HTML(http.StatusOK, layout, data)
	}
}
