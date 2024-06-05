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

func renderTemplate(router *mux.Router, app *App) {
	fs := http.FileServer(http.Dir("../../frontend/srcs"))
    
    http.Handle("/static/", http.StripPrefix("/static/", fs))
	// router.HandleFunc("/", serveTemplate("firstPage.html"))
	router.HandleFunc("/", serveTemplate("presentationPage.html")) // La page de présentation sera servie sur la racine
	router.HandleFunc("/gallery", serveTemplate("gallery.html"))
	router.HandleFunc("/signUp", app.signUp).Methods("POST")
	// router.HandleFunc("/api/login", app.login)

	// fs := http.FileServer(http.Dir("../../frontend/srcs"))
	// http.Handle("/scripts", http.StripPrefix("scripts", fs))
}
