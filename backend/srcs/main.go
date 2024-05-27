package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()
	renderTemplate(mux)

	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, mux)
}
