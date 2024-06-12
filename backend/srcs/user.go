package main 

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"strconv"
	"regexp"
)

func encryptPassword(password string) (string) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(Red + "Error : Encrypt password" + Reset)
		return ""
	}
	return string(crypted)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(Green + "Error : Wrong password" + Reset)
		return false
	}
	return true
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
	fmt.Println(exists)
	return nil, exists
}

func	availableEmail(app *App, email string) (error, bool) {
	var exists bool
	
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	
	err := app.dataBase.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return err, false
	}
	fmt.Println(exists)
	return nil, exists
}

func	availableId(app *App, id int) (error, bool) {
	var exists bool
	
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	
	err := app.dataBase.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return err, false
	}
	fmt.Println(exists)
	return nil, exists
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
	fmt.Println(Yellow + "Sign Up function" + Reset)
	writer.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPost {
		fmt.Println(Red + "Error : Method" + Reset)
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user User

	// NewDecoder.Decode and NewEncoder.Encode encode/décode un JSON -> golang/golang -> JSON. Retourne une structure.
	// Nous permet de travailelr avec du JSON.
	err := json.NewDecoder(request.Body).Decode(&user)
	fmt.Println("Email : " + user.Email)
	if err != nil {
		fmt.Println(Red + "Error : Decode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	encryptPassword := encryptPassword(user.Password)
	fmt.Println(Blue + user.Email, user.Username, user.Password + Reset)
	if (isIdentifierAvailable(app, user) == true ) {
		fmt.Println(Red + "Error : Username or email already in use" + Reset)
		http.Error(writer, "Error : Username or email already in use", http.StatusInternalServerError)
		return		
	}
	if !isValidEmail(user.Email) {
		http.Error(writer, "Invalid email format", http.StatusBadRequest)
		return
	}
	var userID int
	err = app.dataBase.QueryRow("INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id", user.Email, user.Username, string(encryptPassword)).Scan(&userID)
	if err != nil {
		fmt.Println(Red + "Error : insert users to the database" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Id = userID
	app.users = append(app.users, user)
	writer.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "User created successfully",
		"id":      strconv.Itoa(user.Id),
	}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		fmt.Println(Red + "Error : Encode Json object" + Reset)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
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
	try, _ := getUser(user.Username, app)
	if CheckPasswordHash(user.Password, try.Password) == true {
		fmt.Println(Green + "Right password" + Reset)
		redirectPath = "/gallery"
		writer.WriteHeader(http.StatusOK)
	} else {
		fmt.Println(Green + "Wrong password" + Reset)
		redirectPath = "/"
		writer.WriteHeader(http.StatusUnauthorized)
	}

	json.NewEncoder(writer).Encode(map[string]string{"redirectPath": redirectPath})
}