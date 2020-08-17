package db

import (
	"Project_binary/types"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

// Connection function to return gorm db pointer
func Connection(dbTemp string) *gorm.DB {
	dbString := getDBString(dbTemp)
	db, err := gorm.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&(types.Data{}))
	db.AutoMigrate(&(types.User{}))
	db.AutoMigrate(&(types.TokenMeta{}))
	return db
}

func getDBString(dbConStringMain string) string {
	dbName := GoDotEnvVariable("DB_NAME")
	password := GoDotEnvVariable("Password")
	user := GoDotEnvVariable("user")
	return fmt.Sprintf(dbConStringMain, user, password, dbName)
}

// GoDotEnvVariable Return value from env file
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
