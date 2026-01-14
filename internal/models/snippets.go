package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets(title,content,created,expires) VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {

	stmt := `select id, title, content, created, expires from snippets where expires > UTC_TIMESTAMP() AND ID = ?`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}

	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `select id, title, content, created, expires from snippets where expires > UTC_TIMESTAMP() order by id desc limit 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err := rows.Scan(&s.ID, &s.Title, &s.Title, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, err

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
