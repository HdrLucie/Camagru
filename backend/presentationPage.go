package main

import "net/http"

func presentationPageHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/presentationPage/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "presentationPage", p)
}
