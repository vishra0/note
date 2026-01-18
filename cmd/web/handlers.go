package main

import (
	"fmt"
	"html/template"

	"net/http"
	"strconv"
	"vis/note/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notfound(w)
		return
	}
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, snippet := range s {
		fmt.Fprintf(w, "%v \n", snippet)
	}

}
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notfound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notfound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Snippet: s}
	files := []string{
		"./ui/html/base.layout.tmpl",
		"./ui/html/show.page.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	title := "o snail"
	content := "o snail run away \n slowly slowly"
	expires := "21"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
