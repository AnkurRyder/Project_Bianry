package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

// Data struct for storing data to the DBß
type Data struct {
	ID    uuid.UUID
	Value bool   `json:"value"`
	Key   string `json:"key"`
}

// User struct for storing user credentials to DB
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var err error
var dbConStringMain string = "%s:%s@tcp(docker.for.mac.localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local"

func main() {
	db := dbConnection(dbConStringMain)

	defer db.Close()

	router := setupRouter(db)

	router.Run()
}

func dbConnection(dbTemp string) *gorm.DB {
	dbString := getDBString(dbTemp)
	db, err := gorm.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Data{})
	db.AutoMigrate(&User{})
	return db
}

func getDBString(dbConStringMain string) string {
	dbName := goDotEnvVariable("DB_NAME")
	password := goDotEnvVariable("Password")
	user := goDotEnvVariable("user")
	return fmt.Sprintf(dbConStringMain, user, password, dbName)
}

// Return value from env file
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/:id", getData(db))

	router.POST("/", writeData(db))

	router.PATCH(":id", modifyData(db))

	router.DELETE(":id", deleteData(db))

	return router
}

func getData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userData Data
		id := c.Param("id")
		db.Where("Id = ?", id).First(&userData)
		c.JSON(200, userData)
	}
}

func writeData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := Data{ID: guuid.New(), Value: false, Key: ""}
		c.BindJSON(&user)
		db.Create(&user)
		c.JSON(200, user)
	}
}

func modifyData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userData Data
		var user Data
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

func deleteData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user Data
		id := c.Param("id")
		user.ID, err = guuid.Parse(id)
		// Check if data exists with given id
		db.Delete(&user)
		c.String(204, "No Content")
	}
}
