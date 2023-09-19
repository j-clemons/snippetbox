package main

import (
    "flag"
    "log/slog"
    "net/http"
    "os"
)

// define an application struct to hold the app-wide dependencies
type application struct {
    logger *slog.Logger
}

func main() {
    // define command line flag with name 'addr' 
    // default to 4000
    addr := flag.String("addr", ":4000", "HTTP network address")

    // must parse the flag first so it can read the flag and assign
    // to the variable. Must be called *before* using the addr var or it
    // will just be the default. If it errors application will be terminated
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    // initialize a new instance of the application struct
    // containing the dependencies
    app := &application{
        logger: logger,
    }

    // Print a log message to say that the server is starting.
    logger.Info("starting server", "addr", *addr)

    // Use the http.ListenAndServe() function to start a new web server.
    // We pass in two params: the TCP network address to list to no (ex. :4000)
    // and the servemux we just created. If http.ListenAndServe() returns an 
    // error we use the log.Fatal() to log the error and exit.
    // Note any error returned by http.ListenAndServe() is always non-nil.
    err := http.ListenAndServe(*addr, app.routes())

    logger.Error(err.Error())
    os.Exit(1)
}
