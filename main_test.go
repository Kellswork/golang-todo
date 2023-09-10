// import package main
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// create Todo test function
func TestCreateTodo(t *testing.T) {
	// thinsg to test
	// 1. that the response body was decoded successfully
	// 2. that the todo title whne empty returns a response
	// 3. a todo was succesfully added to the DB
	// generally, tests are for testing the response bodies
	// any method that accepts responsewriter also accepts response recorder
	// in the handlerFunc, I was confused if we are passing in the type or the function

	// user input is a json string
	jsonStr := []byte(`{"title": "go to the beach"}`)

	req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(createTodo)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("createTodo returned wrong status code. Got: %v Expected %v", status, http.StatusCreated)
	}

	// expected should be what you set the response body to return in the functiuon handler
	expected := `{"message": "Todo created successfully", ID":""}`

	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v expected: %v", recorder.Body.String(), expected)
	}

	// to test if this is working, make sure to have your TP up and running.
	// I stopped at the part were the test was failing due to a different error response it sgetting
	// when you get back to this, run the code to see the failing test result
	// check chat gpt for the sample testing code it created
}
