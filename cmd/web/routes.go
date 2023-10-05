package main

import (
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/justinas/alice"
)

// the routes() method returns a http.Handler containing app routes
func (app *application) routes() http.Handler {
    // intialize the router
    router := httprouter.New()

    router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.notFound(w)
    })

    // create a file server which serves files out of the "./ui/static"
    // dir. Note that the path given to the http.Dir func is relative
    // to the project directory root.
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

    // create a new middleware chain containing middleware specific
    // to our dynamic application routes. 
    dynamic := alice.New(app.sessionManager.LoadAndSave)

    // Register the other application routes as normal.
    router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
    router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
    router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
    router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

    // create a middleware chain using alice 
    standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

    return standard.Then(router)
}
