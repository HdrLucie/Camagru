package main

import (
	"fmt"
	"net/http"
	_ "strings"
	"encoding/json"
	_ "github.com/golang-jwt/jwt/v5"
)

func (app *App) modifyUsername(writer http.ResponseWriter, request *http.Request) {
	var userData struct {
		Login string `json:"login"`
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
	newUsername := userData.Login
	err = app.setUsername(user.Id, newUsername)
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

func (app *App) modifyEmail(writer http.ResponseWriter, request *http.Request) {
	var userData struct {
		Email string `json:"email"`
	}
	fmt.Println("ModifyEmail")
	user, ok := request.Context().Value("user").(*User)
	if !ok {
		fmt.Println("Here")
		http.Error(writer, "User not found in context", http.StatusInternalServerError)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&userData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	newEmail := userData.Email
	err = app.setEmail(user.Id, newEmail)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}