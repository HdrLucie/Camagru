package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"log"
	"database/sql"
)

var Reset = "\033[0m" 
var Red = "\033[31m" 
var Green = "\033[32m" 
var Yellow = "\033[33m" 
var Blue = "\033[34m" 
var Magenta = "\033[35m" 
var Cyan = "\033[36m" 
var Gray = "\033[37m" 
var White = "\033[97m"

var funcMsg = 0
var usersList = 0
var setterMsg = 0
var getterMsg = 0

type App struct {
	dataBase	*sql.DB
	users		[]User
}

type User struct {
	Id          int `json:id`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	JWT			string `json:"token"`
	AuthToken	string`json:"authToken"`
	authStatus	bool`json:"authStatus"`
}

type TemplateData struct {
	Page string
}

func main() {
	port := os.Getenv("BIND_ADDR")
	if port == "" {
		port = "8080"
	}
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(port)
	db := DBConnection()
	app := &App{dataBase: db}
	router := http.NewServeMux()
	renderTemplate(router, app)

	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, router)
}