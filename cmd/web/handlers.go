package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler function writes a byte slice containing
// "Hello from Snippetbox" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	// Set a new cache-control header. If an existing "Cache-Control" header exists
	// It will be overwritten
	// w.Header().Set("Cache-Control", "public, max-age=31536000")

	// In contrast, the Add() method appends a new "Cache-Control" header and can
	// be called multiple times.
	// w.Header().Add("Cache-Control", "public")
	// w.Header().Add("Cache-Control", "max-age=31536000")

	// Delete all values for the "Cache-Control" header.
	// w.Header().Del("Cache-Control")

	// Retrieve the first value for the "Cache-Control" header
	// w.Header().Get("Cache-Control")

	//Retrieve a slice of all values for the "Cache-Control"
	// w.Header().Values("Cache-Control")

	// Use the Header().Add() method to add a 'Server: Go' header to the
	// response header map. The first parameter is the header name, and
	// the second parameter is the header value.
	w.Header().Add("Server", "Go")

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to a integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.SprintF() function to interpolate the id value with a
	// message, then write it as the HTTP response.
	// msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	// w.Write([]byte(msg))
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Use the w.WriteHeader() method to send a 201 status code.
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new snippet..."))
}
