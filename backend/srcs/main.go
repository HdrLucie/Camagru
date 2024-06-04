package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"log"
)

type App struct {
	dataBase *DBConfig
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
	mux := http.NewServeMux()
	renderTemplate(mux, app)

	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, mux)
}