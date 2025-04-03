package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
)

// ! ||--------------------------------------------------------------------------------||
// ! ||                                  USER GETTERS                                  ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) UserExists(username string) (error) {
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Check if an user exists in DB" + Reset)
	}
	var err error
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`
	app.dataBase.QueryRow(query, username).Scan(&err)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (app *App) getUserByJWT(JWT string) (*User, error) {
	var user User
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get user by JWT" + Reset)
	}
	query := "SELECT id, email, username, password, JWT, authToken, authStatus, avatar FROM Users WHERE JWT = $1"
	row := app.dataBase.QueryRow(query, JWT)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.JWT, &user.AuthToken, &user.AuthStatus, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (app *App) getStatus(id int) (bool) {
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get Status" + Reset)
	}	
	for i, _ := range app.users {
		if app.users[i].Id == id {
			return app.users[i].AuthStatus
		}
	}
	return false
}

func (app *App) getUserByUsername(username string) (*User, error) {
	var user User
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get user by username" + Reset)
	}
	query := "SELECT id, email, username, password, authToken, authStatus FROM Users WHERE username = $1"
	row := app.dataBase.QueryRow(query, username)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.AuthStatus)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (app *App) getUserByEmail(email string) (*User, error) {
	var user User
	if (getterMsg == 1) {
		fmt.Println(Yellow + "Get user by email" + Reset)
	}
	query := "SELECT id, email, username, password, authToken, authStatus FROM Users WHERE email = $1"
	row := app.dataBase.QueryRow(query, email)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.AuthStatus)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                  FRONT GETTERS                                 ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) deserializeUserData(writer http.ResponseWriter, request *http.Request) User {
	var u User

	fmt.Println(Yellow + "deserializeUserData function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
	err := json.NewDecoder(request.Body).Decode(&u)
	fmt.Println(u)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	return u
}

func extractJWTFromRequest(request *http.Request) string {
	JWT := request.Header.Get("Authorization")
	return strings.TrimPrefix(JWT, "Bearer ")
}

func (app *App) getUser(writer http.ResponseWriter, request *http.Request) {
	token := extractJWTFromRequest(request)
	user, _ := app.getUserByJWT(token)
	user.Password = ""
	user.AuthToken = ""
	writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(user)
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                 GETTER STICKERS                                ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) getStickers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(app.stickers)
}

// ! ||--------------------------------------------------------------------------------||
// ! ||                                 PICTURES GETTER                                ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) getPictures(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(app.pictures)
}
