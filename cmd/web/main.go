package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver (blank import to register driver)
	"github.com/shaheerkj/snippetbox/internal/models"
)

// application holds application-wide dependencies and shared resources
// These are injected into handlers via receiver methods
type application struct {
	logger         *slog.Logger                  // Structured logger for application-wide logging
	snippets       *models.SnippetModel          // Database model for snippet operations
	templateCache  map[string]*template.Template // Pre-parsed templates for better performance
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Parse command-line flags for configuration
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "shaheer:110434@/snippetbox?parseTime=true", "MySQL DSN String")

	flag.Parse() // Parse the flags from command line

	// Initialize structured logger that writes to stdout
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open database connection and verify connectivity
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close() // Ensure database connection is closed when main() returns

	// Initialize template cache (pre-parse all templates at startup)
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
	}

	// initialize the form decoder
	formDecoder := form.NewDecoder()

	//session config
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize application dependencies
	// Using & creates a pointer, allowing the struct to be shared across handlers
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db}, // Initialize snippet model with DB connection
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("Starting server", "addr", *addr)

	// Start HTTP server with configured routes
	// This blocks until the server encounters an error
	// err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB creates a database connection pool and verifies connectivity
func openDB(dsn string) (*sql.DB, error) {
	// Create connection pool (doesn't actually connect yet)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection is actually working
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
