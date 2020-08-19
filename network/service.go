package network

import (
	"Project_binary/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var err error

// GetData function to return Handler for get request
func GetData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		err = CheckAuth(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		var userData types.Data
		id := c.Param("id")
		db.Where("Id = ?", id).First(&userData)
		c.JSON(200, userData)
	}
}

// WriteData function to return Handler for POST request
func WriteData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err = CheckAuth(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		user := types.Data{ID: guuid.New(), Value: false, Key: ""}
		c.BindJSON(&user)
		db.Create(&user)
		c.JSON(200, user)
	}
}

// ModifyData function to return Handler for PATCH request
func ModifyData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err = CheckAuth(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		var userData types.Data
		var user types.Data
		id := c.Param("id")
		user.ID, err = guuid.Parse(id)
		if err != nil {
			log.Fatal(err)
		}
		c.BindJSON(&user)
		db.Model(&userData).Where("Id = ?", id).Update(map[string]interface{}{"Value": user.Value, "Key": user.Key})
		c.JSON(200, user)
	}
}

// DeleteData function to return Handler for DELETE request
func DeleteData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err = CheckAuth(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		var user types.Data
		id := c.Param("id")
		user.ID, err = guuid.Parse(id)
		// Check if data exists with given id
		db.Delete(&user)
		c.String(204, "No Content")
	}
}
