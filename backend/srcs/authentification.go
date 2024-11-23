package main

import (
	"github.com/go-mail/mail"
	"math/rand"
	"fmt"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
)

func	generateAuthToken() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(token)
}

func sendMail(user User) {
	// NewMessage creates a new message.
	if (funcMsg == 1) {
		fmt.Println(Yellow + "Send mail function" + Reset)
	}
	m := mail.NewMessage()
	m.SetHeader("From", "camagru@mail.fr")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Authentification")
	content := fmt.Sprintf("Welcome to Camagru %s, click here : <a href=\"http://localhost:8080/verify?token=%s\"> to verify your account\n", user.Username, user.AuthToken)
	m.SetBody("text/html", content)

	dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d")
	err := dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}

func verifyToken(app *App, token string) int {
	var user User
	query := "SELECT id FROM Users WHERE authToken = $1"
	err := app.dataBase.QueryRow(query, token).Scan(&user.Id)
	if err != nil {
		fmt.Println(Magenta + "Verify token error" + Reset)
		fmt.Println("Erreur:", err)
		return -1
	}
	return user.Id
}

func (app *App) verifyAccount(writer http.ResponseWriter, request *http.Request) {
	var token string 
	if (funcMsg == 1) {
		fmt.Println(Yellow + "Verify auth func" + Reset)
	}
	var tokenRequest struct {
		Token string `json:"token"`
	}
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&tokenRequest)
	fmt.Println("Token : " + tokenRequest.Token)
	if err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	token = tokenRequest.Token
	id := verifyToken(app, token)
	if id == -1 {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	} else {
		err := setterStatus(app, id)
		if err != nil {
			fmt.Println("Error details:", err)
		}
	}
	if (usersList == 1) {
		app.printUsers()
	}
	writer.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Authentification successfully confirmed!",
		"id":      strconv.Itoa(id),
		"redirectPath": "/connection",
	}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		fmt.Println(Red + "Error : Encode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}