package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "github.com/j-clemons/snippetbox/internal/models"

    "github.com/julienschmidt/httprouter"
)

// Define a home handler function which write a byte slice containing
// "Hello from Snippetbox" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // can remove check for r.URL.Path != "/" because httprouter 
    // matches path exactly

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // call the newTemplateData() helper to get the templateData struct
    // containing the 'default' data 
    data := app.newTemplateData(r)
    data.Snippets = snippets

    // use the new render helper
    app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// Add a snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())

    // extract the value of the id parameter from query string
    // and attempt to convert to int. If not int or less than 1
    // return 404 error
    id, err := strconv.Atoi(params.ByName("id"))
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

    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("display the form for creating a new snippet..."))
}

// Add a snippetCreate handler function
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    // no longer need to check if request is POST because this is 
    // done automatically by httprouter

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
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
