package main

import (
	"fmt"
	"net/http"
	_ "strings"
	"encoding/json"
	_ "github.com/golang-jwt/jwt/v5"
)

func (app *App) modifyProfile(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(Red + "MODIFY PROFILE" + Reset)
	var userData struct {
		Login string `json:"username"`
		Email string `json:"email"`
		NotifyState bool `json:"notifyState"`
	}
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Modify profile", userData.Login, userData.Email, userData.NotifyState)
	user, ok := request.Context().Value("user").(*User)
	if !ok {
		http.Error(writer, "User not found in context", http.StatusInternalServerError)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&userData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	app.setNotifyState(userData.NotifyState, user);
	app.modifyUsername(user, userData.Login, writer, request);
	app.modifyEmail(user, userData.Email, writer, request);
}

func (app *App) modifyUsername(user *User, login string, writer http.ResponseWriter, request *http.Request) {
	if (login == "") {
		return
	}
	err := app.setUsername(user.Id, login)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return		
	}
}

func (app *App) modifyPassword(writer http.ResponseWriter, request *http.Request) {
	var userData struct {
		Password string `json:"password"`
	}
	user, ok := request.Context().Value("user").(*User)
	if !ok {
		http.Error(writer, "User not found in context", http.StatusInternalServerError)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&userData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	newPassword := userData.Password
	err = app.setPassword(user.Id, newPassword)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (app *App) modifyEmail(user *User, email string, writer http.ResponseWriter, request *http.Request) {
	if (email == "") {
		return;
	}
	err := app.setEmail(user.Id, email)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}
