package main

import (
	"Project_binary/db"
	"Project_binary/network"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var err error

func main() {
	time.Sleep(5 * time.Second)

	db := db.Connection(db.DBConStringMain)

	defer db.Close()

	router := network.SetupRouter(db)

	router.Run(":8080")
}
