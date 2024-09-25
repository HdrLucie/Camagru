package main

import (
	"fmt"
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

// func	getUser(username string, app *App) (*User, error) {
// 	var user User

// 	query := "SELECT id, email, password, authToken FROM Users WHERE username = $1"
// 	row := app.dataBase.QueryRow(query, username)
// 	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.AuthToken)
// 	if err != nil {
// 		fmt.Println(Red + "User doesn't exist" + Reset)
// 		return nil, err
// 	}
// 	return &user, nil
// }

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
		fmt.Println(Yellow + "Get user by JWT" + Reset)
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