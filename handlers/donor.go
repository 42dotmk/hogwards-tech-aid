package handlers

import (
	"github.com/42dotmk/hogwards/models"
	rnd"github.com/42dotmk/hogwards/lib/renderers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DonorController struct {
	r  gin.IRouter
	db *gorm.DB
}

func NewDonorCtrl( db *gorm.DB) *DonorController {
	return &DonorController{db: db}
}

func (self *DonorController) OnRouter(r gin.IRouter) {
    self.r = r
	self.r.GET("/", self.ListDonors)
	self.r.GET("/new", self.ShowDonorForm)
	self.r.GET("/:id", self.ShowDonorDetails)
	self.r.GET("/:id/edit", self.ShowDonorForm)
	self.r.POST("/:id", self.UpsertDonor)
	self.r.PUT("/", self.UpsertDonor)
	self.r.DELETE("/:id", self.DeleteDonor)
}

func (self *DonorController) ListDonors(c *gin.Context) {
	var donors []models.Donor
    self.db.Find(&donors)

    // c.HTML(http.StatusOK, "donors-list.html", gin.H{"Donors": donors})
	rnd.HxRender(c, gin.H{"Donors": donors}, "donors-list.html")
}

func (self *DonorController) ShowDonorDetails(c *gin.Context) {
	if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
		var donor models.Donor
		self.db.First(&donor, id)
		rnd.HxRender(c, gin.H{"Donor": donor}, "donor.html")
	} else {
		c.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID: %d", id))
	}
}

func (self *DonorController) ShowDonorForm(c *gin.Context) {
	if idStr := c.Param("id"); idStr == "" {
		rnd.HxRender(c, gin.H{}, "edit-donor.html")
		return
	} else if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		var donor models.Donor
		self.db.First(&donor, idStr)
		rnd.HxRender(c, gin.H{"Donor": donor}, "edit-donor.html")
	} else {
		c.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID: %d", id))
	}
}

func (self *DonorController) UpsertDonor(c *gin.Context) {
	var donor models.Donor
	donor.Name = c.PostForm("Name")
	donor.Phone = c.PostForm("Phone")
	donor.Email = c.PostForm("Email")
	donor.Address = c.PostForm("Address")
	donor.IsAnonymous = c.PostForm("IsAnonymous") == "on"
	donor.IsLegalEntity = c.PostForm("IsLegalEntity") == "on"
	donor.AnonymousName = c.PostForm("AnonymousName")

	if idStr := c.Param("id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID: %d", id)
			c.String(http.StatusBadRequest, msg)
			c.Abort()
			return
		}
		donor.ID = uint(id)
	}
	self.db.Save(&donor)
	c.Header("HX-Redirect", fmt.Sprintf("/donors/%d", donor.ID))
	c.String(http.StatusOK, "Saved")
	c.Abort()
}

func (self *DonorController) DeleteDonor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid ID: %d", id)
		c.String(http.StatusBadRequest, msg)
		c.Abort()
		return
	}
	self.db.Delete(&models.Donor{}, id)
	c.Header("HX-Redirect", "/donors")
	c.String(http.StatusOK, "Deleted")
	c.Abort()
}
