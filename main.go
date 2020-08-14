package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Data struct for storing data to the DBß
type Data struct {
	Id    uuid.UUID `json:”id”`
	Value bool      `json:"value"`
	Key   string    `json:"key"`
}

func main() {
	db, err := gorm.Open("mysql", "root:Stark9415@tcp(127.0.0.1:3306)/project_binary?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("hi therr", err)
	}
	defer db.Close()

	db.AutoMigrate(&Data{})

	router := gin.Default()

	router.GET("/:id", func(c *gin.Context) {
		var userData Data
		id := c.Param("id")
		db.Where("Id = ?", id).First(&userData)
		c.JSON(200, userData)
	})
	router.POST("/", func(c *gin.Context) {
		user := Data{Id: guuid.New(), Value: false, Key: ""}
		c.BindJSON(&user)
		db.Create(&user)
		c.JSON(200, user)
	})

	router.PATCH(":id", func(c *gin.Context) {
		var userData Data
		var user Data
		id := c.Param("id")
		user.Id, err = guuid.Parse(id)
		if err != nil {
			fmt.Println(err)
		}
		c.BindJSON(&user)
		db.Model(&userData).Where("Id = ?", id).Update(map[string]interface{}{"Value": user.Value, "Key": user.Key})
		c.JSON(200, user)
	})
	router.DELETE(":id", func(c *gin.Context) {
		var user Data
		id := c.Param("id")
		user.Id, err = guuid.Parse(id)
		// Check if data exists with given id
		db.Delete(&user)
		c.String(204, "No Content")
	})
	router.Run()

	// fmt.Println(db.NewRecord(user))
}
