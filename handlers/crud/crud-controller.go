package crud

import (
	"github.com/42dotmk/hogwards/models"
    rnd"github.com/42dotmk/hogwards/lib/renderers"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FormHandler[T models.IModel] func(c *gin.Context) (T, error)
type TmplMapper[T models.IModel] func(c *gin.Context) (string, error)

type CrudCtrl[T models.IModel] struct {
	Router      *gin.RouterGroup
	Db          *gorm.DB
	ModelName   string
	FormHandler FormHandler[T]
}
func (self * CrudCtrl[T]) baseName() string{
	var name = fmt.Sprintf("%T", *new(T))
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	name = strings.ToLower(name)
    return name
}


func NewCtrl[T models.IModel](db *gorm.DB, formHandler FormHandler[T]) *CrudCtrl[T] {

    res:= &CrudCtrl[T]{
		FormHandler: formHandler,
		Db:          db,
	}
    res.ModelName = res.baseName()
    return res
}

func (self *CrudCtrl[T]) DefineRoutes(r *gin.RouterGroup) *CrudCtrl[T] {
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

func (self *CrudCtrl[T]) List(c *gin.Context) {
	var items []T
	self.Db.Find(&items)
	rnd.HxRender(c,
		gin.H{fmt.Sprintf("%sList", self.ModelName): items},
		fmt.Sprintf("%s-list.html", self.ModelName))
}

func (self *CrudCtrl[T]) Details(c *gin.Context) {
	template := fmt.Sprintf("%s-details.html", self.ModelName)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		var item T
		self.Db.First(&item, id)
		rnd.HxRender(c, gin.H{self.ModelName: item}, template)
	} else {
		c.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
	}
}

func (self *CrudCtrl[T]) EditForm(c *gin.Context) {
	template := fmt.Sprintf("%s-edit.html", self.ModelName)
	idStr := c.Param("id")
	if idStr == "" {
		rnd.HxRender(c, gin.H{}, template)
		return
	} else {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			var item T
			self.Db.First(&item, idStr)
			rnd.HxRender(c, gin.H{self.ModelName: item}, template)
		} else {
			c.Error(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
		}
	}
}

func (self *CrudCtrl[T]) Upsert(c *gin.Context) {
	item, err := self.FormHandler(c)
	if err != nil {
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
	self.Db.Save(&item)
	c.Header("HX-Redirect", fmt.Sprintf("/%s/%d", self.Router.BasePath(), item.GetID()))
	c.String(http.StatusOK, "Saved")
	c.Abort()
}

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

