package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/shaheerkj/snippetbox/internal/models"
)

// snippetCreateForm holds form data and validation errors for snippet creation
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string // Maps field names to error messages
}

// home displays the homepage with the latest snippets
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Fetch the 10 most recent non-expired snippets from database
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Prepare template data with snippets and common data (current year, etc.)
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Render the home template with 200 OK status
	app.render(w, r, http.StatusOK, "home.html", data)
}

// snippetView displays a specific snippet by ID
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the id parameter from the URL path and convert to integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		// ID is not a valid positive integer
		http.NotFound(w, r)
		return
	}

	// Fetch the snippet from the database
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// Snippet doesn't exist or has expired
			http.NotFound(w, r)
		} else {
			// Database error
			app.serverError(w, r, err)
		}
		return
	}

	// Prepare template data and render the view
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.html", data)
}

// snippetCreate displays the form for creating a new snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Prepare template data with default form values
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365, // Default to 365 days
	}
	app.render(w, r, http.StatusOK, "create.html", data)
}

// snippetCreatePost handles the POST request to create a new snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request body
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Convert expires field from string to integer
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Extract form values into struct
	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{}, // Initialize empty error map
	}

	// Validate title field
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		// Use RuneCountInString to correctly count multi-byte characters
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long."
	}

	// Validate content field
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Validate expires field - must be 1, 7, or 365 days
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365."
	}

	// If there are validation errors, re-display the form
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
