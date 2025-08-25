package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	// "text/template"

	"github.cesar90.com/internal/models"
	"github.cesar90.com/internal/validator"
)

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// FieldErrors map[string]string
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

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
	// w.Header().Add("Server", "Go")
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

	// app.render(w, r, http.StatusOK, "home.tmpl", templateData{
	// 	Snippets: snippets,
	// })
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)
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
	// app.render(w, r, http.StatusOK, "view.tmpl", templateData{
	// 	Snippet: snippet,
	// })

	// We no longer need to check for the flash message
	// flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// data.Flash = flash
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

// Add a snippetCreate handler function
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Display a form for creating a new snippet..."))
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// Add a snippetCreatePost handler function.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// title := "0 snail"
	// content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n"
	// expires := 7
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// }

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// form := snippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// 	// FieldErrors: map[string]string{},
	// }
	var form snippetCreateForm

	// err = app.formDecoder.Decode(&form, r.PostForm)
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Initialize a map to hold any validation errors for the form fields
	// fieldsErrors := make(map[string]string)

	// Check that the title value is not blank and is not more than 100
	// characteres long. If it fails of those checks, add a message to the
	// errors map using the field name as the key.
	// if strings.TrimSpace(title) == "" {
	// 	form.FieldErrors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(title) > 100 {
	// 	form.FieldErrors["title"] = "This field cannot be more than 100 characteres long"
	// }

	// // Check that the content value isn't blank
	// if strings.TrimSpace(content) == "" {
	// 	form.FieldErrors["content"] = "This field cannot be blank"
	// }

	// // Check the expires value matches one of the permitted values (1, 7 or)
	// // 365
	// if expires != 1 && expires != 7 && expires != 365 {
	// 	form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	// }

	// If there are any errors, dump them in a plain-text HTTP Response and
	// return from the handler
	// if len(form.FieldErrors) > 0 {
	// 	fmt.Fprint(w, form.FieldErrors)
	// }

	// If there are any validation errors, then the create.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the form Form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity.
	// when sending the response to indicate that was a validation error.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// if len(form.FieldErrors) > 0 {
	// 	data := app.newTemplateData(r)
	// 	data.Form = form
	// 	app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
	// 	return
	// }
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// id, err := app.snippets.Insert(title, content, expires)
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	// Use the Put() method to add a string value ("Snippet successfully")
	// created!) and the corresponding key("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// Use the w.WriteHeader() method to send a 201 status code.
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Save a new snippet..."))
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Display a form for signing up a new user...")
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Create a new user...")
	// Declare a zero-valued instance of our userSignupForm struct
	var form userSignupForm

	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characteres long")

	// If there are any errors, redisplay the signup form alogn with a 422
	// status code
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Otherwise send the placeholder response (for now).
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
