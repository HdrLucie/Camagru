package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func (app *App) logout(writer http.ResponseWriter, request *http.Request) {
	if (funcMsg == 1) {
		fmt.Println(Yellow + "logout function" + Reset)
	}
	user, ok := request.Context().Value("user").(*User)
	if !ok {
		http.Error(writer, "User not found in context", http.StatusInternalServerError)
		return
	}
	err := app.removeJWT(user.Username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]string{
		"message": "Logout successful",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
	if (usersList == 1) {
		app.printUsers()
	}
	return
}