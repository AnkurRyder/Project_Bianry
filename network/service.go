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
		var userData types.Data
		id := c.Param("id")
		db.Where("Id = ?", id).First(&userData)
		c.JSON(200, userData)
	}
}

// WriteData function to return Handler for POST request
func WriteData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := types.Data{ID: guuid.New(), Value: false, Key: ""}
		c.BindJSON(&user)
		db.Create(&user)
		c.JSON(200, user)
	}
}

// ModifyData function to return Handler for PATCH request
func ModifyData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		var user types.Data
		id := c.Param("id")
		idUUID, err := guuid.Parse(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Wrong ID")
			return
		}
		db.Where("Id = ?", idUUID).First(&user)
		if user.ID != idUUID {
			c.String(http.StatusBadRequest, "Wrong ID or no data exists")
			return
		}
		db.Delete(&user)
		c.JSON(204, "No Content")
	}
}

// Login function for helping user to login
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u types.User
		var user types.User
		err := c.ShouldBindJSON(&u)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		db.Where("username = ?", u.Username).First(&user)
		//compare the user from the request, with the one we defined:
		if user.Username != u.Username || user.Password != u.Password {
			c.JSON(http.StatusUnauthorized, "Please provide valid login details")
			return
		}
		token, err := CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		CreateAuth(user.ID, token, db)
		c.JSON(http.StatusOK, token.AccessToken)
	}
}

// SignUp function for user dignup
func SignUp(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userData types.User
		var userDbCheck types.User
		err := c.ShouldBindJSON(&userData)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
			return
		}
		db.Where("username = ?", userData.Username).First(&userDbCheck)
		if userDbCheck.Username == userData.Username {
			c.JSON(http.StatusForbidden, "Username already exists")
			return
		}
		db.Create(&userData)
		c.JSON(http.StatusCreated, "Account Created")
	}
}

// Logout function for user to logout
func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		deleted, err := DeleteAuth(au.AccessUuid, db)
		if deleted == 0 || err != nil { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		c.JSON(http.StatusOK, "Successfully logged out")
	}
}
