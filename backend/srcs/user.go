package main 

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func encryptPassword(password string) (string, error) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(crypted), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (app *App)	signUp(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user User
	// NewDecoder.Decode and NewEncoder.Encode encode/décode un JSON -> golang/golang -> JSON. Retourne une structure.
	// Nous permet de travailelr avec du JSON.
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	encryptPassword, err := encryptPassword(user.Password)
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	fmt.Println(encryptPassword)
	result, err := app.dataBase.Exec("INSERT INTO users (email, username, password) VALUES ($1, $2, $3)", user.Email, user.Username, string(encryptPassword))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	userID, _ := result.LastInsertId()
	user.Id = int(userID)
	app.users = append(app.users, user)
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(map[string]string{"message": "User created successfully"})
}

func getUser(app *App, username string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	err := app.dataBase.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (app *App)	login(writer http.ResponseWriter, request *http.Request) {
	var user User

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	exist, err := getUser(app, user.Username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var redirectPath string
	if exist {
		redirectPath = "/gallery"
		writer.WriteHeader(http.StatusOK)
	} else {
		redirectPath = "/"
		writer.WriteHeader(http.StatusUnauthorized)
	}

	json.NewEncoder(writer).Encode(map[string]string{"redirectPath": redirectPath})
}