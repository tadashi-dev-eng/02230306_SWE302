// handlers_test.go
package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

// A helper function to reset the state before each test
func resetState() {
	users = make(map[int]User)
	nextID = 1
}

func TestCreateUserHandler(t *testing.T) {
	resetState()

	// 1. Define the user we want to create
	userPayload := `{"name": "John Doe"}`
	req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(userPayload))
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// 3. Create a router and serve the request
	router := chi.NewRouter()
	router.Post("/users", createUserHandler)
	router.ServeHTTP(rr, req)

	// 4. Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 5. Check the response body
	var createdUser User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatal(err)
	}

	if createdUser.Name != "John Doe" {
		t.Errorf("handler returned unexpected body: got name %v want %v", createdUser.Name, "John Doe")
	}

	if createdUser.ID != 1 {
		t.Errorf("handler returned unexpected body: got id %v want %v", createdUser.ID, 1)
	}
}

func TestGetUserHandler(t *testing.T) {
	resetState()

	// First, create a user to fetch
	users[1] = User{ID: 1, Name: "Jane Doe"}

	// Test case 1: User found
	t.Run("User Found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/users/{id}", getUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var foundUser User
		json.NewDecoder(rr.Body).Decode(&foundUser)
		if foundUser.Name != "Jane Doe" {
			t.Errorf("handler returned unexpected body: got %v want %v", foundUser.Name, "Jane Doe")
		}
	})

	// Test case 2: User not found
	t.Run("User Not Found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/99", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/users/{id}", getUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

// You would continue this pattern for updateUserHandler and deleteUserHandler...

func TestDeleteUserHandler(t *testing.T) {
	resetState()
	users[1] = User{ID: 1, Name: "To Be Deleted"}

	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{id}", deleteUserHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify the user was actually deleted
	if _, ok := users[1]; ok {
		t.Error("user was not deleted from the map")
	}
}
