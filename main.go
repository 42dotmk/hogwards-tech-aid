package main

import (
	"strconv"

	"github.com/42dotmk/hogwards/db"
	"github.com/42dotmk/hogwards/handlers"
	"github.com/42dotmk/hogwards/lib/crud"
	"github.com/42dotmk/hogwards/lib/renderers"
	"github.com/42dotmk/hogwards/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load() // load the .env file

	db := db.Connect()
	models.Migrate(db) // auto migrate the models and set the database

	app := setupApp()       // create the gin app and setup defaults
	registerRoutes(app, db) // register the routes

	setupDefaultRenderer()

	app.Run()
}

func setupApp() (app *gin.Engine) {
	app = gin.Default()
	app.Static("/assets", "./assets")
	app.LoadHTMLGlob("templates/*")
	return app
}

func setupDefaultRenderer() {
	config := renderers.NewConfig().
		WithLayout("home.html").
		WithEnableLayoutOnNonHxRequest(true).
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
	renderers.DefaultRenderSetup(config)
}
func registerRoutes(app *gin.Engine, db *gorm.DB) {

	//V0 just create a simple method handler for a simple route (the method can be pased as an argument it does not have to be a closure)
	app.GET("/", func(c *gin.Context) { renderers.Render(c, gin.H{}, "home.html") })

	// v1. create a controller for the donor model hosted at /donors implementing all its routes there
	// the positive aspect of this is that you can easily see all the routes for the donor model in one place
	// all the dependencies for the donor model are also in one place and passed to the controller explicitly
	handlers.NewDonorCtrl(db).
		OnRouter(app.Group("/donors")) // this method call with attach all the routes to the /donors router group

	// V2.0
	// If you need the basic CRUD operations for a model, you can just use the CrudCtrl (or inherit from it and add/override routes)
	crud.New[models.Tag](db).
		OnRouter(app.Group("/tags"))

	// V2.1
	// the repairman binder is a function that takes a gin context and returns a repairman or an error
	// the repairman binder is used to bind the form data to the repairman model
	crud.New[models.Repairman](db).
		OnRouter(app.Group("/team")).
		WithFormBinder(func(c *gin.Context, out *models.Repairman) error {
			idStr := c.Param("ID")
			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return err
				}
				out.ID = uint(id)
			}
			out.Name = c.PostForm("Name")
			out.Phone = c.PostForm("Phone")
			out.Email = c.PostForm("Email")
			out.Address = c.PostForm("Address")
			out.IsActive = c.PostForm("IsActive") == "on"
			return nil
		})

	// v2.2 you can add the routes manually if you like
	recepients := crud.New[models.Recipient](db)
	app.GET("/recipients", recepients.List)
}
