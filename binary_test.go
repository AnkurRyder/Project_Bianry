package main

import (
	"Project_binary/db"
	"Project_binary/network"
	"Project_binary/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

var ID string
var tokenString string
var dbConString string = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
var host string = "http://%s:%s/"
var dbtest *gorm.DB

func TestMain(m *testing.M) {
	dbtest = db.Connection(dbConString)
	host = fmt.Sprintf(host, db.GoDotEnvVariable("HOST"), db.GoDotEnvVariable("PORT"))
	defer dbtest.Close()
	os.Exit(m.Run())
}

func TestSignup(t *testing.T) {
	request := []byte("{\"username\": \"test_me\", \"password\": \"pass\"}")
	req, err := http.NewRequest("POST", host+"signup", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("POST request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := network.SetupRouter(dbtest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestLogin(t *testing.T) {
	request := []byte("{\"username\": \"test\", \"password\": \"pass\"}")
	req, err := http.NewRequest("POST", host+"login", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("POST request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := network.SetupRouter(dbtest)

	handler.ServeHTTP(rr, req)

	tokenString = rr.Body.String()
	tokenString = tokenString[1 : len(tokenString)-1]

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestPOST(t *testing.T) {
	request := []byte("{\"value\": true, \"key\": \"hi there checking\"}")
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(request))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	if err != nil {
		t.Errorf("POST request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := network.SetupRouter(dbtest)

	handler.ServeHTTP(rr, req)
	var tempData types.Data
	_ = json.Unmarshal([]byte(rr.Body.String()), &tempData)
	ID = tempData.ID.String()
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGet(t *testing.T) {
	url := host + ID

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	if err != nil {
		t.Errorf("GET request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := network.SetupRouter(dbtest)

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
	request := []byte("{\"value\": true,\"key\": \"hi there checking update\"}")

	url := host + ID

	// NewRequest returns a new incoming server
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(request))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	if err != nil {
		t.Errorf("PATCH request failed")
	}
	rr := httptest.NewRecorder()

	handler := network.SetupRouter(dbtest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		fmt.Println(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	s := fmt.Sprintf("{\"ID\":\"%s\",\"value\":true,\"key\":\"hi there checking update\"}", ID)
	if strings.Compare(rr.Body.String(), s) != 0 {
		t.Errorf("Output not as expected, expected this %s got this %s", rr.Body.String(), s)
	}
}

func TestDelete(t *testing.T) {

	url := host + ID

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	if err != nil {
		t.Errorf("DELETE request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := network.SetupRouter(dbtest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != 204 {
		fmt.Println(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, 204)
	}
}
