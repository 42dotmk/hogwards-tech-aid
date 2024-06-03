package main

import (
	"github.com/gin-gonic/gin"
	"github.com/42dotmk/hogwards-tech-aid/handlers"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		handlers.TestHandler(c.Writer, c.Request)
	})
	err := r.Run()
	if err != nil {
		return
	}
}
