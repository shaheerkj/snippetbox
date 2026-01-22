# Snippetbox

A secure web application for sharing code snippets, built with Go. Created while following Alex Edwards' book *Let's Go*.

## Features

- **User Authentication** - Sign up, login, and logout with secure password hashing (bcrypt)
- **Create snippets** - Share code snippets with configurable expiration (1 day, 7 days, or 1 year)
- **View snippets** - Browse and view individual code snippets
- **Auto-expiration** - Snippets automatically expire and are hidden after their set duration
- **Session Management** - Server-side sessions stored in MySQL with 12-hour lifetime
- **CSRF Protection** - Cross-site request forgery protection using nosurf
- **HTTPS/TLS** - Secure connections with TLS 1.2+ and modern cipher suites
- **Template caching** - Pre-parsed templates for better performance
- **Middleware chain** - Request logging, panic recovery, authentication, and security headers
- **Form validation** - Server-side validation with user-friendly error messages

## Tech Stack

- **Go 1.24+** - Backend language
- **MySQL** - Database (snippets, users, sessions)
- **html/template** - Server-side templating
- **Alice** - Middleware chaining
- **SCS** - Session management with MySQL store
- **bcrypt** - Secure password hashing
- **nosurf** - CSRF protection

## Project Structure

```
├── cmd/web/           # Application entry point and handlers
│   ├── main.go        # App initialization and server startup
│   ├── handlers.go    # HTTP handlers for routes
│   ├── helpers.go     # Helper functions (error handling, rendering)
│   ├── routes.go      # Route definitions and middleware setup
│   ├── middleware.go  # Custom middleware (logging, auth, security)
│   └── templates.go   # Template caching and custom functions
├── internal/
│   ├── models/        # Database models
│   │   ├── snippets.go  # Snippet CRUD operations
│   │   ├── users.go     # User authentication operations
│   │   └── errors.go    # Custom error types
│   └── validator/     # Form validation utilities
├── tls/               # TLS certificates (cert.pem, key.pem)
└── ui/                # Frontend assets
    ├── html/          # Go templates
    └── static/        # CSS, JS, images
```

## Running the Application

```bash
# Start the server with HTTPS (default port 4000)
go run ./cmd/web

# With custom port
go run ./cmd/web -addr=":8080"

# With custom database DSN
go run ./cmd/web -dsn="user:pass@/snippetbox?parseTime=true"
```

Access the application at: `https://localhost:4000`

## Database Setup

```sql
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE snippetbox;

CREATE TABLE snippets (
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_created (created)
);

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT users_uc_email UNIQUE (email)
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

## TLS Certificate Setup

Generate self-signed certificates for development:

```bash
cd tls
go run "C:\Program Files\Go\src\crypto\tls\generate_cert.go" --rsa-bits=2048 --host=localhost
```

## Routes

| Method | Path | Description | Auth Required |
|--------|------|-------------|---------------|
| GET | `/` | Homepage with latest snippets | No |
| GET | `/snippet/view/{id}` | View a specific snippet | No |
| GET | `/snippet/create` | Display create form | Yes |
| POST | `/snippet/create` | Create new snippet | Yes |
| GET | `/user/signup` | Display signup form | No |
| POST | `/user/signup` | Register new user | No |
| GET | `/user/login` | Display login form | No |
| POST | `/user/login` | Authenticate user | No |
| POST | `/user/logout` | Log out user | Yes |

## Credits

Built following [Let's Go](https://lets-go.alexedwards.net/) by Alex Edwards.