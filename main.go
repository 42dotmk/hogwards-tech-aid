package main

import (
	"app/models"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var router *gin.Engine

func main() {

	// IP RATE LIMITER
	rate, err := limiter.NewRateFromFormatted("240-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance_iprate := limiter.New(store, rate)
	middleware_iprate := mgin.NewMiddleware(instance_iprate)

	models.Migrate()

	router = gin.Default()

	router.ForwardedByClientIP = true
	router.Use(middleware_iprate)

	//templ := template.Must(template.New("").ParseFS(embeddedFiles, "templates/*"))
	//router.SetHTMLTemplate(templ)
	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("templates/*")

	InitializeRoutes()

	router.Run()

}

// Render one of HTML or JSON based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
