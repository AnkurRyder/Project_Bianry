package network

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SetupRouter function defines the controller
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.POST("/signup", SignUp(db))

	router.POST("/login", Login(db))

	authorized := router.Group("/")

	authorized.Use(CheckAuth(db))
	{
		authorized.POST("/logout", Logout(db))

		authorized.GET(":id", GetData(db))

		authorized.POST("/", WriteData(db))

		authorized.PATCH(":id", ModifyData(db))

		authorized.DELETE(":id", DeleteData(db))
	}

	return router
}
