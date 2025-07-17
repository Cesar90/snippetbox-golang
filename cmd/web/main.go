package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-file flag
	// This reads in the command-line flag value and assigns it to addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any erros are
	// encountered during parsing the application will be terminated
	flag.Parse()

	// Use the http.NewServerMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that path given to the http.Dir function is relative to the project
	// directy root
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.handle() function to register the file server as the handler for
	// all URL paths that start with "/static". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)                      //Restrict this route to exact matches on / only
	mux.HandleFunc("GET /snippet/view/{id}", snippetView) // Add the (id) wildcard segment
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	// Create the new route, which is restricted to POST request only,
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Print a log message to say that the server is starting
	// log.Print("Starting server on :4000")
	log.Printf("starting server on %s", *addr)

	// Use the http.listenAndServe() function to start a new web server. Web pass in
	// two parameters: the TCP network address to listen on (in this case "":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil
	// err := http.ListenAndServe(":4000", mux)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
