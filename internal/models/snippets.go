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

// define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
    DB *sql.DB
}

// insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
    stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    result, err := m.DB.Exec(stmt, title, content, expires)
    if err != nil {
        return 0, err
    }

    // use the LastInsertId() method on the result to get the ID
    // of our newly inserted record in the snippets table
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    // ID is type int64, so needs to be converted before return
    return int(id), nil 
}

// return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (Snippet, error) {
    stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

    // use QueryRow() method on connection pool to execute the statement
    row := m.DB.QueryRow(stmt, id)

    // initialize a new zeroed Snippet struct
    var s Snippet
    
    // use row.Scan() to copy the values from each sql.Row to the 
    // corresponding field in the Snippet struct. NOTICE the args to row.Scan 
    // are *pointers* to the place we want to copy the data. 
    // Number of args must be exactly the same as the columns returned by the statement
    err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
    if err != nil {
        // if the query returned no rows, then row.Scan() will return a
        // sql.ErrNoRows error. Use the errors.Is() func check for that error
        // and return our own ErrNoRecord error instead
        if errors.Is(err, sql.ErrNoRows) {
            return Snippet{}, ErrNoRecord
        } else {
            return Snippet{}, err
        }
    }

    // if everything worked then return the filled Snippet struct
    return s, nil
}

// return the 10 most recent snippets 
// func (m *SnippetModel) Latest() ([]Snippet, error) {
//     return nil, nil
// }
