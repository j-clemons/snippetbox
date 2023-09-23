package main

import (
    "html/template"
    "path/filepath"
    "time"

    "github.com/j-clemons/snippetbox/internal/models"
)

// define a templateData type to act as a holding structure for any
// dynamic data that we want to pass to our HTML templates
type templateData struct {
    CurrentYear int
    Snippet     models.Snippet
    Snippets    []models.Snippet
}

// create a function that returns a formatted time.Time object
func humanDate(t time.Time) string {
    return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
    // initialize a new map to act as the cache
    cache := map[string]*template.Template{}

    // use the filepath.Glob() func to get a slice of all filepaths that
    // match the pattern "./ui/html/pages/*.tmpl".
    // this essentially gives a slice of all file paths for the application
    // 'page' templates
    pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        // extract the file name from the filepath and assign to name var
        name := filepath.Base(page)

        // The template.FuncMap must be registered with the template set before you
        // call the ParseFiles() method. This means we have to use template.New() to
        // create an empty template set, use the Funcs() method to register the
        // template.FuncMap, and then parse the file as normal.
        ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
        if err != nil {
            return nil, err
        }

        // call ParseGlob() *on this template set* to add any partials
        ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
        if err != nil {
            return nil, err
        }

        // call ParseFiles() *on this template set* to add the page templates
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }

        // add the template set to the map using name of page as the key
        cache[name] = ts
    }

    return cache, nil
}
