package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/gorilla/mux"
)

func serveTemplate(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmplPath := filepath.Join("../../frontend/srcs/templates", templateName)
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Could not parse template", http.StatusInternalServerError)
			fmt.Println("Error parsing template:", err)
			return
		}

		data := TemplateData{Page: templateName}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Could not execute template", http.StatusInternalServerError)
			fmt.Println("Error executing template:", err)
		}
	}
}

func	serveStyleFiles() {
    styles := http.FileServer(http.Dir("../../frontend/srcs/stylesheets/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
}

func serveScriptsFiles() {
    scripts := http.FileServer(http.Dir("../../frontend/srcs/scripts/"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
}

func serveImgFiles() {
    assets := http.FileServer(http.Dir("../../frontend/srcs/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
}

func renderTemplate(router *mux.Router, app *App) {
	serveStyleFiles()
	serveScriptsFiles()
	serveImgFiles()
    router.HandleFunc("/", serveTemplate("login.html"))
    router.HandleFunc("/gallery", serveTemplate("gallery.html"))
    router.HandleFunc("/signUp", app.signUp).Methods("POST")
    router.HandleFunc("/login", app.login).Methods("POST")
}