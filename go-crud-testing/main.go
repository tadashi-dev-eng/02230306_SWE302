// main.go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Setup routes
	r.Get("/users", getAllUsersHandler)
	r.Post("/users", createUserHandler)
	r.Get("/users/{id}", getUserHandler)
	r.Put("/users/{id}", updateUserHandler)
	r.Delete("/users/{id}", deleteUserHandler)

	log.Println("Server starting on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
