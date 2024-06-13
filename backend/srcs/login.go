package main

import (
	"github.com/golang-jwt/jwt/v5"
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"time"
)

type Claims struct {
	Username string		`json:"username"`
	UserId   int		`json:"userid"`
	jwt.RegisteredClaims
}

func	getUser(username string, app *App) (*User, error) {
	var user User

	query := "SELECT id, username, email, password FROM Users WHERE username = $1"

	row := app.dataBase.QueryRow(query, username)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		fmt.Println(Red + "User doesn't exist" + Reset)
		return nil, err
	}
	return &user, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(Green + "Error : Wrong password" + Reset)
		return false
	}
	return true
}

func	createToken(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		UserId:   user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func	checkPassword(u User, app *App, writer http.ResponseWriter) (string, int) {

	user, err := getUser(u.Username, app)
	if err != nil {
		redirectPath := "/"
		return redirectPath, http.StatusUnauthorized
	}
	if CheckPasswordHash(u.Password, user.Password) == true {
		fmt.Println(Green + "Right password" + Reset)
		redirectPath := "/gallery"
		return redirectPath, http.StatusOK
	} else {
		fmt.Println(Green + "Wrong password" + Reset)
		redirectPath := "/"
		return redirectPath, http.StatusUnauthorized
	}
}

func (app *App)	login(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(Yellow + "login function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	var user User

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(Red + "Username : " + user.Username + "Password : " + user.Password + Reset)
	err, _ = availableUsername(app, user.Username)
	if err != nil {
		fmt.Println(Red + "Error : wrong username" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	var redirectPath string
	redirectPath, statusCode := checkPassword(user, app, writer)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := createToken(user)
	writer.WriteHeader(statusCode)
    json.NewEncoder(writer).Encode(map[string]string{
        "token": token,
		"redirectPath": redirectPath,
    })
}