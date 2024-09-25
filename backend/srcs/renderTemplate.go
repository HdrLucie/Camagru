package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
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

func	serveStyleFiles(router *http.ServeMux) {
    styles := http.FileServer(http.Dir("../../frontend/srcs/stylesheets/"))
	router.Handle("/styles/", http.StripPrefix("/styles", styles))
}

func serveScriptsFiles(router *http.ServeMux) {
	scripts := http.FileServer(http.Dir("../../frontend/srcs/scripts/"))
	router.Handle("/scripts/", http.StripPrefix("/scripts", scripts))
}

func serveImgFiles(router *http.ServeMux) {
    assets := http.FileServer(http.Dir("../../frontend/srcs/assets/"))
	router.Handle("/assets/", http.StripPrefix("/assets", assets))
}

func mdw(next http.Handler) http.Handler {
    f := func(w http.ResponseWriter, r *http.Request) {
        // Executes middleware logic here...
        fmt.Println()
        fmt.Println(r)
        fmt.Println()
        next.ServeHTTP(w, r) // Pass request to next handler
    }

    return http.HandlerFunc(f)
}

func renderTemplate(router *http.ServeMux, app *App) {
	serveStyleFiles(router)
	serveScriptsFiles(router)
	serveImgFiles(router)
    
	router.HandleFunc("/", serveTemplate("firstPage.html"))
	router.HandleFunc("/connection", serveTemplate("login.html"))
    router.HandleFunc("/gallery", serveTemplate("gallery.html"))
	router.HandleFunc("/takePicture", serveTemplate("picture.html"))
	router.HandleFunc("/authentification", serveTemplate("authentification.html"))
    router.HandleFunc("/signUp", app.signUp)
    router.HandleFunc("/login", app.login)
	router.HandleFunc("/verify", serveTemplate("verify.html"))
	router.HandleFunc("/logout", app.logout)
	router.HandleFunc("/verifyAccount", app.verifyAccount)
}