package main

import (
	"fmt"
	"net/http"
	"encoding/json"
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
	query := "SELECT id, username, email, password, authToken FROM Users WHERE JWT = $1"
	row := app.dataBase.QueryRow(query, JWT)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.AuthToken)
	if err != nil {
		fmt.Println(err)
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
			return app.users[i].authStatus
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
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.authStatus)
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
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.AuthToken, &user.authStatus)
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

	fmt.Println(Yellow + "Send Reset Link function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
	err := json.NewDecoder(request.Body).Decode(&u)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	return u
}