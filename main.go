package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type app struct {
	pages map[string]*template.Template
}

type pageData struct {
	Title string
}

func main() {
	a := &app{pages: loadPages("home", "about", "404")}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /{$}", a.page("home", "Home"))
	mux.HandleFunc("GET /about", a.page("about", "About"))
	mux.HandleFunc("/", a.notFound) // anything that didn't match above

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Each page gets its own template set: base.html + that page.
// That way every page can define its own "content" block.
func loadPages(names ...string) map[string]*template.Template {
	pages := make(map[string]*template.Template)
	for _, name := range names {
		pages[name] = template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/"+name+".html",
		))
	}
	return pages
}

func (a *app) page(name, title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, http.StatusOK, name, pageData{Title: title})
	}
}

func (a *app) notFound(w http.ResponseWriter, r *http.Request) {
	a.render(w, http.StatusNotFound, "404", pageData{Title: "Not found"})
}

// Render into a buffer first, so a template error doesn't send a half-written page.
func (a *app) render(w http.ResponseWriter, status int, name string, data pageData) {
	t, ok := a.pages[name]
	if !ok {
		http.Error(w, "page not found: "+name, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		log.Println("template error:", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}