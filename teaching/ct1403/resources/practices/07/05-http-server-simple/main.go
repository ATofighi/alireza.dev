package main

import (
	"fmt"
	"log"
	"net/http"
)

// homeHandler handles GET requests to "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Setting a header in the response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Welcome to the Home Page!")
}

// formHandler handles POST requests to "/submit"
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data (x-www-form-urlencoded or multipart/form-data)
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		message := r.FormValue("message")

		w.WriteHeader(http.StatusOK)
		// Construct a simple response
		fmt.Fprintf(w, "Form submission received!\n")
		fmt.Fprintf(w, "Name: %s\n", name)
		fmt.Fprintf(w, "Message: %s\n", message)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Register multiple handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/submit", formHandler)

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
