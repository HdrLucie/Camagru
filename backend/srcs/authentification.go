package main

import (
	"github.com/go-mail/mail"
)

func sendMail(user User) {
	// NewMessage creates a new message.
	m := mail.NewMessage()
	m.SetHeader("From", "camagru@mail.fr")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Authentification")
	m.SetBody("text/html", "Test")

	dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d")
	err := dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}