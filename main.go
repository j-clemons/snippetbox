package main

import (
    "log"
    "net/http"
)

// Define a home handler function which write a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
    // check if request URL path matches "/".
    // If not then use http.NotFound() to send a 404
    // must return from the handler or else the handler will continue 
    // executing all code below
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a specific snippet..."))
}

// Add a snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
    // use r.Method to check whether the request is using POST or not.
    if r.Method != http.MethodPost {
        // if it's not, use the w.WriteHeader() method to send a 405 status code
        // and the w.Write() method to write a "Method Not Allowed"
        // response body. We then return from the function so that the 
        // subsequent code is not executed.
        w.Header().Set("Allow", http.MethodPost)

        // use the http.Error() function to send a status code and string body
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Write([]byte("Create a new snippet..."))
}

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
