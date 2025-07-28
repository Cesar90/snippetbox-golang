package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	// "text/template"

	"github.cesar90.com/internal/models"
)

// Define a home handler function writes a byte slice containing
// "Hello from Snippetbox" as the response body

// Change the signature of the handler so it is defined as a method against
// *application

func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
	// w.Write([]byte("Hello from Snippetbox"))

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// the http.Error() function to send and internal Server Error response to the
	// user, and then return the handler so no subsequent code is executed
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// log.Print(err.Error())
	// 	// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// 	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// Create an instance of a templateData struct holding the slice of
	// snippets
	// data := templateData{
	// 	Snippets: snippets,
	// }

	// err = ts.Execute(w, nil)
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	// log.Print(err.Error())
	// 	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// }

	// Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil
	app.render(w, r, http.StatusOK, "home.tmpl", templateData{
		Snippets: snippets,
	})
}

// Add a snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
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
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// log.Print(err.Error())
	// 	// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// 	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// Create an instance of a templateData struct holding the snippet data.
	// data := templateData{
	// 	Snippet: snippet,
	// }

	// err = ts.Execute(w, nil)
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	// log.Print(err.Error())
	// 	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	app.serverError(w, r, err)
	// }

	// fmt.Fprintf(w, "%+v", snippet)
	app.render(w, r, http.StatusOK, "view.tmpl", templateData{
		Snippet: snippet,
	})
}

// Add a snippetCreate handler function
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// Use the w.WriteHeader() method to send a 201 status code.
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Save a new snippet..."))
}
