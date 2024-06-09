package main

import (
	"app/db"
	"app/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BuildDonorRoutes(r *gin.RouterGroup) {
	r.GET("/", ListDonors)
	r.GET("/new", ShowDonorForm)
	r.GET("/:id", ShowDonorDetails)
	r.GET("/:id/edit", ShowDonorForm)
	r.POST("/:id", UpsertDonor)
	r.PUT("/", UpsertDonor)
	r.DELETE("/:id", DeleteDonor)
}
func ListDonors(c *gin.Context) {
	var donors []models.Donor
	db.DB.Find(&donors)
	hxRender(c, gin.H{"Donors": donors}, "donors-list.html")
}

func ShowDonorDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		var donor models.Donor
		db.DB.First(&donor, id)
		hxRender(c, gin.H{"Donor": donor}, "donor.html")
	} else {
		c.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID: %d", id))
	}
}

func ShowDonorForm(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		hxRender(c, gin.H{}, "edit-donor.html")
		return
	} else {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			var donor models.Donor
			db.DB.First(&donor, idStr)
			hxRender(c, gin.H{"Donor": donor}, "edit-donor.html")
		} else {
			c.Error(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID: %d", id))
		}
	}
}

func UpsertDonor(c *gin.Context) {
	var donor models.Donor
	donor.Name = c.PostForm("Name")
	donor.Phone = c.PostForm("Phone")
	donor.Email = c.PostForm("Email")
	donor.Address = c.PostForm("Address")
	donor.IsAnonymous = c.PostForm("IsAnonymous") == "on"
	donor.IsLegalEntity = c.PostForm("IsLegalEntity") == "on"
	donor.AnonymousName = c.PostForm("AnonymousName")

	idStr := c.Param("id")
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID: %d", id)
			c.String(http.StatusBadRequest, msg)
			c.Abort()
			return
		}
		donor.ID = uint(id)
	}
	db.DB.Save(&donor)
	c.Header("HX-Redirect", fmt.Sprintf("/donors/%d", donor.ID))
	c.String(http.StatusOK, "Saved")
	c.Abort()
}

func DeleteDonor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid ID: %d", id)
		c.String(http.StatusBadRequest, msg)
		c.Abort()
		return
	}
	db.DB.Delete(&models.Donor{}, id)
	c.Header("HX-Redirect", "/donors")
	c.String(http.StatusOK, "Deleted")
	c.Abort()
}
