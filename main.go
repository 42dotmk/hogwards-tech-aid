package main

import (
	"github.com/42dotmk/hogwards/db"
	"github.com/42dotmk/hogwards/handlers"
	"github.com/42dotmk/hogwards/lib/renderers"
	"github.com/42dotmk/hogwards/models"
	"github.com/gin-gonic/gin"
)

func main() {
	models.Migrate()

	app := gin.Default()
	//a route for static files and under which route to serve them
	app.Static("/assets", "./assets")
	//tell gin where to search for templates
	app.LoadHTMLGlob("templates/*")

	// Setup the Htmx renderer defaults
	hxConfig := renderers.NewConfig().
		WithLayout("home.html"). //the default layout to use
		WithEnableLayoutOnNonHxRequest(true). 
		//a function that is used to supply the layout with data, in this case the menu items
		WithLayoutDataFunc(func(data gin.H) {
			data["Menu"] = []models.MenuItem{
				{Title: "Home", Uri: "/", IsEnabled: true, IsExternal: true},
				{Title: "Донори", Uri: "/donors", IsEnabled: true, IsExternal: false},
				{Title: "Донации", Uri: "/donations", IsEnabled: true, IsExternal: false},
				{Title: "Тим", Uri: "/team", IsEnabled: true, IsExternal: false},
				{Title: "Опрема", Uri: "/equipment", IsEnabled: true, IsExternal: false},
				{Title: "Конфигурации", Uri: "/bundles", IsEnabled: true, IsExternal: false},
			}
		})
	renderers.DefaultRenderSetup(hxConfig)

	//HANDLERS
	//V0 just create a simple method handler for a simple route (the method can be pased as an argument it does not have to be a closure)
	app.GET("/", func(c *gin.Context) { renderers.Render(c, gin.H{}, "home.html") })

	// v1. create a controller for the donor model hosted at /donors implementing all its routes there
	// the positive aspect of this is that you can easily see all the routes for the donor model in one place
	// all the dependencies for the donor model are also in one place and passed to the controller explicitly
	handlers.NewDonorCtrl(db.DB).
		OnRouter(app.Group("/donors")) // this method call with attach all the routes to the /donors router group

	// V2.0
	// If you need the basic CRUD operations for a model, you can just use the CrudCtrl (or inherit from it and add/override routes)
	handlers.NewCrudCtrl[models.Tag](db.DB).
		OnRouter(app.Group("/tags"))

	// V2.1
	// the repairman binder is a function that takes a gin context and returns a repairman or an error
	// the repairman binder is used to bind the form data to the repairman model
	handlers.NewCrudCtrl[models.Repairman](db.DB).
		OnRouter(app.Group("/team")).
		WithFormBinder(func(c *gin.Context, out *models.Repairman) error {
			out.Name = c.PostForm("Name")
			out.Phone = c.PostForm("Phone")
			out.Email = c.PostForm("Email")
			out.Address = c.PostForm("Address")
			out.IsActive = c.PostForm("IsActive") == "on"
			return nil
		})

	// v2.2 you can add the routes manually if you like
	recepients := handlers.NewCrudCtrl[models.Recipient](db.DB)
	app.GET("/recipients", recepients.List)

	// v2.3 Or you can have a custom implementation of the CRUD controller

	app.Run()
}
