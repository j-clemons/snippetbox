package main

import (
    "database/sql"
    "flag"
    "log/slog"
    "net/http"
    "os"

    "github.com/j-clemons/snippetbox/internal/models"

    _ "github.com/go-sql-driver/mysql"
)

// define an application struct to hold the app-wide dependencies
type application struct {
    logger   *slog.Logger
    snippets *models.SnippetModel
}

func main() {
    // define command line flag with name 'addr' 
    // default to 4000
    addr := flag.String("addr", ":4000", "HTTP network address")

    // define a new command line flag for the MySQL DSN String
    dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")

    // must parse the flag first so it can read the flag and assign
    // to the variable. Must be called *before* using the addr var or it
    // will just be the default. If it errors application will be terminated
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    // defer a call to db.Close() so that the connection pool is closed
    // before the main() function exits
    defer db.Close()

    // initialize a new instance of the application struct
    // containing the dependencies
    app := &application{
        logger:   logger,
        snippets: &models.SnippetModel{DB: db},
    }

    // Print a log message to say that the server is starting.
    logger.Info("starting server", "addr", *addr)

    // Use the http.ListenAndServe() function to start a new web server.
    // We pass in two params: the TCP network address to list to no (ex. :4000)
    // and the servemux we just created. If http.ListenAndServe() returns an 
    // error we use the log.Fatal() to log the error and exit.
    // Note any error returned by http.ListenAndServe() is always non-nil.
    err = http.ListenAndServe(*addr, app.routes())

    logger.Error(err.Error())
    os.Exit(1)
}

// the openDB() func wraps sql.Open() and returns a sql.DB connection
// pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
