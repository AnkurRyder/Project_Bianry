package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

var ID string

var host string = "http://127.0.0.1:8080/"

func TestMain(t *testing.T) {
	db, err := gorm.Open("mysql", "root:Stark9415@tcp(127.0.0.1:3306)/project_binary?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Errorf("DataBase Connection failed")
	}
	defer db.Close()
}

func TestPOST(t *testing.T) {
	db, err := gorm.Open("mysql", "root:Stark9415@tcp(127.0.0.1:3306)/project_binary?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Errorf("DataBase Connection failed")
	}
	defer db.Close()
	request := []byte("{\"value\": true, \"key\": \"hi there checking\"}")
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("POST request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := setupRouter(db)

	handler.ServeHTTP(rr, req)
	var tempData Data
	_ = json.Unmarshal([]byte(rr.Body.String()), &tempData)
	ID = tempData.ID.String()
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGet(t *testing.T) {
	db, err := gorm.Open("mysql", "root:Stark9415@tcp(127.0.0.1:3306)/project_binary?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Errorf("DataBase Connection failed")
	}
	defer db.Close()

	url := host + ID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("GET request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := setupRouter(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	s := fmt.Sprintf("{\"ID\":\"%s\",\"value\":true,\"key\":\"hi there checking\"}", ID)
	if strings.Compare(rr.Body.String(), s) != 0 {
		t.Errorf("Output not as expected, expected this %s got this %s", rr.Body.String(), s)
	}
}

func TestPatch(t *testing.T) {
	db, err := gorm.Open("mysql", "root:Stark9415@tcp(127.0.0.1:3306)/project_binary?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Errorf("DataBase Connection failed")
	}
	defer db.Close()
	request := []byte("{\"value\": true, \"key\": \"hi there checking \"}")

	url := host + ID

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("PATCH request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := setupRouter(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		fmt.Println(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDelete(t *testing.T) {
	db, err := gorm.Open("mysql", "root:Stark9415@tcp(127.0.0.1:3306)/project_binary?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Errorf("DataBase Connection failed")
	}
	defer db.Close()

	url := host + ID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Errorf("DELETE request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := setupRouter(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 204 {
		fmt.Println(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
