package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/shaheerkj/snippetbox/internal/models"
)

// templateData holds dynamic data that's passed to HTML templates
// Provides a consistent structure for all template data
type templateData struct {
	Snippet     models.Snippet   // Single snippet (for view page)
	Snippets    []models.Snippet // Multiple snippets (for home page)
	CurrentYear int              // Current year for footer
	Form        any              // Form data and validation errors
	Flash       string
}

// humanDate formats a time.Time into a human-readable string
// Uses Go's reference time format: Mon Jan 2 15:04:05 MST 2006
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// functions is a map of custom functions available in templates
// Must be registered with template before parsing
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache parses all templates at application startup and caches them
// This improves performance by avoiding re-parsing templates on every request
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize empty cache map
	cache := map[string]*template.Template{}

	// Get all page templates from ./ui/html/pages/
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through each page template
	for _, page := range pages {
		// Extract filename to use as cache key (e.g., "home.html")
		name := filepath.Base(page)

		// Parse template in 3 stages:
		// 1. Create new template set, register custom functions, parse base layout
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// 2. Parse all partials (nav, footer, etc.)
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// 3. Parse the specific page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add parsed template set to cache
		cache[name] = ts
	}

	return cache, nil
}
