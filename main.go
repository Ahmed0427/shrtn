package main

import (
	"log"
	"net/http"
	"fmt"
	"os"

	"github.com/ahmed0427/shrtn/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")
	router := server.NewRouter(os.Getenv("CONN_STR"))
	fmt.Printf("URL shortener service running at http://localhost:%s\n", PORT)
	log.Fatal(http.ListenAndServe(":" + PORT, router))
}
