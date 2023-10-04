package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "github.com/j-clemons/snippetbox/internal/models"
    "github.com/j-clemons/snippetbox/internal/validator"

    "github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires             int    `form:"expires"`
    validator.Validator `from:"-"`
}

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
    data := app.newTemplateData(r)

    data.Form = snippetCreateForm{
        Expires: 365,
    }

    app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// Add a snippetCreate handler function
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    // declare a new empty instace of the snippetCreateForm struct
    var form snippetCreateForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // because the Validator struct is embedded in the snippetCreateForm
    // struct CheckFiled() can be called directly on it
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

    // if there are any errors, dump them in a plain text HTTP response
    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }

    // pass the data to the SnippetModel.Insert() method
    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // redirect the user to the relevant page for the snippet
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
