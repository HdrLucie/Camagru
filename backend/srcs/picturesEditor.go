package main 

import (
	"fmt"
	// "encoding/hex"
	"encoding/json"
	"net/http"
)

func (app *App) downloadImage(writer http.ResponseWriter, request http.Request) {
	var picture Pictures
	// var token string

	fmt.Println(Yellow + "Download image" + Reset)
	writer.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// NewDecoder.Decode and NewEncoder.Encode encode/décode un JSON -> golang/golang -> JSON. Retourne une structure.
	// Nous permet de travailelr avec du JSON.
	err := json.NewDecoder(request.Body).Decode(&picture)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(picture.Path, picture.JWT)
}
