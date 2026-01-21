# Snippetbox

A web application for sharing code snippets, built with Go. Created while following Alex Edwards' book *Let's Go*.

## Features

- **Create snippets** - Share code snippets with configurable expiration (1 day, 7 days, or 1 year)
- **View snippets** - Browse and view individual code snippets
- **Auto-expiration** - Snippets automatically expire and are hidden after their set duration
- **Template caching** - Pre-parsed templates for better performance
- **Middleware chain** - Request logging, panic recovery, and security headers

## Tech Stack

- **Go 1.21+** - Backend language
- **MySQL** - Database
- **html/template** - Server-side templating
- **Alice** - Middleware chaining

## Project Structure

```
├── cmd/web/           # Application entry point and handlers
│   ├── main.go        # App initialization and server startup
│   ├── handlers.go    # HTTP handlers for routes
│   ├── helpers.go     # Helper functions (error handling, rendering)
│   ├── routes.go      # Route definitions and middleware setup
│   ├── middleware.go  # Custom middleware (logging, security headers)
│   └── templates.go   # Template caching and custom functions
├── internal/models/   # Database models
│   ├── snippets.go    # Snippet CRUD operations
│   └── errors.go      # Custom error types
└── ui/                # Frontend assets
    ├── html/          # Go templates
    └── static/        # CSS, JS, images
```

## Running the Application

```bash
# Start the server (default port 4000)
go run ./cmd/web

# With custom port
go run ./cmd/web -addr=":8080"

# With custom database DSN
go run ./cmd/web -dsn="user:pass@/snippetbox?parseTime=true"
```

## Database Setup

```sql
CREATE DATABASE snippetbox;

CREATE TABLE snippets (
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_created (created)
);
```

## Credits

Built following [Let's Go](https://lets-go.alexedwards.net/) by Alex Edwards.