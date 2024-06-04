package main

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
)

type DBConfig struct {
	Host		string
	Port		int
	User		string
	Password	string
	Name		string
	SSLMode		string
}

func getDBConfig() *DBConfig {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		panic("DB_HOST variable not set")
	}
	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		panic("DB_PORT variable not set")
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		panic("DB_USER variable not set")
	}
	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		panic("DB_PASSWORD variable not set")
	}
	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		panic("DB_NAME variable not set")
	}
	sslmode, ok := os.LookupEnv("DB_SSLMODE")
	if !ok {
		panic("DB_SSLMode variable not set")
	}
	return &DBConfig {
		Host: host,
		Port: portInt,
		User: user,
		Password: password,
		Name: dbname,
		SSLMode: sslmode,
	}
}

func DBConnection() *DBConfig {
	dbConfig := getDBConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+ 
	"password=%s dbname=%s sslmode=%s", 
	dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return dbConfig
}
