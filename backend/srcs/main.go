package main

import (
	"fmt"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"size:100;unique"`
}

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:100"`
	Description string
	Deadline    string
	Status      string
	UserID      uint
}

type Project struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100"`
	Description string
	CreatedAt   string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()
	renderTemplate(mux)

	dsn := "username:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	fmt.Println("Successfully connected to the database")

	// Tu peux maintenant utiliser db pour interagir avec ta base de données
	// Exemple de récupération d'un utilisateur
	var user User
	db.First(&user, 1) // Trouve l'utilisateur avec l'ID 1
	fmt.Println(user)

	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, mux)
}
