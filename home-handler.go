package main

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	hxRender(c, gin.H{}, "home.html")
}
