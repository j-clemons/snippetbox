package main

import (
    "fmt"
    "net/http"
    "strconv"
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
    // extract the value of the id parameter from query string
    // and attempt to convert to int. If not int or less than 1
    // return 404 error
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    // use the fmt.Fprintf() function to interpolate the id with our
    // response and write it to the http.ResponseWriter
    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
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
