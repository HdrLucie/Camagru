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
	stickers	[]Stickers
	pictures	[]Pictures
	avatars		[]Avatars
	comments	[]Comments
	likes		[]Likes
}

type User struct {
	Id          int `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	JWT			string `json:"token"`
	AuthToken	string`json:"authToken"`
	AuthStatus	bool`json:"authStatus"`
	Avatar		string `json:"avatar"`
	Notify		bool`json:"notify"`
}

type Stickers struct {
	Id		int `json:"id"`
	Name	string `json:"name"`
	Path	string `json:"path"`
}

type Pictures struct {
	Path		string `json:"path"`
	Id			int `json:"id"`
	UserId		int `json:"userId"`
	UploadTime	string `json:"uploadTime"`
	Likes		int `json:"likes"`
	Comments	int `json:"comments"`
}

type TemplateData struct {
	Page string
}

type Avatars struct {
	Id		int `json:"id"`
	Name	string `json:"name"`
	Path	string `json:"path"`
}

type Comments struct {
	Username	string	`json:"Username"`
	Comment		string	`json:"Comment"`
	PId			int		`json:"pId"`
}

type Likes struct {
	Username	string	`json:"Username"`
	PId			int		`json:"pId"`
	UId			int		`json:"uId"`
}

var MailPwd string;

func main() {
	port := os.Getenv("BIND_ADDR")
	if port == "" {
		port = "8080"
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Pas de fichier .env, utilisation des variables d'environnement Docker")
	}
	pwd, ok := os.LookupEnv("PASSWORD")
	if !ok || len(pwd) == 0 {
		log.Fatal("PASSWORD variable not set")
	}
	MailPwd = pwd;

	db := DBConnection()
	app := &App{dataBase: db}
	router := http.NewServeMux()
	err = app.InsertSticker()
	if err != nil {
		log.Fatal("Error loading stickers directory")
	}
	err = app.createUploadsDirectory();
	if err != nil {
		log.Fatal("Error creating uploads directory");
	}	
	renderTemplate(router, app);

	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, router)
}
