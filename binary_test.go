package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestPOST(t *testing.T) {
// 	// request, err := json.Marshal(map[string]string{
// 	// 	"value": true,
// 	// 	"key": "hi there checking POST",
// 	// })
// 	request := []byte("{\"value\": true, \"key\": \"hi there checking \"}")
// 	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/", bytes.NewBuffer(request))
// 	if err != nil {
// 		t.Errorf("POST request failed")
// 	}
// 	// NewRequest returns a new incoming server
// 	rr := httptest.NewRecorder()

// 	handler := setupRouter()

// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// }

// '{"value": true, "key": "hi there checking POST"}'

func TestGet(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/1114c5af-ddec-40f2-92c0-2d887ddc6cb2", nil)
	if err != nil {
		t.Errorf("POST request failed")
	}
	// NewRequest returns a new incoming server
	rr := httptest.NewRecorder()

	handler := setupRouter()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		fmt.Println(rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
