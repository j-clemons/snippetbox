package models

import (
    "database/sql"
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
// func (m *SnippetModel) Get(id int) (Snippet, error) {
//     return nil, nil
// }

// return the 10 most recent snippets 
// func (m *SnippetModel) Latest() ([]Snippet, error) {
//     return nil, nil
// }
