package main 

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-mail/mail"
)

func (app *App) deserializeUserData(writer http.ResponseWriter, request *http.Request) User {
	var u User

	fmt.Println(Yellow + "Send Reset Link function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
	err := json.NewDecoder(request.Body).Decode(&u)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	return u
}

func (app *App) sendResetLink(writer http.ResponseWriter, request *http.Request) {
	var user User

	if (funcMsg == 1) {
		fmt.Println(Yellow + "Send Reset Link function" + Reset)
	}
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	m := mail.NewMessage()
	m.SetHeader("From", "camagru@mail.fr")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Reset password")
	content := fmt.Sprintf("Hi %s, It looks like you requested to reset your password. No worries — just click the link below to set up a new one: <a href=\"http://localhost:8080/resetPassword?email=%s\">here</a> to reset your password\n", user.Username, user.Email)
	m.SetBody("text/html", content)

	dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d")
	err = dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}

func (app *App) resetPassword(writer http.ResponseWriter, request *http.Request) {

	u := app.deserializeUserData(writer, request)
	user, err := app.getUserByEmail(u.Email)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
        response := map[string]string{"error": "User not found"}
        json.NewEncoder(writer).Encode(response)
		return 
	}
	printUsers(app)
	err = app.newPassword(user.Id, u.Password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
        response := map[string]string{"error": "Impossible to change password"}
        json.NewEncoder(writer).Encode(response)
	} else {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{
			"message": "User updated successfully",
			"redirectPath": "/connection",
		})
	}
	printUsers(app)	
}