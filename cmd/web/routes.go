package main

import (
	"net/http"

	"github.com/justinas/alice" // Middleware chaining library
)

// routes sets up the application's HTTP routes and middleware chain
func (app *application) routes() http.Handler {
	// Create a new router/mux
	mux := http.NewServeMux()

	// Serve static files (CSS, JS, images) from ./ui/static/
	// StripPrefix removes "/static" from the URL before looking up the file
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	//creating a new middleware chain containing the middleware specific to our
	//dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Application routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))                          // Homepage (exact match only)
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))     // View individual snippet
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))      // Display create form
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost)) // Process form submission

	// Create middleware chain (executed in order: recoverPanic -> logRequest -> commonHeaders)
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Wrap the mux with the middleware chain
	return standard.Then(mux)
}
