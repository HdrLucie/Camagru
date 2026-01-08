package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-mail/mail"
	"database/sql"
)

type LikeResponse struct {
	Liked     bool `json:"liked"`
}

type LikeRequest struct {
	Username string `json:"Username"`
	Id       int    `json:"Id"`
	Photo    int    `json:"Photo"`
}

func (app *App)	checkLikeValidity(Uid int, Pid int) (bool, error) {
	var exists int
	err := app.dataBase.QueryRow("SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2", Uid, Pid,).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (app *App) insertLikeIntoDB(Uid int, Pid int, user *User, r LikeRequest) (error, bool) {
	liked, err := app.checkLikeValidity(Uid, Pid);
	if err != nil {
		return err, liked
	}
	if liked {
		_, err = app.dataBase.Exec("DELETE FROM likes WHERE user_id = $1 AND post_id = $2", Uid, Pid);
		if err != nil {
			fmt.Println(Red + "Error : delete like to the database" + Reset)
			return err, liked
		}
		_, err = app.dataBase.Exec("UPDATE images SET like_count = like_count - 1 WHERE id = $1", Pid);
		if err != nil {
			fmt.Println(Red + "Error : delete like_count to the database" + Reset)
			return err, liked
		}

	} else {
		fmt.Println(Green + "Insert Like" + Reset)
		_, err = app.dataBase.Exec("INSERT INTO likes (post_id, user_id) VALUES ($1, $2)", Pid, Uid);
		if err != nil {
			fmt.Println(Red + "Error : insert like to the database" + Reset)
			return err, liked
		}
		_, err = app.dataBase.Exec("UPDATE images SET like_count = like_count + 1 WHERE id = $1", Pid)
		if err != nil {
			fmt.Println(Red + "Error : set count likes" + Reset)
			return err, liked
		}
		if (user.Notify == true) {
			fmt.Println(Red + "Email" + Reset)
			m := mail.NewMessage();
			m.SetHeader("From", "camagru@mail.fr")
			m.SetHeader("To", user.Email)
			m.SetHeader("Subject", "Someone liked your post") 
			content := fmt.Sprintf("Hi %s, %s liked your post : http://localhost:8080/photo/%d", user.Username, r.Username, r.Photo)
			m.SetBody("text/html", content) 
			dialer := mail.NewDialer("smtp.mail.fr", 587, "camagru@mail.fr", "12hdkHUDH![d") 
			err = dialer.DialAndSend(m) 
			if err != nil {
				panic(err);
			}
		}
		return nil, liked
	}
	return nil, liked
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
	err, liked := app.insertLikeIntoDB(r.Id, r.Photo, user, r);

	if err != nil {
		fmt.Println(Red + "Error with like management" + Reset);
		return 
	}
	resp := LikeResponse{
		Liked:     liked,
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}
