package main

import (
	"github.com/go-mail/mail"
	"math/rand"
	"fmt"
	"encoding/hex"
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
	content := fmt.Sprintf("Welcome to Canagru, click here : <a href=\"https://localhost:8080/connection?token=%s\"> to verify your account\n", user.AuthToken)
	m.SetBody("text/html", content)

	dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d")
	err := dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}