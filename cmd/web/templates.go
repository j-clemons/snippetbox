package main

import (
    "html/template"
    "path/filepath"

    "github.com/j-clemons/snippetbox/internal/models"
)

// define a templateData type to act as a holding structure for any
// dynamic data that we want to pass to our HTML templates
type templateData struct {
    Snippet  models.Snippet
    Snippets []models.Snippet
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

        files := []string{
            "./ui/html/base.tmpl",
            "./ui/html/partials/nav.tmpl",
            page,
        }

        // parse the files into a template set
        ts, err := template.ParseFiles(files...)
        if err != nil {
            return nil, err
        }

        // add the template set to the map using name of page as the key
        cache[name] = ts
    }

    return cache, nil
}
