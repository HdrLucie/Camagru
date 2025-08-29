package main

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"strconv"
	"regexp"
)

var secretKey = []byte("secret-key")

func encryptPassword(password string) (string) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(Red + "Error : Encrypt password" + Reset)
		return ""
	}
	return string(crypted)
}

func isValidEmail(email string) bool {
    re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    return re.MatchString(email)
}

func availableUsername(app *App, username string) (error, bool) {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	err := app.dataBase.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return err, false
	}
	return nil, exists
}

func	availableEmail(app *App, email string) (error, bool) {
	var exists bool
	
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	
	err := app.dataBase.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return err, false
	}
	return nil, exists
}

func	availableId(app *App, id int) (error, bool) {
	var exists bool
	
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	
	err := app.dataBase.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return err, false
	}
	return nil, exists
}

func isIdentifierAvailable(app *App, user User) bool {
	err, id := availableId(app, user.Id)
	if err != nil {
		fmt.Println("Error checking ID availability:", err)
		return false
	}
	err, email := availableEmail(app, user.Email)
	if err != nil {
		fmt.Println("Error checking email availability: ", err)
		return false
	}
	err, username := availableUsername(app, user.Username)
	if err != nil {
		fmt.Println("Error checking username availability:", err)
		return false
	}
	if (id == true || email == true || username == true) {
		return true
	}
	return false
}

func (app *App)	signUp(writer http.ResponseWriter, request *http.Request) {
	var userID int
	var user User
	var token string
	if (funcMsg == 1) {
		fmt.Println(Yellow + "Sign Up function" + Reset)
	}
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// NewDecoder.Decode and NewEncoder.Encode encode/décode un JSON -> golang/golang -> JSON. Retourne une structure.
	// Nous permet de travailelr avec du JSON.
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	encryptPassword := encryptPassword(user.Password)
	user.Password = encryptPassword
	if (isIdentifierAvailable(app, user) == true ) {
		http.Error(writer, "Error : Username or email already in use", http.StatusConflict)
		return	
	}
	if !isValidEmail(user.Email) {
		http.Error(writer, "Invalid email format", http.StatusBadRequest)
		return
	}
	token = generateAuthToken()
	user.AuthToken = token
	fmt.Println(app.stickers[0].Path);
	err = app.dataBase.QueryRow("INSERT INTO users (email, username, password, authToken, authStatus, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Email, user.Username, string(encryptPassword), string(user.AuthToken), 0, "astronaute_avatar.png").Scan(&userID)
	if err != nil {
		fmt.Println(Red + "Error : insert users to the database" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Id = userID
	user.Notify = true
	app.users = append(app.users, user)
	writer.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Account created successfully! Please check your email to verify your account.",
		"id":      strconv.Itoa(user.Id),
		"redirectPath": "/authentification",
	}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		fmt.Println(Red + "Error : Encode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	sendMail(user)
	if (usersList == 1) {
		app.printUsers()
	}
}
