package main 

import (
	"fmt"
	// "encoding/hex"
	// "encoding/json"
	"net/http"
)

func (app *App) downloadImage(writer http.ResponseWriter, request *http.Request) {

	fmt.Println(Yellow + "Download image" + Reset)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    err := request.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(writer, "Erreur parsing formulaire: "+err.Error(), http.StatusBadRequest)
        return
    }

    file, _, err := request.FormFile("image")
    if err != nil {
        http.Error(writer, "Erreur récupération image: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
	timeStamp := request.FormValue("timestamp");
	userId := request.FormValue("id");
	fmt.Println(timeStamp, userId);

	
}
