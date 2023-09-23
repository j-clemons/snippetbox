package main

import (
    "bytes"
    "fmt"
    "net/http"
)

// the serverError helper writes a log entry at Error level
// then sends a generic 500 Interal Server Error response
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
    var (
        method = r.Method
        uri = r.URL.RequestURI()
    )
    
    app.logger.Error(err.Error(), "method", method, "uri", uri)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper send a specific status code and description
func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

// for consistency we'll implement a notFound helper. 
// simply a wrapper for the 404 response
func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.serverError(w, r, err)
        return
    }

    // initialize a new buffer
    buf := new(bytes.Buffer)

    // write the template to the buffer instead of straight to the
    // http.ResponseWriter. If there is an error, call our serverError()
    // and then return
    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // if the template is written the buffer without any errors
    // we write the HTTP status code to http.ResponseWriter
    w.WriteHeader(status)

    // write the contents of the buffer to the http.ResponseWriter
    // note this is a time where we pass our http.ResponseWriter to a
    // func that takes an io.Writer
    buf.WriteTo(w)
}
