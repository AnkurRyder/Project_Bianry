package db

import (
	"Project_binary/types"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func DbConnection(dbTemp string) *gorm.DB {
	dbString := getDBString(dbTemp)
	db, err := gorm.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&(types.Data{}))
	db.AutoMigrate(&(types.User{}))
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
