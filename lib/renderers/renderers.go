package renderers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MenuItem struct {
	Title     string
	Uri       string
	IsEnabled bool
    IsExternal bool
}

var DefaultLayoutDataFunc = func(data gin.H) {
	data["Menu"] = []MenuItem{
		{"Home", "/", true, true},
		{"Донори", "/donors", true, false}, 
		{"Донации", "/donations", true, false}, 
		{"Тим", "/team", true, false}, 
		{"Опрема", "/equipment", true, false}, 
		{"Конфигурации", "/bundles", true, false}, 
		{"Приматели", "/recepients", true, false}, 
	}
}

type HxConfig struct {
	DefaultLayout              string
	LayoutDataFunc             func(data gin.H)
	EnableLayoutOnNonHxRequest bool
}

var hxConfig = DefaulHxConfig()

func DefaulHxConfig() *HxConfig {
	return &HxConfig{
		DefaultLayout:              "home.html",
		LayoutDataFunc:             DefaultLayoutDataFunc,
		EnableLayoutOnNonHxRequest: true,
	}
}

func HxRenderSetup(options *HxConfig) {
	if options.DefaultLayout != "" {
		hxConfig.DefaultLayout = options.DefaultLayout
	}

	if options.LayoutDataFunc != nil {
		hxConfig.LayoutDataFunc = options.LayoutDataFunc
	}
}

func HxRender(c *gin.Context, data gin.H, templateName string) {
	HxRenderWithLayout(c, data, templateName, hxConfig.DefaultLayout)
}

func HxRenderWithLayout(c *gin.Context, data gin.H, templateName string, layout string) {
	if c.Request.Header.Get("Hx-Request") == "true" {
		c.HTML(http.StatusOK, templateName, data)
	} else {
		// extend the data
		if hxConfig.LayoutDataFunc != nil {
			hxConfig.LayoutDataFunc(data)
		}
		render(c, data, layout)
	}
}

func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data)
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
