package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateData struct {
	Page string
}

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

func renderTemplate(mux *http.ServeMux) {
	mux.HandleFunc("/", serveTemplate("firstPage.html"))
	mux.HandleFunc("/presentationPage", serveTemplate("presentationPage.html")) // La page de présentation sera servie sur la racine
	mux.HandleFunc("/register", serveTemplate("register.html"))

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("../../frontend/srcs"))
	http.Handle("/scripts", http.StripPrefix("scripts", fs))
}
