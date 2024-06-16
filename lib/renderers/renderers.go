package renderers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Config is used to configure the renderer and set the default layout and layout data function to use
// This config is than used by the Render function to render the data with the correct layout and data
type Config struct {
    //the layout to use on the templates for full page rendering
	DefaultLayout              string

    //a function that is used to supply the layout with data
	LayoutDataFunc             func(data gin.H)

    //if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
	EnableLayoutOnNonHxRequest bool
}


//the layout to use on the templates for full page rendering
func (c *Config) WithLayout(layout string) *Config {
    c.DefaultLayout = layout
    return c
}

//a function that is used to supply the layout with data
func (c *Config) WithLayoutDataFunc(layoutDataFunc func(data gin.H)) *Config {
    c.LayoutDataFunc = layoutDataFunc
    return c
}

//if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
func (c *Config) WithEnableLayoutOnNonHxRequest(enableLayoutOnNonHxRequest bool) *Config {
    c.EnableLayoutOnNonHxRequest = enableLayoutOnNonHxRequest
    return c
}

// NewConfig returns a new renderer Config with the default values
func NewConfig() *Config {
	return &Config{
		DefaultLayout:              "index.html",
		LayoutDataFunc:             nil,
		EnableLayoutOnNonHxRequest: true,
	}
}

// config is the default configuration for the Render function
var config = NewConfig()

// sets the default configuration for the Render function
func DefaultRenderSetup(options *Config) {
	if options.DefaultLayout != "" {
		config.DefaultLayout = options.DefaultLayout
	}

	if options.LayoutDataFunc != nil {
		config.LayoutDataFunc = options.LayoutDataFunc
	}
}


// Render is a function that renders a template with the given data
// If the request is request accepts application/json it will return the data as json
// If the request is an Htmx request it will render the template with the data
// If the request is not an Htmx request it will use the layout to render the data
//  - The layout should be aware of the data that is passed to it and conditionally render that template
// it uses the default hxConfig to render the template
// See `RenderWithConfig` for more control over the rendering
func Render(c *gin.Context, data gin.H, templateName string) {
	hxAwareRender(c, data, templateName, config.DefaultLayout, config)
}

// RenderWithConfig is a function that renders a template with the given data and the render configuration
// See Render for more information on the rendering
// See Config for more information on the configuration
func RenderWithConfig(c *gin.Context, data gin.H, templateName string, conf *Config) {
	hxAwareRender(c, data, templateName, config.DefaultLayout, conf)
}

func hxAwareRender(c *gin.Context, data gin.H, templateName string, layout string, conf *Config) {
	if conf.EnableLayoutOnNonHxRequest && c.Request.Header.Get("Hx-Request") == "true" {
		c.HTML(http.StatusOK, templateName, data)
	} else {
		// extend the data
		if conf.LayoutDataFunc != nil {
			config.LayoutDataFunc(data)
		}
		render(c, data, layout)
	}
}

// render is a helper function that renders the data based on the request accept header
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data)
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
