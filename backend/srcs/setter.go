package main

import (
	"fmt"
)

// ! ||--------------------------------------------------------------------------------||
// ! ||                                  USER SETTERS                                  ||
// ! ||--------------------------------------------------------------------------------||

func (app *App) removeJWT(username string) error {
	if (setterMsg == 1) {
		fmt.Println(Yellow + "Remove JWT function" + Reset)
	}
	_, err := app.dataBase.Exec("UPDATE users SET JWT = $1 WHERE username = $2", "", username) 
	if err != nil {
		fmt.Println(Red + "Error : set confirmed status" + Reset)
		fmt.Println("Error details:", err)
		return err
	}
	for i, u := range app.users {
		if u.Username == username {
			app.users[i].JWT = ""
		}
	}
	return nil
}

func setterStatus(app *App, id int) error {
	result, err := app.dataBase.Exec("UPDATE users SET authStatus = $1 WHERE id = $2", 1, id)
	if err != nil {
		fmt.Println(Red + "Error : set confirmed status" + Reset)
		fmt.Println("Error details:", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(Red + "Error getting rows affected" + Reset)
		return err
	}
	if rowsAffected == 0 {
		fmt.Println(Yellow + "Warning: No rows were updated" + Reset)
	}
	for i, _ := range app.users {
		if app.users[i].Id == id {
			app.users[i].AuthStatus = true
		}
	}
	return nil
}

func (app *App) setPassword(id int, newPassword string) error {

	encryptPassword := encryptPassword(newPassword)
	result, err := app.dataBase.Exec("UPDATE users SET password = $1 WHERE id = $2", encryptPassword, id)
	if err != nil {
		fmt.Println(Red + "Error : set new password" + Reset)
		fmt.Println("Error details:", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(Red + "Error getting rows affected" + Reset)
		return err
	}
	if rowsAffected == 0 {
		fmt.Println(Yellow + "Warning: No rows were updated" + Reset)
	}
	for i, _ := range app.users {
		if app.users[i].Id == id {
			app.users[i].Password = encryptPassword
		}
	}
	return nil
}

func (app *App) setUsername(id int, newUsername string) error {
	_, error := app.getUserByUsername(newUsername)
	if error == nil {
		return error
	}
	result, err := app.dataBase.Exec("UPDATE users SET username = $1 WHERE id = $2", newUsername, id)
	if err != nil {
		fmt.Println(Red + "Error : set username" + Reset)
		fmt.Println("Error details:", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(Red + "Error getting rows affected" + Reset)
		return err
	}
	if rowsAffected == 0 {
		fmt.Println(Yellow + "Warning: No rows were updated" + Reset)
	}
	for i, _ := range app.users {
		if app.users[i].Id == id {
			app.users[i].Username = newUsername
		}
	}
	return nil
}

func (app *App) setEmail(id int, email string) error {
	fmt.Println(Red + "Set email" + Reset)
	_, error := app.getUserByEmail(email)
	if error == nil {
		return error
	}
	result, err := app.dataBase.Exec("UPDATE users SET email = $1 WHERE id = $2", email, id)
	if err != nil {
		fmt.Println(Red + "Error : set email" + Reset)
		fmt.Println("Error details:", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
    if err != nil {
        fmt.Println(Red + "Error getting rows affected" + Reset)
        return err
    }
    if rowsAffected == 0 {
        fmt.Println(Yellow + "Warning: No rows were updated" + Reset)
    }
	for i, _ := range app.users {
		if app.users[i].Id == id {
			app.users[i].Email = email
		}
	}
	return nil
}