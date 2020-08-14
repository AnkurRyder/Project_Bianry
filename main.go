package main

import (
	"fmt"

	"github.com/google/uuid"
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

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

	user := Data{Id: guuid.New(), Value: false, Key: "hi there"}

	// fmt.Println(db.NewRecord(user)) // => returns `true` as primary key is blank

	db.Create(&user)

	fmt.Println(db.NewRecord(user))
}
