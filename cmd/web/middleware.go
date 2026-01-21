package main

import (
	"fmt"
	"net/http"
)

// commonHeaders sets security headers on all responses
// This middleware runs for every request
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers to protect against common web vulnerabilities
		w.Header().Set("Content-Security-Policy", "default-src 'self';style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff") // Prevent MIME type sniffing
		w.Header().Set("X-Frame-Options", "deny")           // Prevent clickjacking
		w.Header().Set("X-XSS-Protection", "0")             // Disable legacy XSS filter
		w.Header().Set("Server", "Go")                      // Set server header

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// logRequest logs details about each HTTP request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr       // Client IP address
			proto  = r.Proto            // HTTP protocol version
			method = r.Method           // HTTP method (GET, POST, etc.)
			uri    = r.URL.RequestURI() // Requested URI
		)
		// Log with structured logging
		app.logger.Info("Received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

// recoverPanic recovers from panics in handlers and returns a 500 error
// Prevents the server from crashing when a handler panics
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer function runs after the handler (even if it panics)
		defer func() {
			if err := recover(); err != nil {
				// Close connection after sending error response
				w.Header().Set("Connection", "close")
				// Convert panic value to error and log it
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
