package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet represents a code snippet stored in the database
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel wraps a database connection pool
// All database operations for snippets are methods on this type
type SnippetModel struct {
	DB *sql.DB
}

// Insert adds a new snippet to the database and returns its ID
// The expires parameter is the number of days until expiration
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL statement with placeholders (?) to prevent SQL injection
	stmt := `INSERT INTO snippets(title,content,created,expires) 
	         VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`

	// Execute the SQL statement with parameters
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get retrieves a specific snippet by ID
// Returns ErrNoRecord if the snippet doesn't exist or has expired
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// Only return snippets that haven't expired yet
	stmt := `SELECT id, title, content, created, expires 
	         FROM snippets 
	         WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// QueryRow returns at most one row
	row := m.DB.QueryRow(stmt, id)

	// Initialize empty Snippet struct
	var s Snippet

	// Scan the result into the struct fields
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No matching record found
			return Snippet{}, ErrNoRecord
		}
		// Some other database error
		return Snippet{}, err
	}

	return s, nil
}

// Latest returns the 10 most recently created non-expired snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	// Get the 10 most recent snippets that haven't expired
	stmt := `SELECT id, title, content, created, expires 
	         FROM snippets 
	         WHERE expires > UTC_TIMESTAMP() 
	         ORDER BY id DESC 
	         LIMIT 10`

	// Query returns multiple rows
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed when function returns

	// Initialize empty slice to hold snippets
	var snippets []Snippet

	// Iterate through all returned rows
	for rows.Next() {
		var s Snippet
		// Scan each row into a Snippet struct
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Add snippet to slice
		snippets = append(snippets, s)
	}

	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

	// stmt := `select id, title, content, created, expires from snippets where expires > UTC_TIMESTAMP() ORDER BY ID DESC LIMIT 10`

	// row, err := m.DB.Query(stmt)

	// if err != nil {
	// 	return nil, err
	// }

	// defer row.Close()

	// var snippets []Snippet

	// for row.Next() {
	// 	var s Snippet

	// 	err = row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	snippets = append(snippets, s)

	// }
	// if err = row.Err(); err != nil {
	// 	return nil, err
	// }
	// return snippets, err
}
