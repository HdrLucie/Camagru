package main

import (
	"fmt"
	"net/http"
)

type Claims struct {
	Username string		`json:"username"`
	UserId   int		`json:"userid"`
	jwt.RegisteredClaims
}

func (app *App) updateUsername(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(Yellow + "Send Reset Link function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}

}