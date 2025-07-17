package main

import (
	"log"
	"net/http"
)

func main() {
	// Use the http.NewServerMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)                      //Restrict this route to exact matches on / only
	mux.HandleFunc("GET /snippet/view/{id}", snippetView) // Add the (id) wildcard segment
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	// Create the new route, which is restricted to POST request only,
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Print a log message to say that the server is starting
	log.Print("Starting server on :4000")

	// Use the http.listenAndServe() function to start a new web server. Web pass in
	// two parameters: the TCP network address to listen on (in this case "":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
