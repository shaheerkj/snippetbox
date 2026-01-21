package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shaheerkj/snippetbox/internal/models"
	"github.com/shaheerkj/snippetbox/internal/validator"
)

// snippetCreateForm holds form data and validation errors for snippet creation
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // "-" tells to ignore this field during decoding
	// removing the explicit fieldErrors struct field and instead
	// embedding the validator struct.
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

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate title field
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Cannot be more than 100 characters long.")

	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field can only be 1, 7 or 365")

	if !form.Valid() {
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
