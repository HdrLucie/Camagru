package main

import (
	"github.com/go-mail/mail"
	"math/rand"
	"fmt"
	"encoding/hex"
	"encoding/json"
	"net/http"
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
	m := mail.NewMessage()
	m.SetHeader("From", "camagru@mail.fr")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Authentification")
	content := fmt.Sprintf("Welcome to Canagru, click here : <a href=\"http://localhost:8080/authentification?token=%s\"> to verify your account\n", user.AuthToken)
	m.SetBody("text/html", content)

	dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d")
	err := dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}

func verifyToken(app *App, token string) (int, string) {
	fmt.Println(Yellow + "GET TOKEN" + Reset)
	var user User
	query := "SELECT id, username FROM Users WHERE authToken = $1"
	err := app.dataBase.QueryRow(query, token).Scan(&user.Id, &user.Username)
	fmt.Println(user.Id, user.Username)
	if err != nil {
		return -1, ""
	}
	fmt.Println(user.Id, user.Username)
	return user.Id, user.Username
}

func setConfirmed(app *App, id int) {
	result, err := app.dataBase.Exec("UPDATE users SET confirmed = $1 WHERE id = $2", 1, id)
    if err != nil {
        fmt.Println(Red + "Error : set confirmed status" + Reset)
        fmt.Println("Error details:", err)
    }
	rowsAffected, err := result.RowsAffected()
    if err != nil {
        fmt.Println(Red + "Error getting rows affected" + Reset)
    }
    fmt.Println("Rows affected:", rowsAffected)
}

func (app *App) verifyAuth(writer http.ResponseWriter, request *http.Request) {
	var token string 
	fmt.Println(Red + "VERIFICATION FUNCTION" + Reset)
	var tokenRequest struct {
        Token string `json:"token"`
    }
    err := json.NewDecoder(request.Body).Decode(&tokenRequest)
    if err != nil {
        http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
    }
	token = tokenRequest.Token
	id, username := verifyToken(app, token)
	fmt.Printf("Id : %d\n", id)
	if id == -1 {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
        return
	} else {
		setConfirmed(app, id)
		user, _ := getUser(username, app)
		fmt.Println(Red + "SET CONFIRMED" + Reset)
		fmt.Println(user)
	}
}