package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type d struct {
	Username	string	`json:"Username"`
	Id			int		`json:"Id"`
	PId			int		`json:"Photo"`
	Comment		string	`json:"Comment"`
}

func (app *App) manageComment(writer http.ResponseWriter, request *http.Request) {
	var comment	d

	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&comment)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(comment)

}
