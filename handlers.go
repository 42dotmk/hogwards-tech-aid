package main

import (
	"app/db"
	"app/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)
func ListDonors(c *gin.Context) {
    var donors []models.Donor
    db.DB.Find(&donors)
    render(c, gin.H{"donors": donors}, "templates/donors.html")
}
func HandleDonor(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		if c.Param("mode") == "edit" {
			ShowDonorForm(c)
		} else {
			ShowDonorDetails(c)
		}
	case "POST":
		UpsertDonor(c)
	}
}

func ShowDonorDetails(c *gin.Context) {
	id := c.Param("id")
	var donor models.Donor
	db.DB.First(&donor, id)
	render(c, gin.H{"donor": donor}, "templates/donor.html")
}

func ShowDonorForm(c *gin.Context) {
	id := c.Param("id")
	var donor models.Donor
	db.DB.First(&donor, id)
	render(c, gin.H{"donor": donor}, "templates/edit-donor.html")
}

func UpsertDonor(c *gin.Context) {
	idStr := c.Param("id")

	var donor models.Donor
	donor.Name = c.PostForm("name")
	donor.Phone = c.PostForm("phone")
	donor.Email = c.PostForm("email")
	donor.Address = c.PostForm("address")
	donor.IsAnonymous = c.PostForm("isAnonymous") == "on"
	donor.IsLegalEntity = c.PostForm("isLegalEntity") == "on"
	donor.AnonymousName = c.PostForm("anonymousName")
	if idStr == "" {
		id := db.DB.Create(&donor)
		idStr = fmt.Sprintf("%d", id)
	} else {
		val, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			db.DB.Save(&donor)
		} else {
			msg := fmt.Sprintf("Invalid ID: %d", val)
			c.Redirect(http.StatusBadRequest, fmt.Sprintf("/error?msg=%s", msg))
			c.Abort()
			return
		}
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%d", idStr))
	c.Abort()
}

func Home(c *gin.Context) {
	render(c, gin.H{}, "home.html")
}

// func ShowList(c *gin.Context) {
// 	url := c.Param("url")
// 	var list models.List
// 	//var histories []models.CatChange
// 	//var deleteHistories []models.CatChange
// 	var copyAlert bool // DEFAULT -> FALSE
//
// 	db.DB.Preload("Categories.Items").Where("url = ?", url).First(&list)
//
// 	if list.ID != 0 {
//
// 		for i, cat := range list.Categories {
// 			db.DB.Where("category_id = ?", cat.ID).Order("updated_at desc").Offset(5).Find(&list.Categories[i].CatChanges)
// 			db.DB.Delete(&list.Categories[i].CatChanges) // NO ACUMULA MAS DE 5 EVENTOS
//
// 			db.DB.Where("category_id = ?", cat.ID).Order("updated_at desc").Find(&list.Categories[i].CatChanges)
// 		}
//
// 		if true {
// 			copyAlert = true
// 		}
//
// 		render(c, gin.H{"title": "Grocery List", "list": list, "copyAlert": copyAlert}, "list.html")
//
// 	}
//
// }
//
// func Category(c *gin.Context) {
//
// 	url := c.Param("url")
//
// 	switch c.Request.Method {
//
// 	case "POST":
// 		{
// 			name := c.PostForm("name")
// 			newCat := models.Category{
// 				Name: name,
// 				Url:  url,
// 			}
// 			db.DB.Create(&newCat)
//
// 			c.HTML(http.StatusOK, "category.html", newCat)
// 		}
//
// 	case "DELETE":
// 		{
// 			id := c.Param("id")
// 			var category models.Category
// 			db.DB.Delete(&category, id)
// 		}
//
// 	case "PUT":
// 		{
// 			id := c.Param("id")
// 			name := c.PostForm("name")
//
// 			var category models.Category
// 			db.DB.First(&category, id)
//
// 			oldname := category.Name
// 			category.Name = name
// 			db.DB.Save(&category)
//
// 			// -------------  ADD EVENT TO HISTORY --------------
//
// 			var catChange models.CatChange
// 			catChange.Title = "'" + oldname + "'" + " changed to '" + name + "'"
// 			catChange.Url = url
// 			catChange.TypeChange = 1
// 			catChange.CategoryID = category.ID
// 			db.DB.Create(&catChange)
//
// 			// --------------------------------------------------
//
// 			c.HTML(http.StatusOK, "history.html", catChange)
// 		}
//
// 	}
//
// }
//
// func Item(c *gin.Context) {
// 	url := c.Param("url")
//
// 	switch c.Request.Method {
// 	case "POST":
// 		{
// 			u64, _ := strconv.ParseUint(c.Param("idcat"), 10, 64)
// 			id := uint(u64)
// 			name := c.PostForm("name")
// 			item := models.Item{
// 				Name:       name,
// 				CategoryID: id,
// 				Url:        url,
// 			}
// 			db.DB.Create(&item)
//
// 			c.HTML(http.StatusOK, "item.html", item)
// 		}
// 	case "DELETE":
// 		{
// 			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
// 			id := uint(u64)
// 			var item models.Item
// 			db.DB.Where("id = ?", id).First(&item)
// 			db.DB.Unscoped().Delete(&item)
//
// 			// -------------  ADD EVENT TO HISTORY --------------
// 			var catChange models.CatChange
// 			catChange.Title = item.Name
// 			catChange.Url = url
// 			catChange.TypeChange = 2
// 			catChange.CategoryID = item.CategoryID
// 			db.DB.Create(&catChange)
// 			// --------------------------------------------------
// 			c.HTML(http.StatusOK, "history.html", catChange)
// 		}
// 	case "PATCH":
// 		{
// 			u64, _ := strconv.ParseUint(c.Param("iditem"), 10, 64)
// 			id := uint(u64)
//
// 			var item models.Item
// 			db.DB.Where("id = ?", id).First(&item)
// 			db.DB.Unscoped().Delete(&item)
//
// 			// -------------  ADD EVENT TO HISTORY ---------------
// 			var catChange models.CatChange
// 			catChange.Title = item.Name
// 			catChange.Url = url
// 			catChange.TypeChange = 3
// 			catChange.CategoryID = item.CategoryID
// 			db.DB.Create(&catChange)
// 			// --------------------------------------------------
//
// 			c.HTML(http.StatusOK, "history.html", catChange)
// 		}
// 	}
// }
