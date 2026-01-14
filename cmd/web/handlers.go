package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shaheerkj/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{"./ui/html/pages/base.html", "./ui/html/pages/home.html", "./ui/html/pages/nav.html"}
	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	log.Print(err.Error())
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return

	// }

}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", snippet)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "O snail"
	content := "Climb mount fuji"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
