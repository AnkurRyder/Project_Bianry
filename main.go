package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type data struct {
	id    string
	value bool
	key   string
}

func main() {
	var obj1 data
	obj1.id = "sadnsjkds"
	obj1.value = true
	obj1.key = "name"

	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	// Creates a gin router with default middleware
	router := gin.Default()

	// A handler for GET request on /example
	router.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"id":    obj1.id,
			"value": obj1.value,
			"key":   obj1.key,
		}) // gin.H is a shortcut for map[string]interface{}
	})
	router.Run() // listen and serve on port 8080
}
