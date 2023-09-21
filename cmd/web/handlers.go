package main

import (
    "errors"
    "fmt"
    "html/template"
    "net/http"
    "strconv"

    "github.com/j-clemons/snippetbox/internal/models"
)

// Define a home handler function which write a byte slice containing
// "Hello from Snippetbox" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // check if request URL path matches "/".
    // If not then use notFound() helper to send a 404
    // must return from the handler or else the handler will continue 
    // executing all code below
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // intialize a slice containing the paths to the template files
    // IMPORTANT. the base template must be *first* file in the slice
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    // use the template.ParseFIles() func to read the template file
    // if there is an error, we log the detailed error message and
    // use the serverError() helper
    ts, err := template.ParseFiles(files...)
    if err != nil {
        // because home handler is now a method against the application
        // struct it can access its fields, included the structured logger.
        app.serverError(w, r, err)
        return
    }

    data := templateData{
        Snippets: snippets,
    }

    // then use the Execute() method on the template set to write the
    // template conent as the response body. The last param to Execute()
    // represents any dynamic data that we want to pass in, which for now
    // we leave as nil
    err = ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.serverError(w, r, err)
    }
}

// Add a snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    // extract the value of the id parameter from query string
    // and attempt to convert to int. If not int or less than 1
    // return 404 error
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }
    
    // Use the SnippetModel's Get() method to retrieve the data for a specific
    // record based on its ID. If no matching record, then return 404
    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    // initialize a slice containing the paths to the view.tmpl file,
    // plus the base layout and navigation partial
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/view.tmpl",
    }

    // parse the template files
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // create instance of a templateData struct holding the snippet data
    data := templateData{
        Snippet: snippet,
    }

    // execute templates
    err = ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.serverError(w, r, err)
    }
}

// Add a snippetCreate handler function
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    // use r.Method to check whether the request is using POST or not.
    if r.Method != http.MethodPost {
        // if it's not, use the w.WriteHeader() method to send a 405 status code
        // and the w.Write() method to write a "Method Not Allowed"
        // response body. We then return from the function so that the 
        // subsequent code is not executed.
        w.Header().Set("Allow", http.MethodPost)

        // use the clientError() helper
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

    // creating dummy data
    title := "0 snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7

    // pass the data to the SnippetModel.Insert() method
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // redirect the user to the relevant page for the snippet
    http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
    // w.Write([]byte("Create a new snippet..."))
}
