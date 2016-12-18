package main

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once sync.Once
	filename string
	temp1 *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w,nil)
}

func main() {
	http.Handle("/", &templateHandler{filename : "chat.html"})

	// Start Web Server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
