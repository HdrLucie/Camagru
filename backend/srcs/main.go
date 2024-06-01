package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DBConnection() {
    cfg := mysql.Config{
        User:   os.Getenv("DB_USER"),
        Passwd: os.Getenv("DB_PASSWORD"),
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: os.Getenv("DB_NAME"),
    }
	fmt.Println(cfg.User)
	fmt.Println(cfg.Passwd)
	fmt.Println(cfg.Net)
	fmt.Println(cfg.Addr)
	fmt.Println(cfg.DBName)
	var err error
	fmt.Println(cfg.FormatDSN())
    db, err = sql.Open("mysql", "hlucie:mdp@tcp(127.0.0.1:3306)/camagru")
    if err != nil {
        log.Fatal(err)
    }

	fmt.Println("2")
    pingErr := db.Ping()
	fmt.Println("3")
    if pingErr != nil {
		fmt.Println("4")
		log.Fatal(pingErr)
    }
    fmt.Println("Connected!")
}

func main() {
	port := os.Getenv("BIND_ADDR")
	if port == "" {
		port = "8080"
	}

	fmt.Println(port)
	DBConnection()
	mux := http.NewServeMux()
	renderTemplate(mux)

	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, mux)
}
