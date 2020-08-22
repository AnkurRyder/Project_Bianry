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
		resp, err := GetDataHelper(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(200, resp)
	}
}

// WriteData function to return Handler for POST request
func writeData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := WriteDataHelper(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(200, resp)
	}
}

// ModifyData function to return Handler for PATCH request
func modifyData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := ModifyDataHelper(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(200, resp)
	}
}

// DeleteData function to return Handler for DELETE request
func deleteData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := ModifyDataHelper(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(204, resp)
	}
}

// Login function for helping user to login
func login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := LoginHelper(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// SignUp function for user dignup
func signUp(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := SignUpHelper(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusCreated, resp)
	}
}

// Logout function for user to logout
func logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := LogoutHelepr(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
