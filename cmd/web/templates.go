package main

import (
    "html/template"
    "io/fs"
    "path/filepath"
    "time"

    "github.com/j-clemons/snippetbox/internal/models"
    "github.com/j-clemons/snippetbox/ui"
)

// define a templateData type to act as a holding structure for any
// dynamic data that we want to pass to our HTML templates
type templateData struct {
    CurrentYear     int
    Snippet         models.Snippet
    Snippets        []models.Snippet
    Form            any
    Flash           string
    IsAuthenticated bool
    CSRFToken       string
}

// create a function that returns a formatted time.Time object
func humanDate(t time.Time) string {
    if t.IsZero() {
        return ""
    }

    return t.UTC().Format("02 Jan 2006 at 15:04")
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
    pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        // extract the file name from the filepath and assign to name var
        name := filepath.Base(page)

        patterns := []string{
            "html/base.tmpl",
            "html/partials/*.tmpl",
            page,
        }

        // The template.FuncMap must be registered with the template set before you
        // call the ParseFiles() method. This means we have to use template.New() to
        // create an empty template set, use the Funcs() method to register the
        // template.FuncMap, and then parse the file as normal.
        ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
        if err != nil {
            return nil, err
        }

        // add the template set to the map using name of page as the key
        cache[name] = ts
    }

    return cache, nil
}
