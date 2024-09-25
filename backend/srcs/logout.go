package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func (app *App) logout(writer http.ResponseWriter, request *http.Request) {
	var user *User
	var JWT string
	var tokenRequest struct {
		Token string `json:"token"`
	}
	if (funcMsg == 1) {
		fmt.Println(Yellow + "logout function" + Reset)
	}
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	JWT = tokenRequest.Token
	user, err = app.getUserByJWT(JWT)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	} 
	err = app.UserExists(user.Username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return 
	}
	err = app.removeJWT(user.Username)
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
		printUsers(app)
	}
	return
}