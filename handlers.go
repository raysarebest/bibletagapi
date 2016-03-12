package main

import (
	"fmt"
	"net/http"
)

// Index is the handler for the root URL
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to this example API\n")
}
