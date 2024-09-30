package main

import (
	"fmt"
	"strconv"
)

func (app *App) printUsers() {
	fmt.Println(Cyan + "List of user : " + Reset)
	fmt.Println(Gray + "- - - - - - - - - - - - - - - - -" + Reset)
	for i, _ := range app.users {
		fmt.Println("Id : " + strconv.Itoa(app.users[i].Id))
		fmt.Println("Email : " + app.users[i].Email)
		fmt.Println("Username : " + app.users[i].Username)
		fmt.Println("Password : " + app.users[i].Password)
		fmt.Println("JWT : " + app.users[i].JWT)
		fmt.Println("authToken : " + app.users[i].AuthToken)
		fmt.Printf("Status : %d\n", app.users[i].authStatus)
		fmt.Println(Gray + "- - - - - - - - - - - - - - - - -" + Reset)
	}
}