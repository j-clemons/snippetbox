package main

import (
    "log"
    "net/http"
)

func main() {
    // Use the http.NewServeMux() function to initialize a new servemux, then
    // register the home function as the handler for the "/" URL pattern.
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet/view", snippetView)
    mux.HandleFunc("/snippet/create", snippetCreate)

    // Print a log message to say that the server is starting.
    log.Print("starting server on :4000")

    // Use the http.ListenAndServe() function to start a new web server.
    // We pass in two params: the TCP network address to list to no (ex. :4000)
    // and the servemux we just created. If http.ListenAndServe() returns an 
    // error we use the log.Fatal() to log the error and exit.
    // Note any error returned by http.ListenAndServe() is always non-nil.
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
