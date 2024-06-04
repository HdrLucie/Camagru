package main 

import (
	"encoding/json"
    "net/http"
)

var users []User

type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func (app *App)	createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Insérer l'utilisateur dans la base de données
    result, err := app.dataBase.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", user.Email, user.Username, user.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Récupérer l'ID de l'utilisateur inséré
    userID, _ := result.LastInsertId()
    user.ID = int(userID)

    // Ajouter l'utilisateur à la slice
    users = append(users, user)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}