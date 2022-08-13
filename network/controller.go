package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SetupRouter function defines the controller
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.POST("/signup", signUp(db))

	router.POST("/login", login(db))

	authorized := router.Group("/")

	authorized.Use(CheckAuth(db))
	{
		authorized.POST("/logout", logout(db))

		authorized.GET(":id", getData(db))

		authorized.POST("/", writeData(db))

		authorized.PATCH(":id", modifyData(db))

		authorized.DELETE(":id", deleteData(db))
	}

	return router
}

// GetData function to return Handler for get request
func getData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := GetDataHelper(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(200, resp)
	}
}

// WriteData function to return Handler for POST request
func writeData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := WriteDataHelper(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(200, resp)
	}
}

// ModifyData function to return Handler for PATCH request
func modifyData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := ModifyDataHelper(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(200, resp)
	}
}

// DeleteData function to return Handler for DELETE request
func deleteData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := ModifyDataHelper(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(204, resp)
	}
}

// Login function for helping user to login
func login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := LoginHelper(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// SignUp function for user dignup
func signUp(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := SignUpHelper(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}

// Logout function for user to logout
func logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, errResp := LogoutHelepr(c, db)
		if errResp.Err != nil {
			c.JSON(errResp.HTTPCode, errResp.Err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
