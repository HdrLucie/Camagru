package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-mail/mail"
)

type LikeRequest struct {
	Username string `json:"Username"`
	Id       int    `json:"Id"`
	Photo    int    `json:"Photo"`
}

func (app *App) sendLikes(writer http.ResponseWriter, request *http.Request) {
	var r LikeRequest;

	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&r)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()
	user, err := app.getUserByPhotoId(r.Photo);
	fmt.Println(user) 
	if (user.Notify == true) {
		m := mail.NewMessage();
		m.SetHeader("From", "camagru@mail.fr")
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "Someone liked your post") 
		content := fmt.Sprintf("Hi %s, %s liked your post.", user.Username, r.Username)
		m.SetBody("text/html", content) 
		dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d") 
		err = dialer.DialAndSend(m) 
		if err != nil {
			panic(err);
		}
	}
}
