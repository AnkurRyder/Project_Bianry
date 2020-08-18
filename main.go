package main

import (
	"Project_binary/db"
	"Project_binary/network"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var err error
var dbConStringMain string = "%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local"

func main() {
	db := db.Connection(dbConStringMain)

	defer db.Close()

	router := setupRouter(db)

	router.Run()
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/login", network.Login(db))

	router.GET("/:id", network.GetData(db))

	router.POST("/", network.WriteData(db))

	router.PATCH(":id", network.ModifyData(db))

	router.DELETE(":id", network.DeleteData(db))

	return router
}
