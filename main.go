package main

import (
	"app/models"

	"github.com/gin-gonic/gin"

	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	rate, err := limiter.NewRateFromFormatted("240-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance_iprate := limiter.New(store, rate)
	middleware_iprate := mgin.NewMiddleware(instance_iprate)

	models.Migrate()

	r := gin.Default()

	r.ForwardedByClientIP = true
	r.Use(middleware_iprate)
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")
	InitializeRoutes(r)
	r.Run()
}

func InitializeRoutes(r *gin.Engine) {
	r.GET("/", Home)

	BuildDonorRoutes(r.Group("/donors"))
	r.GET("/donations", func(c *gin.Context) { hxRender(c, gin.H{}, "donations-list.html") })
	r.GET("/inventory", func(c *gin.Context) { hxRender(c, gin.H{}, "inventory-list.html") })
	r.GET("/configurations", func(c *gin.Context) { hxRender(c, gin.H{}, "configurations-list.html") })
	r.GET("/team", func(c *gin.Context) { hxRender(c, gin.H{}, "team.html") })
}
