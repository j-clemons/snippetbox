package main

import "net/http"

// the routes() method returns a http.Handler containing app routes
func (app *application) routes() http.Handler {
    // Use the http.NewServeMux() function to initialize a new servemux, then
    // register the home function as the handler for the "/" URL pattern.
    mux := http.NewServeMux()

    // create a file server which serves files out of the "./ui/static"
    // dir. Note that the path given to the http.Dir func is relative
    // to the project directory root.
    fileServer := http.FileServer(http.Dir("./ui/static/"))

    // use the mux.Handle() func to register the file server as the handler
    // for all URL paths that start with "/static/".
    // For matching paths, we strip the "/static" prefix before the request
    // reaches the file server.
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    // Register the other application routes as normal.
    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snippet/view", app.snippetView)
    mux.HandleFunc("/snippet/create", app.snippetCreate)

    // pass the servermux to the secureHeaders middleware
    // since the middleware returns a http.Handler nothing else is needed
    return secureHeaders(mux)
}
