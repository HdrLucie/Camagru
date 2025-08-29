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

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(Green + "Error : Wrong password" + Reset)
		return false
	}
	return true
}

func	createToken(user *User) (string, error) {
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

func (app *App)	checkPassword(user *User, pass string, writer http.ResponseWriter) (string, int) {
	if CheckPasswordHash(pass, user.Password) == true {
		redirectPath := "/gallery"
		return redirectPath, http.StatusOK
	} else {
		redirectPath := "/"
		return redirectPath, http.StatusUnauthorized
	}
}

func addTokenToDb(app *App, user *User, token string) error {
	var exists bool
	
	if (funcMsg == 1) {
		fmt.Println(Yellow + "Add token to database" + Reset)
	}
	err := app.dataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", user.Username).Scan(&exists)
	if err != nil {
		fmt.Println("Erreur lors de la vérification de l'existence de l'utilisateur:", err)
		return err
	}
	if !exists {
		fmt.Printf("Aucun utilisateur trouvé avec l'Username : %d\n", user.Username)
		return fmt.Errorf("utilisateur non trouvé")
	}
    result, err := app.dataBase.Exec("UPDATE users SET JWT = $1 WHERE username = $2", token, user.Username)
    if err != nil {
        fmt.Println(Red + "Error : add token to database" + Reset)
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
    for i, u := range app.users {
        if u.Id == user.Id {
            app.users[i].JWT = token
        }
    }
	return nil
}

func (app *App) manageLoginError(pass string, user *User, writer http.ResponseWriter) (string, string, int) {
	err, _ := availableUsername(app, user.Username)
	if err != nil {
		fmt.Println(Red + "Error : wrong username" + Reset)
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return "", "/connection", 401
	}
	redirectPath, statusCode := app.checkPassword(user, pass, writer)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return "", "/connection", 401
	}
	if user.AuthStatus == true {
		token, err := createToken(user)
		if err != nil {
			fmt.Println(Red + "Error : creating token" + Reset)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return "", "", 500
		}
		return token, redirectPath, statusCode
	} else {
		return "", "/connection", 403
	}
}

func (app *App)	login(writer http.ResponseWriter, request *http.Request) {
	var user *User
	var data User
	if (funcMsg == 1) {
		fmt.Println(Yellow + "login function" + Reset)
	}
	if (usersList == 1) {
		app.printUsers()
	}
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	user, err = app.getUserByUsername(data.Username)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}
	token, redirectPath, statusCode := app.manageLoginError(data.Password, user, writer)
	addTokenToDb(app, user, token)
	writer.WriteHeader(statusCode)
    json.NewEncoder(writer).Encode(map[string]string{
        "token": token,
		"redirectPath": redirectPath,
    })
}
