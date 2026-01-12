package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/go-mail/mail"
)

type d struct {
	Username	string	`json:"Username"`
	Id			int		`json:"Id"`
	PId			int		`json:"Photo"`
	Comment		string	`json:"Comment"`
}

type response struct {
	Username	string `json:"Username"`
	Comment		string `json:"Comment"`
}

func (app *App) insertCommentIntoDB(comment d) {
	_, err := app.dataBase.Exec("INSERT INTO comments (comment, post_id, user_id) VALUES ($1, $2, $3)", comment.Comment, comment.PId, comment.Id)
	if err != nil {
		fmt.Println(Red + "Error insert comment" + Reset);
	}
	_, err = app.dataBase.Exec("UPDATE images SET comment_count = comment_count + 1 WHERE id = $1", comment.PId);
}

func	sendEmail(user *User) {
	m := mail.NewMessage();
	m.SetHeader("From", "camagru@mail.fr")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "New comment") 
	content := fmt.Sprintf("Hi %s, someone commented your post", user.Username)
	m.SetBody("text/html", content) 
	dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d") 
	err := dialer.DialAndSend(m) 
	if err != nil {
		panic(err);
	}
}

func (app *App) manageComment(writer http.ResponseWriter, request *http.Request) {
	var comment	d;
	var c Comments;

	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&comment)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	c.Comment = comment.Comment;
	c.Username = comment.Username;
	app.comments = append(app.comments, c);
	r := response{
		Username:	comment.Username,
		Comment:	comment.Comment,
	}
	app.insertCommentIntoDB(comment);
	user, err := app.getUserByPhotoId(comment.PId);
	if (user.Notify == true) {
		sendEmail(user);
	}
	writer.Header().Set("Content-Type", "application/json");
	json.NewEncoder(writer).Encode(r);
}
