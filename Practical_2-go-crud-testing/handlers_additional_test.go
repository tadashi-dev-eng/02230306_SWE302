// handlers_additional_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetAllUsersHandler(t *testing.T) {
	resetState()

	// empty list
	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Get("/users", getAllUsersHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	var list []User
	if err := json.NewDecoder(rr.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d", len(list))
	}

	// add users and check
	users[1] = User{ID: 1, Name: "A"}
	users[2] = User{ID: 2, Name: "B"}

	req2, _ := http.NewRequest("GET", "/users", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr2.Code)
	}
	var list2 []User
	if err := json.NewDecoder(rr2.Body).Decode(&list2); err != nil {
		t.Fatal(err)
	}
	if len(list2) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list2))
	}
}

func TestCreateUserHandler_BadJSON(t *testing.T) {
	resetState()

	bad := "{\"name\": \"MissingEnd"
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(bad))
	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Post("/users", createUserHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json, got %d", rr.Code)
	}
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	resetState()
	req, _ := http.NewRequest("GET", "/users/abc", nil)
	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Get("/users/{id}", getUserHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id, got %d", rr.Code)
	}
}

func TestUpdateUserHandler_Scenarios(t *testing.T) {
	resetState()

	// Setup an existing user
	users[1] = User{ID: 1, Name: "Original"}

	// Successful update
	payload := `{"name":"Updated"}`
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(payload))
	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Put("/users/{id}", updateUserHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d", rr.Code)
	}
	var u User
	if err := json.NewDecoder(rr.Body).Decode(&u); err != nil {
		t.Fatal(err)
	}
	if u.ID != 1 || u.Name != "Updated" {
		t.Fatalf("unexpected updated user: %+v", u)
	}

	// Not found
	req2, _ := http.NewRequest("PUT", "/users/99", bytes.NewBufferString(payload))
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for update not found, got %d", rr2.Code)
	}

	// Bad JSON
	req3, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString("{badjson"))
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)
	if rr3.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad json on update, got %d", rr3.Code)
	}

	// Invalid ID
	req4, _ := http.NewRequest("PUT", "/users/abc", bytes.NewBufferString(payload))
	rr4 := httptest.NewRecorder()
	router.ServeHTTP(rr4, req4)
	if rr4.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id on update, got %d", rr4.Code)
	}
}

func TestDeleteUserHandler_Errors(t *testing.T) {
	resetState()

	// Not found
	req, _ := http.NewRequest("DELETE", "/users/99", nil)
	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Delete("/users/{id}", deleteUserHandler)
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when deleting missing user, got %d", rr.Code)
	}

	// Invalid id
	req2, _ := http.NewRequest("DELETE", "/users/abc", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id on delete, got %d", rr2.Code)
	}
}
