package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// newTemplateData creates a templateData struct populated with common data
// that's needed across all templates (like current year for footer)
// Takes *http.Request as parameter for future expansion (sessions, auth, etc.)
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(), // Used in footer copyright
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// serverError logs the error with request details and sends a 500 response to the user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	// Log the error with structured logging (includes method and URI for debugging)
	app.logger.Error(err.Error(), "method", method, "uri", uri)

	// Send generic 500 error response to user (don't leak error details)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific HTTP status code and error message to the user
// Used for 4xx errors (bad request, not found, etc.)
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// render executes a template from the cache and writes the response
// Uses a buffer to catch template execution errors before writing to the client
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve template from cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Execute template into buffer first (not directly to ResponseWriter)
	// This way if there's an error, we haven't sent partial response yet
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If template executed successfully, write status and content
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// Helper to check for errors in the form parsing, e.g if we pass a
// non-nil dst pointer to the form.Decode() function, it will cause a different kind
// of error that can't be classified simply as a 4XX/400
func (app *application) decodePostForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return err
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
