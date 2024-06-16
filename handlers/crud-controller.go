package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	rnd "github.com/42dotmk/hogwards/lib/renderers"
	"github.com/42dotmk/hogwards/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FormHandler is a function that binds the form data to a model
type FormHandler[T models.IModel] func(c *gin.Context, out *T) error

// ICrudCtrl is an interface that defines the basic CRUD operations for a model
type ICrudCtrl [T models.IModel] interface {
    List(c *gin.Context) ICrudCtrl[T]
    Details(c *gin.Context) ICrudCtrl[T]
    EditForm(c *gin.Context) ICrudCtrl[T]
    Upsert(c *gin.Context) ICrudCtrl[T]
    Delete(c *gin.Context) ICrudCtrl[T]
}

// CrudCtrl is a controller that implements the basic CRUD operations for a model with a gorm backend
type CrudCtrl[T models.IModel] struct {
	Db          *gorm.DB
	Router      *gin.RouterGroup
	ModelName   string
	FormHandler FormHandler[T]
}
// NewCrudCtrl creates a new CRUD controller for the provided model
func NewCrudCtrl[T models.IModel](db *gorm.DB) *CrudCtrl[T] {
	var name = fmt.Sprintf("%T", *new(T))
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	name = strings.ToLower(name)

	return &CrudCtrl[T]{
		FormHandler: DefaultFormHandler[T], // default form handler is used if none is provided
		ModelName:   name,
		Db:          db,
	}
}
// OnRouter attaches the CRUD routes to the provided router
func (self *CrudCtrl[T]) OnRouter(r *gin.RouterGroup) *CrudCtrl[T] {
	r.GET("/", self.List)
	r.GET("/new", self.EditForm)
	r.GET("/:id", self.Details)
	r.GET("/:id/edit", self.EditForm)
	r.POST("/:id", self.Upsert)
	r.PUT("/", self.Upsert)
	r.DELETE("/:id", self.Delete)
	self.Router = r
	return self
}

// WithFormBinder sets the form binder for the controller to be used when binding form data to the model on the POST and PUT requests.
// It is used in the Upsert method of the controller
// If not set, the default form binder is used which assumes that the form field names are the same as the model field names(case sensitive)
func (self *CrudCtrl[T]) WithFormBinder(handler FormHandler[T]) *CrudCtrl[T] {
    self.FormHandler = handler
    return self
}
// List is a handler that lists all the items of the model
// it is a GET request 
// !Requres the template to be named as modelName-list.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) List(c *gin.Context) {
	var items []T
	self.Db.Find(&items)
	rnd.Render(c,
		gin.H{fmt.Sprintf("%sList", self.ModelName): items},
		fmt.Sprintf("%s-list.html", self.ModelName))
}

// Details is a handler that shows the details of a single item of the model
// it is a GET request
// !Requires the template to be named as modelName-details.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) Details(c *gin.Context) {
	template := fmt.Sprintf("%s-details.html", self.ModelName)
	modelName := fmt.Sprintf("%s", self.ModelName)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		var item T
		self.Db.First(&item, id)
		rnd.Render(c, gin.H{modelName: item}, template)
	} else {
		c.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
	}
}

// EditForm is a handler that shows the form for editing an item of the model
// it is a GET request
// !Requires the template to be named as modelName-edit.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) EditForm(c *gin.Context) {
	template := fmt.Sprintf("%s-edit.html", self.ModelName)
	modelName := fmt.Sprintf("%s", self.ModelName)
	idStr := c.Param("id")
	if idStr == "" {
		rnd.Render(c, gin.H{}, template)
		return
	} else {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			var item T
			self.Db.First(&item, idStr)
			rnd.Render(c, gin.H{modelName: item}, template)
		} else {
			c.Error(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
		}
	}
}

// Upsert is a handler that saves an item of the model
// it is a POST or PUT request depending on the presence of the id parameter
// !Requires the form fields to be named as the model field names(case sensitive)
// It redirects to the details page of the saved item
func (self *CrudCtrl[T]) Upsert(c *gin.Context) {
    var item T
	if err := self.FormHandler(c, &item); err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}
	idStr := c.Param("id")
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID: %d", id)
			c.String(http.StatusBadRequest, msg)
			c.Abort()
			return
		}
		item.SetID(uint(id))
	}
    print("saving item")
    res:= self.Db.Save(&item)
    if res.Error != nil {
        c.Error(res.Error)
        c.String(http.StatusBadRequest, res.Error.Error())
        c.Abort()
        return
    }

	c.Header("HX-Redirect", fmt.Sprintf("/%s/%d", self.Router.BasePath(), item.GetID()))
	c.String(http.StatusOK, "Saved")
	c.Abort()
}

// Delete is a handler that deletes an item of the model
// it is a DELETE request
func (self *CrudCtrl[T]) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid ID: %d", id)
		c.String(http.StatusBadRequest, msg)
		c.Abort()
		return
	}
	var item T
	self.Db.Delete(&item, id)
	c.Header("HX-Redirect", fmt.Sprintf("%s/%s", self.Router.BasePath(), self.ModelName))
	c.String(http.StatusOK, "Deleted")
	c.Abort()
}



func BindForm[T any](r *http.Request, out *T) error {
    // Parse the form data from the request.
    if err := r.ParseForm(); err != nil {
        return err
    }
    // Reflect on the struct to set values.
    val := reflect.ValueOf(out).Elem()
    typ := val.Type()
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        // Get form value for the field name.
        formValue := r.FormValue(fieldType.Name)
        // Check if the field can be set and if the form value is not empty.
        if field.CanSet() && formValue != "" {
            // Convert form values to the appropriate field types.
            // This example assumes all fields are strings for simplicity.
            // You might need to convert this based on the field type.
            field.SetString(formValue)
        }
    }
    return nil
}

func DefaultFormHandler[T models.IModel](c *gin.Context, out *T) error {
    if err := c.Request.ParseForm(); err != nil {
        return err
    }
    val := reflect.ValueOf(out).Elem()
    typ := val.Type()
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        formValue := c.PostForm(fieldType.Name)

        // we only set the field if we are able to do so and the form value is not empty
        if field.CanSet() && formValue != "" {
            field.SetString(formValue)
        }
    }
    return nil
}

